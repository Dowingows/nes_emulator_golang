package main

//Status is a uint8
type Status uint8

//Flags used by P (Status) register
const (
	C Status = 1 << iota // carry flag
	Z                    // zero flag
	I                    // interrupt disable
	D                    // decimal mode
	B                    // break command
	U                    // -UNUSED-
	V                    // overflow flag
	N                    // negative flag
)

//Registers of processor
type Registers struct {
	A  byte
	X  byte
	Y  byte
	P  Status
	SP byte
	PC uint16
}

func (reg *Registers) init() *Registers {
	if reg == nil {
		reg = &Registers{}
		reg.reset()
		return reg
	}

	reg.reset()
	return reg

}

func (reg *Registers) reset() {
	reg.A = 0
	reg.X = 0
	reg.Y = 0
	reg.P = I
	reg.SP = 0xfd //253
	//reg.PC = 0xfffc //65532
	reg.PC = 0xff00 //65532
}
