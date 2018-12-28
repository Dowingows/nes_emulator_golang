package main

import "fmt"

//max Memory
const (
	DefaultMemorySize uint32 = 65536
)

//Memory has 65536 8-bit adresses
type Memory [DefaultMemorySize]uint8

func (memory *Memory) init() *Memory {
	if memory == nil {
		memory := Memory{}
		return &memory
	}
	memory.reset()
	return memory
}

func (memory *Memory) reset() {

}

func (memory *Memory) loadCode(code []byte, initAddr uint16) {
	for i := 0; i < len(code); i++ {
		memory[initAddr+uint16(i)] = code[i]
	}
}

func (memory *Memory) fetch(address uint16) uint8 {
	return memory[address]
}

func (memory *Memory) store(address uint16, value uint8) {
	memory[address] = value
}

func (memory *Memory) print() {
	max := 10
	fmt.Println()
	for i := 0; i < max; i++ {
		fmt.Printf(" %x ", memory[i])
	}
	fmt.Println()
}
