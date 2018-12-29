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

	//STA
	for _, val := range []uint8{0x85, 0x8D} {
		opcode := val
		instructions.add(Instruction{opcode, "STA", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.store(addr, cpu.registers.A)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//STX
	for _, val := range []uint8{0x86} {
		opcode := val
		instructions.add(Instruction{opcode, "STX", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.store(addr, cpu.registers.X)
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

	//LSR
	for _, val := range []uint8{0x4A} {
		opcode := val
		instructions.add(Instruction{opcode, "LSR", func(cpu *CPU) {
			if instrsMode[opcode] == modeAccumulator {
				cpu.LsrA()
				fmt.Printf(" | [A:] %04x |", cpu.registers.A)
			} else {
				addr := cpu.solveTypeAddress(opcode)
				cpu.store(addr, cpu.registers.A)
				fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
			}

		}})
	}

	//BCS
	for _, val := range []uint8{0xB0} {
		opcode := val
		instructions.add(Instruction{opcode, "BCS", func(cpu *CPU) {
			addr := cpu.relativeAddress()
			cpu.Bcs(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//BCC
	for _, val := range []uint8{0x90} {
		opcode := val
		instructions.add(Instruction{opcode, "BCC", func(cpu *CPU) {
			addr := cpu.relativeAddress()
			cpu.Bcc(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//BEQ
	for _, val := range []uint8{0xF0} {
		opcode := val
		instructions.add(Instruction{opcode, "BEQ", func(cpu *CPU) {
			addr := cpu.relativeAddress()
			cpu.Beq(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//BIT
	for _, val := range []uint8{0x24} {
		opcode := val
		instructions.add(Instruction{opcode, "BIT", func(cpu *CPU) {
			addr := cpu.zeroPageAddress()
			cpu.Bit(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//BVS
	for _, val := range []uint8{0x70} {
		opcode := val
		instructions.add(Instruction{opcode, "BVS", func(cpu *CPU) {

		}})
	}

	//BPL
	for _, val := range []uint8{0x10} {
		opcode := val
		instructions.add(Instruction{opcode, "BPL", func(cpu *CPU) {
			addr := cpu.registers.PC
			cpu.registers.PC += 1
			//addr := cpu.immediateAddress()
			//cpu.Bcs(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr+2))
		}})
	}

	//BNE
	for _, val := range []uint8{0xD0} {
		opcode := val
		instructions.add(Instruction{opcode, "BNE", func(cpu *CPU) {
			addr := cpu.relativeAddress()
			cpu.Bne(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//JMP
	for _, val := range []uint8{0x4C} {
		opcode := val
		instructions.add(Instruction{opcode, "JMP", func(cpu *CPU) {
			//oldValue := cpu.registers.PC //Remover
			addr := cpu.solveTypeAddress(opcode)
			cpu.Jmp(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
			//cpu.registers.PC = oldValue + 2 //Remover (Coloquei pq ele dá jump)
		}})
	}

	//DEC
	for _, val := range []uint8{0xC6} {
		opcode := val
		instructions.add(Instruction{opcode, "DEC", func(cpu *CPU) {

			addr := cpu.solveTypeAddress(opcode)
			cpu.Dec(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))

		}})
	}

	//INC
	for _, val := range []uint8{0xE6} {
		opcode := val
		instructions.add(Instruction{opcode, "INC", func(cpu *CPU) {

			addr := cpu.solveTypeAddress(opcode)
			cpu.Inc(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))

		}})
	}

	//CMP
	for _, val := range []uint8{0xC9} {
		opcode := val
		instructions.add(Instruction{opcode, "CMP", func(cpu *CPU) {
			/*FALTANDO IMPLEMENTAR*/
			addr := cpu.solveTypeAddress(opcode)

			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//SEI
	for _, val := range []uint8{0x78} {
		opcode := val
		instructions.add(Instruction{opcode, "SEI", func(cpu *CPU) {
			cpu.Sei()
		}})
	}

	//TXS
	for _, val := range []uint8{0x9A} {
		opcode := val
		instructions.add(Instruction{opcode, "TXS", func(cpu *CPU) {
			cpu.Txs()
		}})
	}

	//NOP
	for _, val := range []uint8{0xEA} {
		opcode := val
		instructions.add(Instruction{opcode, "NOP", func(cpu *CPU) {
			cpu.Nop()
		}})
	}

	//SEC
	for _, val := range []uint8{0x38} {
		opcode := val
		instructions.add(Instruction{opcode, "SEC", func(cpu *CPU) {
			cpu.Sec()
		}})
	}

	//CLC
	for _, val := range []uint8{0x18} {
		opcode := val
		instructions.add(Instruction{opcode, "CLC", func(cpu *CPU) {
			cpu.Clc()
		}})
	}

	return instructions
}
