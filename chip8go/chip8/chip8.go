package chip8

import (
	"fmt"
	"log"
	"math/rand"
	"os"
)

type Chip8 struct {
	opcode      uint16
	memory      [4096]byte
	V           [16]byte
	I           uint16
	pc          uint16
	gfx         [64 * 32]byte
	delay_timer byte
	sound_timer byte
	stack       [16]uint16
	sp          uint16
	key         [16]byte
	drawFlag    bool
}

var fontset []uint8 = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func New() *Chip8 {
	return &Chip8{
		opcode: 0,
		pc:     0x200,
		I:      0,
		sp:     0,
	}
}

func (chip *Chip8) Initialize() {
	// Clear display
	for i := range chip.gfx {
		chip.gfx[i] = 0
	}
	// Clear stack
	for i := range chip.stack {
		chip.stack[i] = 0
	}
	// Clear registers V0-VF
	for i := range chip.V {
		chip.V[i] = 0
	}
	// Clear memory
	for i := range chip.memory {
		chip.memory[i] = 0
	}
	// Load fontset
	for i := range 80 {
		chip.memory[i] = fontset[i]

	}

	// Reset timers
	chip.delay_timer = 0
	chip.sound_timer = 0

	// draw clear screen at initialization
	chip.drawFlag = true
}

func (chip *Chip8) LoadProgram(path string) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		log.Panicf("an error ocurred when reading the source code of the provided path %v", err)
	}

	for i := range len(buffer) {
		chip.memory[i+512] = buffer[i]
	}
}

func (chip *Chip8) EmulateCycle() {
	chip.opcode = uint16(chip.memory[chip.pc])<<8 | uint16(chip.memory[chip.pc+1])
	switch chip.opcode & 0xF000 {
	// perform opcode translation here
	case 0xA000:
		chip.I = chip.opcode & 0x0FFF
		chip.pc += 2
	case 0x2000:
		chip.stack[chip.sp] = chip.pc
		chip.sp++
		chip.pc = chip.opcode & 0x0FFF
	case 0x1000:
		chip.pc = chip.opcode & 0x0FFF
		chip.pc += 2
	case 0x3000:
		if uint16(chip.V[(chip.opcode&0x0F00>>8)]) == (chip.opcode & 0x00FF) {
			chip.pc += 4
		} else {
			chip.pc += 2
		}
	case 0x4000:
		if uint16(chip.V[(chip.opcode&0x0F00>>8)]) == (chip.opcode & 0x00FF) {
			chip.pc += 2
		} else {
			chip.pc += 4
		}
	case 0x5000:
		if chip.V[(chip.opcode&0x0F00>>8)] == chip.V[(chip.opcode&0x00F0>>4)] {
			chip.pc += 4
		} else {
			chip.pc += 2
		}
	case 0x6000:
		chip.V[(chip.opcode & 0x0F00 >> 8)] = byte(chip.opcode & 0x00FF)
		chip.pc += 2
	case 0x7000:
		chip.V[(chip.opcode & 0x0F00 >> 8)] += byte(chip.opcode & 0x00FF)
		chip.pc += 2
	case 0x0000:
		switch chip.opcode & 0x000F {
		case 0x0000:
			// clear screen
			for i := range chip.gfx {
				chip.gfx[i] = 0
			}
			chip.drawFlag = true
			chip.pc += 2
		case 0x000E:
			// return from subroutine
			chip.sp--
			chip.pc = chip.stack[chip.sp] + 2
		default:
			fmt.Printf("unkown opcode 0x%X\n", chip.opcode)
		}
		chip.pc += 2
	case 0x8000:
		switch chip.opcode & 0x000F {
		case 0x0000:
			chip.V[(chip.opcode&0x0F00)>>8] = chip.V[(chip.opcode&0x00F0)>>4]
			chip.pc += 2
		case 0x0001:
			chip.V[(chip.opcode&0x0F00)>>8] |= chip.V[(chip.opcode&0x00F0)>>4]
			chip.pc += 2
		case 0x0002:
			chip.V[(chip.opcode&0x0F00)>>8] &= chip.V[(chip.opcode&0x00F0)>>4]
			chip.pc += 2
		case 0x0003:
			chip.V[(chip.opcode&0x0F00)>>8] ^= chip.V[(chip.opcode&0x00F0)>>4]
			chip.pc += 2
		case 0x0004:
			if chip.V[(chip.opcode&0x00F0)>>4] > (0xFF - chip.V[(chip.opcode&0x0F00)>>8]) {
				chip.V[0xF] = 1
			} else {
				chip.V[0xF] = 0
			}
			chip.V[(chip.opcode&0x0F00)>>8] += chip.V[(chip.opcode&0x00F0)>>4]
			chip.pc += 2
		case 0x0005:
			if chip.V[(chip.opcode&0x00F0)>>4] > chip.V[(chip.opcode&0x0F00)>>8] {
				chip.V[0xF] = 0
			} else {
				chip.V[0xF] = 1
			}
			chip.V[(chip.opcode&0x0F00)>>8] -= chip.V[(chip.opcode&0x00F0)>>4]
			chip.pc += 2
		case 0x0006:
			chip.V[0xF] = chip.V[(chip.opcode&0x0F00)>>8] & 0x1
			chip.V[(chip.opcode&0x0F00)>>8] = chip.V[(chip.opcode&0x0F00)>>8] >> 1
			chip.pc += 2
		case 0x0007:
			if chip.V[(chip.opcode&0x00F0)>>4] < chip.V[(chip.opcode&0x0F00)>>8] {
				chip.V[0xF] = 0
			} else {
				chip.V[0xF] = 1
			}
			chip.V[(chip.opcode&0x0F00)>>8] = chip.V[(chip.opcode&0x00F0)>>4] - chip.V[(chip.opcode&0x0F00)>>8]
			chip.pc += 2
		case 0x000E:
			chip.V[0xF] = (chip.V[(chip.opcode&0x0F00)>>8] >> 7) & 0x1
			chip.V[(chip.opcode&0x0F00)>>8] = chip.V[(chip.opcode&0x0F00)>>8] << 1
			chip.pc += 2
		}

	case 0x9000:
		if chip.V[(chip.opcode&0x00F0)>>4] != chip.V[(chip.opcode&0x0F00)>>8] {
			chip.pc += 4
		} else {
			chip.pc += 2
		}
	case 0xB000:
		chip.pc = uint16(chip.V[0x0]) + (chip.opcode & 0x0FFF)
	case 0xC000:
		chip.V[(chip.opcode&0x00F0)>>8] = byte(rand.Intn(256)) & byte((chip.opcode & 0x00FF))
		chip.pc += 2
	case 0xF000:
		switch chip.opcode & 0x00FF {
		case 0x0033:
			chip.memory[chip.I] = chip.V[(chip.opcode&0x0F00)>>8] / 100
			chip.memory[chip.I+1] = (chip.V[(chip.opcode&0x0F00)>>8] / 10) % 10
			chip.memory[chip.I+2] = (chip.V[(chip.opcode&0x0F00)>>8] % 100) % 10
			chip.pc += 2
		case 0x0007:
			chip.V[(chip.opcode&0x0F00)>>8] = chip.delay_timer
			chip.pc += 2
		case 0x000A:
			keyPressed := false
			for i := range chip.key {
				if chip.key[i] != 0 {
					chip.V[(chip.opcode&0x0F00)>>8] = chip.key[i]
					keyPressed = true
				}
			}
			if !keyPressed {
				return
			}
			chip.pc += 2
		case 0x0015:
			chip.delay_timer = chip.V[(chip.opcode&0x0F00)>>8]
			chip.pc += 2
		case 0x0018:
			chip.sound_timer = chip.V[(chip.opcode&0x0F00)>>8]
			chip.pc += 2
		case 0x001E:
			chip.I += uint16(chip.V[(chip.opcode&0x0F00)>>8])
			chip.pc += 2
		case 0x0029:
			chip.I = uint16(chip.V[(chip.opcode&0x0F00)>>8] * 0x5)
			chip.pc += 2
		case 0x0055:
			for i := uint16(0); i <= (chip.opcode&0x0F00)>>8; i++ {
				chip.memory[chip.I+i] = chip.V[i]
			}
			chip.pc += 2
		case 0x0065:
			for i := uint16(0); i <= (chip.opcode&0x0F00)>>8; i++ {
				chip.V[i] = chip.memory[chip.I+i]
			}
			chip.pc += 2
		}

	case 0xD000:
		x := uint16(chip.V[(chip.opcode&0x0F00)>>8])
		y := uint16(chip.V[(chip.opcode&0x00F0)>>4])
		height := chip.opcode & 0x000F
		var pixel uint16
		chip.V[0xF] = 0
		for yline := range height {
			pixel = uint16(chip.memory[chip.I+yline])
			for xline := range uint16(8) {
				if pixel&(0x80>>xline) != 0 {
					if chip.gfx[(x+xline+((y+yline)*64))] == 1 {
						chip.V[0xF] = 1
					}
					chip.gfx[(x + xline + ((y + yline) * 64))] ^= 1
				}
			}
		}
		chip.drawFlag = true
		chip.pc += 2
	case 0xE000:
		switch chip.opcode & 0x00FF {
		case 0x009E:
			if chip.key[chip.V[(chip.opcode&0x0F00)>>8]] != 0 {
				chip.pc += 4
			} else {
				chip.pc += 2
			}
		case 0x00A1:
			if chip.key[chip.V[(chip.opcode&0x0F00)>>8]] == 0 {
				chip.pc += 4
			} else {
				chip.pc += 2
			}
		}

	default:
		fmt.Printf("unkown opcode 0x%X\n", chip.opcode)
	}

	if chip.delay_timer > 0 {
		chip.delay_timer--
	}
	if chip.sound_timer > 0 {
		if chip.sound_timer == 1 {
			fmt.Printf("SIMULATING SOUND: BEEP\n")
		}
		chip.sound_timer--
	}

}
