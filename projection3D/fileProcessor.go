package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ObjModel struct {
	Points []Point
	Faces  [][]int
}

// LoadOBJ reads an .obj file and returns points and faces in the format used by this project.
// Note: OBJ files use 1-based indices, this converts them to 0-based.
func LoadOBJ(filename string) (*ObjModel, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	model := &ObjModel{
		Points: []Point{},
		Faces:  [][]int{},
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "v":
			// Vertex: v x y z
			if len(fields) < 4 {
				return nil, fmt.Errorf("line %d: vertex needs 3 coordinates", lineNum)
			}
			x, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid x coordinate: %w", lineNum, err)
			}
			y, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid y coordinate: %w", lineNum, err)
			}
			z, err := strconv.ParseFloat(fields[3], 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid z coordinate: %w", lineNum, err)
			}
			model.Points = append(model.Points, Point{x, y, z})

		case "f":
			// Face: f v1/vt1/vn1 v2/vt2/vn2 v3/vt3/vn3 ...
			// We only care about the vertex index (first number before /)
			if len(fields) < 4 {
				return nil, fmt.Errorf("line %d: face needs at least 3 vertices", lineNum)
			}
			face := make([]int, 0, len(fields)-1)
			for _, f := range fields[1:] {
				// Split by "/" and take only the first part (vertex index)
				parts := strings.Split(f, "/")
				idx, err := strconv.Atoi(parts[0])
				if err != nil {
					return nil, fmt.Errorf("line %d: invalid vertex index: %w", lineNum, err)
				}
				// OBJ indices are 1-based, convert to 0-based
				face = append(face, idx-1)
			}
			model.Faces = append(model.Faces, face)
		}
		// Ignore other lines (vt, vn, etc.)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return model, nil
}

// Center moves the model so its center is at the origin
func (m *ObjModel) Center() {
	if len(m.Points) == 0 {
		return
	}

	var sumX, sumY, sumZ float64
	for _, p := range m.Points {
		sumX += p.x
		sumY += p.y
		sumZ += p.z
	}

	n := float64(len(m.Points))
	centerX, centerY, centerZ := sumX/n, sumY/n, sumZ/n

	for i := range m.Points {
		m.Points[i].x -= centerX
		m.Points[i].y -= centerY
		m.Points[i].z -= centerZ
	}
}

// Scale multiplies all coordinates by the given factor
func (m *ObjModel) Scale(factor float64) {
	for i := range m.Points {
		m.Points[i].x *= factor
		m.Points[i].y *= factor
		m.Points[i].z *= factor
	}
}

// Normalize scales the model to fit within a unit cube (-0.5 to 0.5)
func (m *ObjModel) Normalize() {
	if len(m.Points) == 0 {
		return
	}

	// Find bounding box
	minX, minY, minZ := m.Points[0].x, m.Points[0].y, m.Points[0].z
	maxX, maxY, maxZ := minX, minY, minZ

	for _, p := range m.Points {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
		if p.z < minZ {
			minZ = p.z
		}
		if p.z > maxZ {
			maxZ = p.z
		}
	}

	// Find the largest dimension
	dx, dy, dz := maxX-minX, maxY-minY, maxZ-minZ
	maxDim := dx
	if dy > maxDim {
		maxDim = dy
	}
	if dz > maxDim {
		maxDim = dz
	}

	if maxDim == 0 {
		return
	}

	// Center and scale
	m.Center()
	m.Scale(1.0 / maxDim)
}
