package main

import (
	"fmt"
)

// TypeU defines the struct of the regiter type U
type TypeU struct {
	fn         string
	Rx, Ry, Rz uint8
}

func (t *TypeU) decode(i uint32) *TypeU {
	t.Rz = uint8((0x00007C00 & i) >> 10)
	t.Rx = uint8((0x000003E0 & i) >> 5)
	t.Ry = uint8(0x0000001F & i)
	return t
}

// Add is the addition operation
func (t *TypeU) Add() {
	Result = append(Result, fmt.Sprintf("add r%d, r%d, r%d\n", t.Rz, t.Rx, t.Ry))

	sum := uint64(processor.R[t.Rx]) + uint64(processor.R[t.Ry])
	processor.checkOverflow(sum)
	processor.R[t.Rz] = processor.R[t.Rx] + processor.R[t.Ry]

	Result = append(Result, fmt.Sprintf("[U] FR = 0x%.8X, R%d = R%d + R%d = 0x%.8X\n",
		processor.FR, t.Rz, t.Rx, t.Ry, processor.R[t.Rz]))
}

// And logic operation
func (t *TypeU) And() {
	Result = append(Result, fmt.Sprintf("and r%d, r%d, r%d\n", t.Rz, t.Rx, t.Ry))
	processor.R[t.Rz] = processor.R[t.Rx] & processor.R[t.Ry]
	Result = append(Result, fmt.Sprintf("[U] R%d = R%d & R%d = 0x%.8X\n",
		t.Rz, t.Rx, t.Ry, processor.R[t.Rz]))
}

// Div divider two values within registers. May set OVERFLOW to 0 and ZD to 1 or 0
func (t *TypeU) Div() {
	Result = append(Result, fmt.Sprintf("div r%d, r%d, r%d\n", t.Rz, t.Rx, t.Ry))
	ok := true
	if processor.R[t.Ry] == 0 {
		processor.writeFlag("ZD", 1)

		ok = processor.checkSoftInt()
		if !ok {
			processor.CR = 1
		}

	} else {
		processor.writeFlag("ZD", 0)
		processor.writeFlag("OV", 0)
		processor.R[t.Rz] = processor.R[t.Rx] / processor.R[t.Ry]

		if mod := processor.R[t.Rx] % processor.R[t.Ry]; mod > 0 {
			processor.ER = mod
		} else {
			processor.ER = 0
		}
	}
	Result = append(Result, fmt.Sprintf("[U] FR = 0x%.8X, ER = 0x%.8X, R%d = R%d / R%d = 0x%.8X\n",
		processor.FR, processor.ER, t.Rz, t.Rx, t.Ry, processor.R[t.Rz]))

	if !ok {
		Result = append(Result, fmt.Sprintf("[SOFTWARE INTERRUPTION]\n"))
	}
}

// Mul instructions multiplacates of two value within register
func (t *TypeU) Mul() {
	Result = append(Result, fmt.Sprintf("mul r%d, r%d, r%d\n", t.Rz, t.Rx, t.Ry))

	mul := uint64(processor.R[t.Rx]) * uint64(processor.R[t.Ry])
	processor.checkOverflow(mul)

	processor.R[t.Rz] = processor.R[t.Rx] * processor.R[t.Ry]

	Result = append(Result, fmt.Sprintf("[U] FR = 0x%.8X, ER = 0x%.8X, R%d = R%d * R%d = 0x%.8X\n",
		processor.FR, processor.ER, t.Rz, t.Rx, t.Ry, processor.R[t.Rz]))
}

// Not logic operation
func (t *TypeU) Not() {
	Result = append(Result, fmt.Sprintf("not r%d, r%d\n", t.Rx, t.Ry))
	processor.R[t.Rx] = ^processor.R[t.Ry]
	Result = append(Result, fmt.Sprintf("[F] R%d = ~R%d = 0x%.8X\n",
		t.Rx, t.Ry, processor.R[t.Rx]))
}

// Or logic operation
func (t *TypeU) Or() {
	Result = append(Result, fmt.Sprintf("or r%d, r%d, r%d\n", t.Rz, t.Rx, t.Ry))
	processor.R[t.Rz] = processor.R[t.Rx] | processor.R[t.Ry]
	Result = append(Result, fmt.Sprintf("[U] R%d = R%d | R%d = 0x%.8X\n",
		t.Rz, t.Rx, t.Ry, processor.R[t.Rz]))
}

// Shl shifts to left
func (t *TypeU) Shl() {
	Result = append(Result, fmt.Sprintf("shl r%d, r%d, %d\n", t.Rz, t.Rx, t.Ry))

	temp := (uint64(processor.ER) << 32) | uint64(processor.R[t.Rx])
	temp = temp << (t.Ry + 1)

	ER, R := processor.getExtension(temp)

	processor.ER = ER
	processor.R[t.Rz] = R

	Result = append(Result, fmt.Sprintf("[U] ER = 0x%.8X, R%d = R%d << %d = 0x%.8X\n",
		processor.ER, t.Rz, t.Rx, t.Ry+1, processor.R[t.Rz]))
}

// Shr shifts to right
func (t *TypeU) Shr() {
	Result = append(Result, fmt.Sprintf("shr r%d, r%d, %d\n", t.Rz, t.Rx, t.Ry))
	//Cocatenate
	temp := uint64(processor.ER) << 32
	temp = temp | uint64(processor.R[t.Rx])

	//Execute
	temp = temp >> (t.Ry + 1)

	//Split
	ER, R := processor.getExtension(temp)

	processor.ER = ER
	processor.R[t.Rz] = R

	Result = append(Result, fmt.Sprintf("[U] ER = 0x%.8X, R%d = R%d >> %d = 0x%.8X\n",
		processor.ER, t.Rz, t.Rx, t.Ry+1, processor.R[t.Rz]))
}

func (t *TypeU) Sub() {
	Result = append(Result, fmt.Sprintf("sub r%d, r%d, r%d\n", t.Rz, t.Rx, t.Ry))

	sum := uint64(processor.R[t.Rx]) - uint64(processor.R[t.Ry])
	processor.checkOverflow(sum)
	processor.R[t.Rz] = processor.R[t.Rx] - processor.R[t.Ry]

	Result = append(Result, fmt.Sprintf("[U] FR = 0x%.8X, R%d = R%d - R%d = 0x%.8X\n",
		processor.FR, t.Rz, t.Rx, t.Ry, processor.R[t.Rz]))
}

// Xor logic operation
func (t *TypeU) Xor() {
	Result = append(Result, fmt.Sprintf("xor r%d, r%d, r%d\n", t.Rz, t.Rx, t.Ry))
	processor.R[t.Rz] = processor.R[t.Rx] ^ processor.R[t.Ry]
	Result = append(Result, fmt.Sprintf("[U] R%d = R%d ^ R%d = 0x%.8X\n",
		t.Rz, t.Rx, t.Ry, processor.R[t.Rz]))

}
