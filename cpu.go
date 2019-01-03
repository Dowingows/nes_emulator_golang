package main

import (
	"fmt"
	"io/ioutil"
	"log"
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
	value := cpu.memory.fetch(addr)
	cpu.setZNFlags(value)
	*register = value
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

func (cpu *CPU) relativeAddress() (addr uint16) {
	value := uint16(cpu.memory.fetch(cpu.registers.PC))
	cpu.registers.PC++

	var offset uint16

	if value > 0x7f {
		offset = -(0x0100 - value)
	} else {
		offset = value
	}

	addr = cpu.registers.PC + offset

	return
}

func (cpu *CPU) run() {

	i := 0
	//stopIn := 386
	fmt.Println()
	for {
		//cpu.printRegisters()
		i++
		opcode := cpu.memory[cpu.registers.PC]
		instr := cpu.instructionsTable.get(opcode)
		fmt.Printf("%d - [%2x]: %s ", i, opcode, instr.mneumonic)
		cpu.registers.PC++

		if instr.fetch == nil {
			fmt.Println("\nInstrução não encontrada!! Abortar <{Opcode: 0x%02x}>", opcode)
			break
		}

		instr.fetch(cpu)
		fmt.Printf("   <{PC: %04x | A:%02x | X: %02x | Y: %02x | SP: %02x | P: %08b (%02x) }> ", cpu.registers.PC, cpu.registers.A, cpu.registers.X, cpu.registers.Y, cpu.registers.SP, cpu.registers.P, cpu.registers.P)
		if instr.opcode == 0x00 /*|| i >= stopIn */ {
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

func (cpu *CPU) Txs() {
	cpu.registers.SP = cpu.registers.X
}

//Ldy Load value from a memory addr to Y register
func (cpu *CPU) Ldy(addr uint16) {
	cpu.load(addr, &cpu.registers.Y)
}

//Sta Store  A register in a memory address
func (cpu *CPU) Sta(addr uint16) {
	cpu.memory.store(addr, cpu.registers.A)
}

//Stx Store  X register in a memory address
func (cpu *CPU) Stx(addr uint16) {
	cpu.memory.store(addr, cpu.registers.X)
}

//Nop doenst anything
func (cpu *CPU) Nop() {

}

func (cpu *CPU) Sei() {
	cpu.registers.P = setBit(cpu.registers.P, I)
}

//Sec sets to one the 'C' flag the register P
func (cpu *CPU) Sec() {
	cpu.registers.P = setBit(cpu.registers.P, C)
}

//Sec sets to one the 'D' flag the register P
func (cpu *CPU) Sed() {
	cpu.registers.P = setBit(cpu.registers.P, D)
}

func (cpu *CPU) Php() {
	// Não sei pq setei o bit B, usei só para ficar igual ao resultado do nestest.log
	cpu.push(uint8(setBit(cpu.registers.P, B)))
}

func (cpu *CPU) Pha() {
	// Não sei pq setei o bit B, usei só para ficar igual ao resultado do nestest.log
	cpu.push(cpu.registers.A)
}

func (cpu *CPU) Pla() {
	cpu.registers.A = cpu.setZNFlags(cpu.pull()) //DÚVIDA!! NÃO SEI SE É Z OU ZN
}

func (cpu *CPU) Plp() {
	value := Status(cpu.pull())
	cpu.registers.P = clearBit(value, B)
	cpu.registers.P = setBit(value, U)
}

// Sets the bit at pos in the integer n.
func setBit(n Status, pos Status) Status {
	n |= (1 << pos)
	return n
}

// Clears the bit at pos in n.
func clearBit(n Status, pos Status) Status {
	mask := ^(1 << uint(pos))
	n &= Status(mask)
	return n
}

//Clc clear the 'C' flag the register P
func (cpu *CPU) Clc() {
	cpu.registers.P = clearBit(cpu.registers.P, C)
}

//Cld clear the 'd' flag the register P
func (cpu *CPU) Cld() {
	cpu.registers.P = clearBit(cpu.registers.P, D)
}

//Clc clear the 'V' flag the register P
func (cpu *CPU) Clv() {
	cpu.registers.P = clearBit(cpu.registers.P, V)
}

/*Jsr pushes the address-1 of the next operation on to the stack
 *	before transferring program control to the following address
 */
func (cpu *CPU) Jsr(address uint16) {
	value := cpu.registers.PC - 1
	cpu.push16(value)
	cpu.registers.PC = address
}

//Rts instruction is used at the end of a subroutine to return to the calling routine. It pulls the program counter (minus one) from the stack.
func (cpu *CPU) Rts() {
	cpu.registers.PC = cpu.pull16() + 1
}

//Rti return from interrupt
func (cpu *CPU) Rti() {
	cpu.registers.P = Status(cpu.pull())
	cpu.registers.P = setBit(cpu.registers.P, U) //uma gambiarra básica, remover depois
	cpu.registers.PC = cpu.pull16()
}

/*Lsr shifts all bits right one position.
 *0 is shifted into bit 7 and the original bit 0 is shifted into the Carry.
 */
func (cpu *CPU) Lsr(addr uint16) {
	value := cpu.memory.fetch(addr)

	oldC := getBit(uint8(value), uint8(C))

	value >>= 1

	if getBit(uint8(value), uint8(N)) == 1 {
		cpu.registers.P = setBit(cpu.registers.P, N)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, N)
	}

	if value == 0 {
		cpu.registers.P = setBit(cpu.registers.P, Z)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, Z)
	}

	if oldC == 1 {
		cpu.registers.P = setBit(cpu.registers.P, C)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, C)
	}

	cpu.memory.store(addr, value)
}

func (cpu *CPU) Bcs(addr uint16) {
	if getBit(uint8(cpu.registers.P), uint8(C)) != 0 {
		cpu.registers.PC = addr
	}
}

func (cpu *CPU) Bcc(addr uint16) {
	if getBit(uint8(cpu.registers.P), uint8(C)) == 0 {
		cpu.registers.PC = addr
	}
}

//Bvs If the overflow flag is set then add the relative displacement to the program counter to cause a branch to a new location.
func (cpu *CPU) Bvs(addr uint16) {
	if getBit(uint8(cpu.registers.P), uint8(V)) != 0 {
		cpu.registers.PC = addr
	}
}

//Bvc If the overflow flag is clear then add the relative displacement to the program counter to cause a branch to a new location.
func (cpu *CPU) Bvc(addr uint16) {
	if getBit(uint8(cpu.registers.P), uint8(V)) == 0 {
		cpu.registers.PC = addr
	}
}

//Bpl If the negative flag is clear then add the relative displacement to the program counter to cause a branch to a new location.
func (cpu *CPU) Bpl(addr uint16) {
	if getBit(uint8(cpu.registers.P), uint8(N)) == 0 {
		cpu.registers.PC = addr
	}
}

//Bne If the zero flag is clear then add the relative displacement to the program counter to cause a branch to a new location.
func (cpu *CPU) Bne(addr uint16) {
	if cpu.registers.P&(Z+1) == 0 {
		cpu.registers.PC = addr
	}
}

//Beq If the zero flag is set then add the relative displacement to the program counter to cause a branch to a new location.
func (cpu *CPU) Beq(addr uint16) {
	if cpu.registers.P&(Z+1) != 0 {
		cpu.registers.PC = addr
	}
}

//Bmi If the negative flag is set then add the relative displacement to the program counter to cause a branch to a new location.
func (cpu *CPU) Bmi(addr uint16) {
	if getBit(uint8(cpu.registers.P), uint8(N)) != 0 {
		cpu.registers.PC = addr
	}
}

func getBit(number uint8, pos uint8) uint8 {
	return ((number >> pos) & 1)
}

//Bit Necessário dá uma refatorada depois para ficar mais clara
func (cpu *CPU) Bit(addr uint16) {
	value := cpu.memory.fetch(addr)
	cpu.setZFlag(value & cpu.registers.A)

	memN := getBit(value, 7)
	memV := getBit(value, 6)
	registerP := cpu.registers.P
	if memN == 1 {
		registerP = setBit(registerP, 7)
	} else {
		registerP = clearBit(registerP, 7)
	}

	if memV == 1 {
		registerP = setBit(registerP, 6)
	} else {
		registerP = clearBit(registerP, 6)
	}
	cpu.registers.P = registerP
}

/*LsrA is a LSR in mode Accumulator
 *
 */
func (cpu *CPU) LsrA() {
	value := cpu.shift(1, cpu.registers.A)
	cpu.setZNFlags(value)
	cpu.registers.A = value
}

//Asla
func (cpu *CPU) AslA() {
	value := cpu.registers.A

	oldN := getBit(uint8(value), uint8(N))

	value <<= 1

	if getBit(uint8(value), uint8(N)) == 1 {
		cpu.registers.P = setBit(cpu.registers.P, N)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, N)
	}

	if value == 0 {
		cpu.registers.P = setBit(cpu.registers.P, Z)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, Z)
	}

	if oldN == 1 {
		cpu.registers.P = setBit(cpu.registers.P, C)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, C)
	}

	cpu.registers.A = value
}

func (cpu *CPU) RorA() {
	value := cpu.registers.A
	oldValue := value
	value >>= 1

	if value == 0x00 {
		if getBit(uint8(oldValue), uint8(C)) == 1 {
			value = byte(setBit(Status(value), N))
		} else {
			value = byte(clearBit(Status(value), N))
		}
	}

	//fmt.Printf("\n Eu estou aqui! %04x\n", value)

	if getBit(uint8(oldValue), uint8(C)) == 1 {
		cpu.registers.P = setBit(Status(cpu.registers.P), C)
	} else {
		cpu.registers.P = clearBit(Status(cpu.registers.P), C)
	}

	cpu.setZNFlags(value)

	cpu.registers.A = value
}

func (cpu *CPU) RolA() {
	value := cpu.registers.A
	oldValue := value
	value <<= 1

	if value == 0x00 {
		if getBit(uint8(oldValue), uint8(N)) == 1 {
			value = byte(setBit(Status(value), C))
		} else {
			value = byte(clearBit(Status(value), C))
		}
	}

	//fmt.Printf("\n Eu estou aqui! %04x\n", value)

	if getBit(uint8(oldValue), uint8(N)) == 1 {
		cpu.registers.P = setBit(Status(cpu.registers.P), C)
	} else {
		cpu.registers.P = clearBit(Status(cpu.registers.P), C)
	}

	cpu.setZNFlags(value)

	cpu.registers.A = value
}

func (cpu *CPU) shift(direction int, value uint8) uint8 {
	c := Status(0)

	switch direction {
	case 0:
		c = Status((value & uint8(N+1)) >> 7)
		value <<= 1
	case 1:
		c = Status(value & uint8(C+1))
		value >>= 1
	}

	cpu.registers.P &= ^(C + 1)
	cpu.registers.P |= c

	return value
}

func (cpu *CPU) Jmp(addr uint16) {
	cpu.registers.PC = addr
}

func (cpu *CPU) Dec(addr uint16) {
	value := cpu.memory.fetch(addr)
	cpu.memory.store(addr, cpu.setZNFlags(value-1))
}

func (cpu *CPU) setZFlag(value uint8) uint8 {
	if value == 0 {
		cpu.registers.P = Status(setBit(cpu.registers.P, Z))
	} else {
		cpu.registers.P = Status(clearBit(cpu.registers.P, Z))
	}

	return value
}

func (cpu *CPU) setNFlag(value uint8) uint8 {
	if getBit(uint8(value), uint8(N)) != 0 {
		cpu.registers.P = Status(setBit(cpu.registers.P, N))
	} else {
		cpu.registers.P = Status(clearBit(cpu.registers.P, N))
	}

	return value
}

func (cpu *CPU) setZNFlags(value uint8) uint8 {
	cpu.setZFlag(value)
	cpu.setNFlag(value)
	return value
}

//Inc sets
func (cpu *CPU) Inc(addr uint16) {
	value := cpu.memory.fetch(addr)
	cpu.memory.store(addr, cpu.setZNFlags(value+1))
}

func (cpu *CPU) Sbc(addr uint16) {
	value := uint16(cpu.memory.fetch(addr))

	if getBit(uint8(cpu.registers.P), uint8(D)) == 0 {
		value ^= 0xff
	} else {
		value = 0x99 - value
	}

	registerA := cpu.registers.A
	if getBit(uint8(cpu.registers.P), uint8(D)) == 0 {
		result := cpu.setCFlagAddition(int(registerA) + int(value) + int(getBit(uint8(cpu.registers.P), uint8(C))))
		cpu.registers.A = cpu.setZNFlags(cpu.setVFlagAddition(uint16(registerA), uint16(value), uint16(result)))
	} else {

		low := uint16(registerA&0x000f) + uint16(value&0x000f) + uint16(getBit(uint8(cpu.registers.P), uint8(C)))
		high := uint16(registerA&0x00f0) + uint16(value&0x00f0)

		result := cpu.setCFlagAddition(int(high | (low & 0x000f)))

		result = cpu.setVFlagAddition(uint16(registerA), uint16(value), uint16(result))
		result = cpu.setZNFlags(result)
		cpu.registers.A = result
	}
}

//Adc add with carry | Função dando bugs, consertar depois!!!!!!!!
func (cpu *CPU) Adc(addr uint16) {

	value := uint16(cpu.memory.fetch(addr))
	registerA := cpu.registers.A
	if getBit(uint8(cpu.registers.P), uint8(D)) == 0 {
		result := cpu.setCFlagAddition(int(registerA) + int(value) + int(getBit(uint8(cpu.registers.P), uint8(C))))
		cpu.registers.A = cpu.setZNFlags(cpu.setVFlagAddition(uint16(registerA), uint16(value), uint16(result)))
	} else {

		low := uint16(registerA&0x000f) + uint16(value&0x000f) + uint16(getBit(uint8(cpu.registers.P), uint8(C)))
		high := uint16(registerA&0x00f0) + uint16(value&0x00f0)

		/*if low >= 0x000a {
			low -= 0x000a
			high += 0x0010
		}

		if high >= 0x00a0 {
			high -= 0x00a0
		}*/
		result := cpu.setCFlagAddition(int(high | (low & 0x000f)))

		result = cpu.setVFlagAddition(uint16(registerA), uint16(value), uint16(result))
		result = cpu.setZNFlags(result)
		cpu.registers.A = result
	}
}

func (cpu *CPU) setCFlagAddition(value int) uint8 {
	if value > 0xFF {
		cpu.registers.P = setBit(cpu.registers.P, C)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, C)
	}
	return uint8(value)
}

func (cpu *CPU) setVFlagAddition(num1 uint16, num2 uint16, result uint16) uint8 {
	if ((num1^num2)&0x80 == 0x0) && ((num1^result)&0x80 == 0x80) {
		cpu.registers.P = setBit(cpu.registers.P, V)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, V)
	}
	return uint8(result)
}

func (cpu *CPU) And(addr uint16) {
	value := cpu.memory.fetch(addr)
	cpu.registers.A = cpu.setZNFlags(cpu.registers.A & value)

}

func (cpu *CPU) Ora(addr uint16) {
	value := cpu.memory.fetch(addr)
	cpu.registers.A = cpu.setZNFlags(cpu.registers.A | value)
}

func (cpu *CPU) Eor(addr uint16) {
	value := cpu.memory.fetch(addr)
	cpu.registers.A = cpu.setZNFlags(cpu.registers.A ^ value)
}

//compareMemReg is a generic function to compare memory values with registers
func (cpu *CPU) compareMemReg(addr uint16, register *byte) {
	value := cpu.memory.fetch(addr)

	result := *register - value

	if *register >= value {
		cpu.registers.P = setBit(cpu.registers.P, C)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, C)
	}

	if *register == value {
		cpu.registers.P = setBit(cpu.registers.P, Z)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, Z)
	}

	if getBit(uint8(result), uint8(N)) == 1 {
		cpu.registers.P = setBit(cpu.registers.P, N)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, N)
	}
}

//Cmp instruction compares the contents of the accumulator with another memory held value and sets the zero and carry flags as appropriate.
func (cpu *CPU) Cmp(addr uint16) {
	cpu.compareMemReg(addr, &cpu.registers.A)
}

//Cpy compares mem with register Y
func (cpu *CPU) Cpy(addr uint16) {
	cpu.compareMemReg(addr, &cpu.registers.Y)
}

//Cpx compares mem with register X
func (cpu *CPU) Cpx(addr uint16) {
	cpu.compareMemReg(addr, &cpu.registers.X)
}

func (cpu *CPU) incrementRegister(register *byte) {
	*register++
	if *register == 0 {
		cpu.registers.P = setBit(cpu.registers.P, Z)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, Z)
	}

	if getBit(uint8(*register), uint8(N)) == 1 {
		cpu.registers.P = setBit(cpu.registers.P, N)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, N)
	}
}

func (cpu *CPU) decrementRegister(register *byte) {
	*register--
	if *register == 0 {
		cpu.registers.P = setBit(cpu.registers.P, Z)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, Z)
	}

	if getBit(uint8(*register), uint8(N)) == 1 {
		cpu.registers.P = setBit(cpu.registers.P, N)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, N)
	}
}

//Iny
func (cpu *CPU) Iny() {
	cpu.incrementRegister(&cpu.registers.Y)
}

//Inx
func (cpu *CPU) Inx() {
	cpu.incrementRegister(&cpu.registers.X)
}

//Dey
func (cpu *CPU) Dey() {
	cpu.decrementRegister(&cpu.registers.Y)
}

//Dey
func (cpu *CPU) Dex() {
	cpu.decrementRegister(&cpu.registers.X)
}

//Ta is a transfer accumulator
func (cpu *CPU) TaReg(register *byte) {
	*register = cpu.registers.A
	if *register == 0 {
		cpu.registers.P = setBit(cpu.registers.P, Z)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, Z)
	}

	if getBit(uint8(*register), uint8(N)) == 1 {
		cpu.registers.P = setBit(cpu.registers.P, N)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, N)
	}
}

//Tay copies the current contents of the accumulator into the Y register
func (cpu *CPU) Tay() {
	cpu.TaReg(&cpu.registers.Y)
}

//Tax copies the current contents of the accumulator into the X register
func (cpu *CPU) Tax() {
	cpu.TaReg(&cpu.registers.X)
}

//TRegA copies the contents of the register to accumulator
func (cpu *CPU) TRegA(register *byte) {
	cpu.registers.A = *register

	if *register == 0 {
		cpu.registers.P = setBit(cpu.registers.P, Z)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, Z)
	}

	if getBit(uint8(*register), uint8(N)) == 1 {
		cpu.registers.P = setBit(cpu.registers.P, N)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, N)
	}
}

//Tya copies the current contents of the Y register into the accumulator
func (cpu *CPU) Tya() {
	cpu.TRegA(&cpu.registers.Y)
}

//Txa copies the current contents of the X register into the accumulator
func (cpu *CPU) Txa() {
	cpu.TRegA(&cpu.registers.X)
}

//TSpReg copies the current contents of the stack register into the register
func (cpu *CPU) TSpReg(register *byte) {

	*register = cpu.registers.SP

	if *register == 0 {
		cpu.registers.P = setBit(cpu.registers.P, Z)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, Z)
	}

	if getBit(uint8(*register), uint8(N)) == 1 {
		cpu.registers.P = setBit(cpu.registers.P, N)
	} else {
		cpu.registers.P = clearBit(cpu.registers.P, N)
	}
}

func (cpu *CPU) Tsx() {
	cpu.TSpReg(&cpu.registers.X)
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

	path := "tests/nestest.nes"
	data, _ := ioutil.ReadFile(path)

	if string(data[:3]) != "NES" {
		log.Fatalf("Invalid ROM file" + string(data[:3]))
	}
	//header := data[:16]

	//fmt.Print(header)
	//fmt.Printf("\n %s | %02X | PRG-ROM: %02X | CHR-ROM: %02X \n| ROM CTRL BYTE 1: %02X | ROM CTRL BYTE 2: %02X |  N RAM BANKS: %02X\n", string(header[:3]), header[3:4], header[4:5], header[5:6], header[6:7], header[7:8], header[7:9])
	NumPRG := uint16(data[5])
	//NumCHR := header[5:6][0]

	prg := data[16 : 16+NumPRG*16384]

	//cpu.printRegisters()
	cpu.registers.P = 36
	cpu.registers.PC = 0xC000
	cpu.memory.loadCode(prg, cpu.registers.PC)
	cpu.run()

	fmt.Printf("\n Número de instruções: %d\n", len(cpu.instructionsTable))
	//cpu.printRegisters()

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
