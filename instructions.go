package main

import (
	"fmt"
)

// addressing modes
const (
	_ = iota
	modeAbsolute
	modeAbsoluteX
	modeAbsoluteY
	modeAccumulator
	modeImmediate
	modeImplied
	modeIndexedIndirect
	modeIndirect
	modeIndirectIndexed
	modeRelative
	modeZeroPage
	modeZeroPageX
	modeZeroPageY
)

var instrsMode = [256]byte{
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	1, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 8, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 13, 13, 6, 3, 6, 3, 2, 2, 3, 3,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 13, 13, 6, 3, 6, 3, 2, 2, 3, 3,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
}

//Instruction is a Struct to save opcode, mnemonic and func
type Instruction struct {
	opcode    uint8
	mneumonic string
	fetch     func(*CPU)
}

//InstructionsTable is a table with valid instructions to 6502 processor
type InstructionsTable map[uint8]Instruction

func (instruction *Instruction) toString() string {
	return fmt.Sprintf("\n [%02x] %s ", instruction.opcode, instruction.mneumonic)
}

func (instructionsTable InstructionsTable) add(instr Instruction) {
	instructionsTable[instr.opcode] = instr
}

//opcode is a hexadecimal value. First, we need to cast hex to int, after we get the value from map
func (instructionsTable InstructionsTable) get(opcode byte) Instruction {
	return instructionsTable[uint8(opcode)]
}

func newInstructionsTable() InstructionsTable {
	instructions := InstructionsTable{}

	//JSR
	for _, val := range []uint8{0x20} {
		opcode := val
		instructions.add(Instruction{opcode, "JSR", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.Jsr(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//LDA

	for _, val := range []uint8{0xA5, 0xA9} {
		opcode := val
		instructions.add(Instruction{opcode, "LDA", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.Lda(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//LDX

	for _, val := range []uint8{0xA2} {
		opcode := val
		instructions.add(Instruction{opcode, "LDX", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.Ldx(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//LDY
	for _, val := range []uint8{0xA0} {
		opcode := val
		instructions.add(Instruction{opcode, "LDY", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.Ldy(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//STY
	for _, val := range []uint8{0x8C} {
		opcode := val
		instructions.add(Instruction{opcode, "STY", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.store(addr, cpu.registers.Y)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}
	//STA
	for _, val := range []uint8{0x85, 0x8D} {
		opcode := val
		instructions.add(Instruction{opcode, "STA", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.store(addr, cpu.registers.A)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//LSR
	for _, val := range []uint8{0x4A} {
		opcode := val
		instructions.add(Instruction{opcode, "LSR", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.store(addr, cpu.registers.A)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	return instructions
}
