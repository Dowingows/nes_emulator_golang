package main

import (
	"fmt"
	"io/ioutil"
)

//CPU of our NES emulator
type CPU struct {
	registers         Registers
	memory            Memory
	instructionsTable InstructionsTable
}

func (cpu *CPU) init(memory Memory) {
	cpu.registers.init()
	cpu.memory = memory
	cpu.instructionsTable = newInstructionsTable()
}

func (cpu *CPU) printRegisters() {
	fmt.Printf("\n|A: %02x X: %02x Y: %02x P: %02x SP: %02x PC: %04x|\n", cpu.registers.A, cpu.registers.X, cpu.registers.Y, cpu.registers.P, cpu.registers.SP, cpu.registers.PC)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (cpu *CPU) load(addr uint16, register *uint8) {
	*register = cpu.memory.fetch(addr)
}

func (cpu *CPU) store(addr uint16, value uint8) {
	cpu.memory.store(addr, value)
}

func (cpu *CPU) execute(opcode byte) {
	instr := cpu.instructionsTable.get(opcode)
	if instr.fetch == nil {
		fmt.Printf("\nInstrução não existente no processador: {%04x}!\n", opcode)
	} else {
		fmt.Println("\nExecutando instrução... ")
		fmt.Print(instr)
		instr.fetch(cpu)
	}
}

///Adressing modes

func (cpu *CPU) solveTypeAddress(opcode uint8) uint16 {

	addr := uint16(0)

	switch instrsMode[opcode] {
	case modeAbsolute:
		addr = cpu.absoluteAddress()
		break
	case modeImmediate:
		//this is immediate addressing
		addr = cpu.immediateAddress()
		break
	case modeZeroPage:
		addr = cpu.zeroPageAddress()
		break
	default:
		fmt.Printf("\n* * * * Não encontrado esse MODE!! %d * * * * \n", instrsMode[opcode])
		return 0
	}
	return addr
}

func (cpu *CPU) immediateAddress() uint16 {
	addr := cpu.registers.PC
	cpu.registers.PC++
	return addr
}

func (cpu *CPU) zeroPageAddress() uint16 {

	addr := uint16(cpu.memory.fetch(cpu.registers.PC))
	cpu.registers.PC++
	return addr
}

func (cpu *CPU) absoluteAddress() uint16 {
	//Absolute addressing. We must use two bytes to operand
	addr := uint16(cpu.memory.fetch(cpu.registers.PC+1))<<8 | uint16(cpu.memory.fetch(cpu.registers.PC))
	cpu.registers.PC += 2
	return addr
}

func (cpu *CPU) run() {
	i := 0
	stopIn := 22
	fmt.Println()
	for {
		i++
		opcode := cpu.memory[cpu.registers.PC]
		instr := cpu.instructionsTable.get(opcode)
		fmt.Printf(" [%2x]: %s ", opcode, instr.mneumonic)
		cpu.registers.PC++

		if instr.fetch == nil {
			break
		}

		instr.fetch(cpu)

		if instr.opcode == 0x00 || i >= stopIn {
			break
		}
		fmt.Println()
	}
	fmt.Printf("\n* Acabou!! *\n")
}

//Stack (Location: 0x100 -> 0x200)

func (cpu *CPU) push(value uint8) {
	cpu.memory.store(0x0100|uint16(cpu.registers.SP), value)
	cpu.registers.SP--
}

func (cpu *CPU) pull() uint8 {
	cpu.registers.SP++
	value := cpu.memory.fetch(0x0100 | uint16(cpu.registers.SP))
	return value
}

func (cpu *CPU) push16(value uint16) {
	cpu.push(uint8(value >> 8))
	cpu.push(uint8(value))
}

func (cpu *CPU) pull16() (value uint16) {
	low := cpu.pull()
	high := cpu.pull()
	value = (uint16(high) << 8) | uint16(low)
	return
}

//Instructions

//Lda Load value from a memory addr to A register
func (cpu *CPU) Lda(addr uint16) {
	cpu.load(addr, &cpu.registers.A)
}

//Ldx Load value from a memory addr to A register
func (cpu *CPU) Ldx(addr uint16) {
	cpu.load(addr, &cpu.registers.X)
}

//Ldy Load value from a memory addr to Y register
func (cpu *CPU) Ldy(addr uint16) {
	cpu.load(addr, &cpu.registers.Y)
}

//Sta Store  A register in a memory address
func (cpu *CPU) Sta(addr uint16) {
	cpu.memory.store(addr, cpu.registers.A)
}

/*Jsr pushes the address-1 of the next operation on to the stack
 *	before transferring program control to the following address
 */
func (cpu *CPU) Jsr(address uint16) {
	value := cpu.registers.PC - 1
	cpu.push16(value)
	//Descomentar depois!!! cpu.registers.PC = address
}

func main() {

	var m Memory
	cpu := CPU{}
	cpu.init(m)

	/*for _, value := range cpu.instructionsTable {
		fmt.Printf(value.toString())
	}*/

	//instr := cpu.instructionsTable.get(0)
	//instr.fetch(&cpu)
	//fmt.Println(cpu.instructionsTable.get(0))

	path := "tests/cpu_test.asm"
	data, _ := ioutil.ReadFile(path)
	cpu.printRegisters()
	cpu.memory.loadCode(data, cpu.registers.PC)
	cpu.run()
	cpu.printRegisters()
	/*if string(data[:3]) != "NES" {
		log.Fatalf("Invalid ROM file" + string(data[:3]))
	}*/

	//path := "roms/galaga.nes"
	//path := "tests/cpu_dummy_reads.nes"
	/*path := "tests/nestest.nes"
	data, _ := ioutil.ReadFile(path)

	if string(data[:3]) != "NES" {
		log.Fatalf("Invalid ROM file" + string(data[:3]))
	}
	cpu.run(data)*/
	/*NPRGROMBanks := data[4]

	prgBeginning := uint64(16)
	prgEnd := 16 + uint64(NPRGROMBanks)*0x4000

	PRGROM := data[prgBeginning:prgEnd]

	opcode := cpu.memory.fetch(cpu.registers.PC)

	fmt.Println(PRGROM)
	fmt.Printf("\n (%x) = %x \n", cpu.registers.PC, opcode)*/
	/*fmt.Println()
	for i := 0; i < len(data); i++ {
		fmt.Printf("\n%02x ", data[i])
	}
	fmt.Println()*/
}
