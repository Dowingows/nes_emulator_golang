package main

import (
	"fmt"
)

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
	instructions.add(Instruction{0x00, "BRK", func(cpu *CPU) {

	}})
	instructions.add(Instruction{0x01, "ORA", func(cpu *CPU) {

	}})

	//LDX

	instructions.add(Instruction{0xA2, "LDX", func(cpu *CPU) {
		addr := cpu.immediateAddress()
		cpu.Ldx(addr)
		fmt.Printf(" #$%02X", cpu.memory.fetch(addr))
	}})
	instructions.add(Instruction{0x05, "ORA", func(cpu *CPU) {

	}})
	instructions.add(Instruction{0x06, "ASL", func(cpu *CPU) {

	}})
	instructions.add(Instruction{0x08, "PHP", func(cpu *CPU) {

	}})
	instructions.add(Instruction{0x09, "ORA", func(cpu *CPU) {

	}})
	instructions.add(Instruction{0x0A, "ASL", func(cpu *CPU) {

	}})
	instructions.add(Instruction{0x0B, "0D", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0xEA, "NOP", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0x20, "JSR", func(cpu *CPU) {
		addr := cpu.absoluteAddress()
		cpu.Jsr(addr)
		fmt.Printf(" $%04X", addr)
	}})

	instructions.add(Instruction{0xA9, "LDA", func(cpu *CPU) {
		//this is immediate addressing
		addr := cpu.immediateAddress()
		cpu.Lda(addr)
		fmt.Printf(" #$%02X", cpu.memory.fetch(addr))
	}})
	instructions.add(Instruction{0xA5, "LDA", func(cpu *CPU) {
		//this is zero page addressing
		addr := cpu.zeroPageAddress()
		cpu.Lda(addr)
		fmt.Printf(" $%02X", addr)
	}})
	instructions.add(Instruction{0xB5, "LDA", func(cpu *CPU) {

	}})
	instructions.add(Instruction{0xAD, "LDA", func(cpu *CPU) {

	}})
	instructions.add(Instruction{0xBD, "LDA", func(cpu *CPU) {

	}})
	instructions.add(Instruction{0xB9, "LDA", func(cpu *CPU) {

	}})
	instructions.add(Instruction{0xA1, "LDA", func(cpu *CPU) {

	}})
	instructions.add(Instruction{0xB1, "LDA", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0xA0, "LDY", func(cpu *CPU) {
		//this is immediate addressing
		addr := cpu.immediateAddress()
		cpu.Ldy(addr)
		fmt.Printf(" #$%02X", cpu.memory.fetch(addr))
	}})

	instructions.add(Instruction{0x8C, "STY", func(cpu *CPU) {
		//Absolute addressing. We must use two bytes to operand
		addr := cpu.absoluteAddress()
		cpu.store(addr, cpu.registers.Y)
		fmt.Printf(" $%04X", addr)
	}})

	instructions.add(Instruction{0x60, "RTS", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0x48, "PHA", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0x81, "STA", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0x85, "STA", func(cpu *CPU) {
		// Zero page addressing.
		addr := cpu.zeroPageAddress()
		cpu.store(addr, cpu.registers.A)
		fmt.Printf(" $%02X", addr)
	}})

	instructions.add(Instruction{0x8D, "STA", func(cpu *CPU) {
		//Absolute addressing. We must use two bytes to operand
		addr := cpu.absoluteAddress()
		cpu.store(addr, cpu.registers.A)
		fmt.Printf(" $%04X", addr)
	}})

	instructions.add(Instruction{0x91, "STA", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0x95, "STA", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0x99, "STA", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0x9D, "STA", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0x68, "PLA", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0x30, "BMI", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0x4C, "JMP", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0x6C, "JMP", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0x10, "BPL", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0x78, "SEI", func(cpu *CPU) {

	}})

	instructions.add(Instruction{0x18, "CLC", func(cpu *CPU) {
		fmt.Printf("\nCLC\n")
	}})

	instructions.add(Instruction{0x6C, "JMP", func(cpu *CPU) {
		fmt.Printf("\nJMP\n")
	}})

	return instructions
}
