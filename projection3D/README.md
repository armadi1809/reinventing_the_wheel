# 3D Projection Demo

Minimal graphics programming project that loads a Wavefront `.obj` mesh and renders it as a rotating wireframe using the Ebiten 2D game library for Go.

## Overview

- Description: Loads an `.obj` file (vertices + faces), normalizes the model, applies a simple perspective projection, and draws face edges as lines in a 640×480 window.

## Prerequisites

- Go (1.20+ recommended).
- Ebiten will be fetched automatically from the module when you run or build.

## Run

From the repository root or the `projection3D` directory, pass a single `.obj` file path:

```bash
cd projection3D
go run . data/teapot.obj
```

Other sample models included:

```bash
go run . data/penger.obj
go run . data/real-penger.obj
```

## Build

```bash
cd projection3D
go build -o projection3d-app
./projection3d-app data/teapot.obj
```

## Controls / Usage notes

- Rotate:
  - `←` rotate left
  - `→` rotate right
- Window size: 640×480.
- Rendering: wireframe (edges of each face) drawn in green.

### OBJ support

The loader supports:

- `v x y z` vertex lines
- `f ...` face lines (uses the vertex index; supports `v/vt/vn` style tokens)

Other OBJ record types (`vt`, `vn`, etc.) are ignored.

## Troubleshooting

- If dependencies are missing or build fails, run:

```bash
cd projection3D
go mod tidy
```

## Files

- `projection3D/main.go` — Ebiten game loop + projection + wireframe rendering.
- `projection3D/fileProcessor.go` — minimal OBJ loader + normalization helpers.
- `projection3D/data/` — sample `.obj` meshes.
