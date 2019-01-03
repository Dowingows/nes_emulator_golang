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

	//RTS
	for _, val := range []uint8{0x60} {
		opcode := val
		instructions.add(Instruction{opcode, "RTS", func(cpu *CPU) {
			cpu.Rts()
		}})
	}

	//RTI
	for _, val := range []uint8{0x40} {
		opcode := val
		instructions.add(Instruction{opcode, "RTI", func(cpu *CPU) {
			cpu.Rti()
		}})
	}

	//LDA

	for _, val := range []uint8{0xA5, 0xA9, 0XAD} {
		opcode := val
		instructions.add(Instruction{opcode, "LDA", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.Lda(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//LDX

	for _, val := range []uint8{0xA2, 0XAE} {
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
	for _, val := range []uint8{0x86, 0x8E} {
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
				/*addr := cpu.solveTypeAddress(opcode)
				cpu.store(addr, cpu.registers.A)
				fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))*/
			}

		}})
	}

	//ASL
	for _, val := range []uint8{0x0A} {
		opcode := val
		instructions.add(Instruction{opcode, "ASL", func(cpu *CPU) {
			if instrsMode[opcode] == modeAccumulator {
				cpu.AslA()
				fmt.Printf(" | [A:] %04x |", cpu.registers.A)
			} else {
				fmt.Printf("\n Não implementada\n")
			}

		}})
	}

	//ROR
	for _, val := range []uint8{0x6A} {
		opcode := val
		instructions.add(Instruction{opcode, "ROR", func(cpu *CPU) {
			if instrsMode[opcode] == modeAccumulator {
				cpu.RorA()
				fmt.Printf(" | [A:] %04x |", cpu.registers.A)
			} else {
				fmt.Printf("\n Não implementada\n")
			}

		}})
	}

	//ROR
	for _, val := range []uint8{0x2A} {
		opcode := val
		instructions.add(Instruction{opcode, "ROR", func(cpu *CPU) {
			if instrsMode[opcode] == modeAccumulator {
				cpu.RolA()
				fmt.Printf(" | [A:] %04x |", cpu.registers.A)
			} else {
				fmt.Printf("\n Não implementada\n")
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
			addr := cpu.relativeAddress()
			cpu.Bvs(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//BVC
	for _, val := range []uint8{0x50} {
		opcode := val
		instructions.add(Instruction{opcode, "BVC", func(cpu *CPU) {
			addr := cpu.relativeAddress()
			cpu.Bvc(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//BPL
	for _, val := range []uint8{0x10} {
		opcode := val
		instructions.add(Instruction{opcode, "BPL", func(cpu *CPU) {
			addr := cpu.relativeAddress()
			cpu.Bpl(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//BMI
	for _, val := range []uint8{0x30} {
		opcode := val
		instructions.add(Instruction{opcode, "BMI", func(cpu *CPU) {
			addr := cpu.relativeAddress()
			cpu.Bmi(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
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
			addr := cpu.solveTypeAddress(opcode)
			cpu.Cmp(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//CPX
	for _, val := range []uint8{0xE0} {
		opcode := val
		instructions.add(Instruction{opcode, "CPX", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.Cpx(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//Cpy
	for _, val := range []uint8{0xC0} {
		opcode := val
		instructions.add(Instruction{opcode, "CPY", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.Cpy(addr)
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

	//SED
	for _, val := range []uint8{0xF8} {
		opcode := val
		instructions.add(Instruction{opcode, "SED", func(cpu *CPU) {
			cpu.Sed()
		}})
	}

	//PHP
	for _, val := range []uint8{0x08} {
		opcode := val
		instructions.add(Instruction{opcode, "PHP", func(cpu *CPU) {
			cpu.Php()
		}})
	}

	//PHA
	for _, val := range []uint8{0x48} {
		opcode := val
		instructions.add(Instruction{opcode, "PHA", func(cpu *CPU) {
			cpu.Pha()
		}})
	}

	//PLA
	for _, val := range []uint8{0x68} {
		opcode := val
		instructions.add(Instruction{opcode, "PLA", func(cpu *CPU) {
			cpu.Pla()
		}})
	}

	//PLP
	for _, val := range []uint8{0x28} {
		opcode := val
		instructions.add(Instruction{opcode, "PLP", func(cpu *CPU) {
			cpu.Plp()
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

	//AND
	for _, val := range []uint8{0x29} {
		opcode := val
		instructions.add(Instruction{opcode, "AND", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.And(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//ORA
	for _, val := range []uint8{0x09} {
		opcode := val
		instructions.add(Instruction{opcode, "ORA", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.Ora(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//EOR
	for _, val := range []uint8{0x49} {
		opcode := val
		instructions.add(Instruction{opcode, "EOR", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.Eor(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
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

	//CLD
	for _, val := range []uint8{0xD8} {
		opcode := val
		instructions.add(Instruction{opcode, "CLD", func(cpu *CPU) {
			cpu.Cld()
		}})
	}

	//CLV
	for _, val := range []uint8{0xB8} {
		opcode := val
		instructions.add(Instruction{opcode, "CLV", func(cpu *CPU) {
			cpu.Clv()
		}})
	}

	//ADC
	for _, val := range []uint8{0x69} {
		opcode := val
		instructions.add(Instruction{opcode, "ADC", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.Adc(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}
	//SBC
	for _, val := range []uint8{0xE9} {
		opcode := val
		instructions.add(Instruction{opcode, "SBC", func(cpu *CPU) {
			addr := cpu.solveTypeAddress(opcode)
			cpu.Sbc(addr)
			fmt.Printf(" | [%04X] %04x |", addr, cpu.memory.fetch(addr))
		}})
	}

	//DEY
	for _, val := range []uint8{0x88} {
		opcode := val
		instructions.add(Instruction{opcode, "DEY", func(cpu *CPU) {
			cpu.Dey()
		}})
	}

	//DEX
	for _, val := range []uint8{0xCA} {
		opcode := val
		instructions.add(Instruction{opcode, "DEX", func(cpu *CPU) {
			cpu.Dex()
		}})
	}

	//INY
	for _, val := range []uint8{0xC8} {
		opcode := val
		instructions.add(Instruction{opcode, "INY", func(cpu *CPU) {
			cpu.Iny()
		}})
	}

	//INX
	for _, val := range []uint8{0xE8} {
		opcode := val
		instructions.add(Instruction{opcode, "INX", func(cpu *CPU) {
			cpu.Inx()
		}})
	}

	//Tay
	for _, val := range []uint8{0xA8} {
		opcode := val
		instructions.add(Instruction{opcode, "TAY", func(cpu *CPU) {
			cpu.Tay()
		}})
	}

	//Tax
	for _, val := range []uint8{0xAA} {
		opcode := val
		instructions.add(Instruction{opcode, "TAX", func(cpu *CPU) {
			cpu.Tax()
		}})
	}

	//Tya
	for _, val := range []uint8{0x98} {
		opcode := val
		instructions.add(Instruction{opcode, "TYA", func(cpu *CPU) {
			cpu.Tya()
		}})
	}

	//Txa
	for _, val := range []uint8{0x8A} {
		opcode := val
		instructions.add(Instruction{opcode, "TXA", func(cpu *CPU) {
			cpu.Txa()
		}})
	}

	//Tsx
	for _, val := range []uint8{0xBA} {
		opcode := val
		instructions.add(Instruction{opcode, "TSX", func(cpu *CPU) {
			cpu.Tsx()
		}})
	}

	return instructions
}
