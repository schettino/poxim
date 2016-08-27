package main

import "fmt"

// TypeF defines the struct of the regiter type U
type TypeF struct {
	fn     string
	Im16   uint16
	Rx, Ry uint8
}

func (t *TypeF) decode(i uint32) {
	t.Im16 = uint16((0x03FFFC00 & i) >> 10)
	t.Rx = uint8((0x000003E0 & i) >> 5)
	t.Ry = uint8(0x0000001F & i)
}

func (t *TypeF) Sfr() {
	Result = append(Result, fmt.Sprintf("sfr r%d\n", t.Rx))
	processor.FR = processor.R[t.Rx]
	Result = append(Result, fmt.Sprintf("[F] FR = R%d = 0x%.8X\n", t.Rx, processor.R[t.Rx]))
}

func (t *TypeF) Lfr() {
	Result = append(Result, fmt.Sprintf("lfr r%d\n", t.Rx))
	processor.R[t.Rx] = processor.FR
	Result = append(Result, fmt.Sprintf("[F] R%d = FR = 0x%.8X\n", t.Rx, processor.R[t.Rx]))
}

func (t *TypeF) Isr() {
	Result = append(Result, fmt.Sprintf("isr r%d, r%d, 0x%.4X\n", t.Rx, t.Ry, t.Im16))
	processor.R[t.Rx] = processor.IPC
	processor.R[t.Ry] = processor.CR
	processor.PC = uint32(t.Im16)
	newBranch = true
	Result = append(Result, fmt.Sprintf("[F] R%d = IPC >> 2 = 0x%.8X, R%d = CR = 0x%.8X, PC = 0x%.8X\n",
		t.Rx, (processor.R[t.Rx]), t.Ry, processor.R[t.Ry], processor.PC<<2))
}

func (t *TypeF) Andi() {
	Result = append(Result, fmt.Sprintf("andi r%d, r%d, %d\n", t.Rx, t.Ry, t.Im16))
	processor.R[t.Rx] = processor.R[t.Ry] & uint32(t.Im16)
	Result = append(Result, fmt.Sprintf("[F] R%d = R%d & 0x%.4X = 0x%.8X\n",
		t.Rx, t.Ry, t.Im16, processor.R[t.Rx]))
}

func (t *TypeF) Noti() {
	Result = append(Result, fmt.Sprintf("noti r%d, %d\n", t.Rx, t.Im16))
	processor.R[t.Rx] = ^uint32(t.Im16)
	Result = append(Result, fmt.Sprintf("[F] R%d = ~0x%.4X = 0x%.8X\n",
		t.Rx, t.Im16, processor.R[t.Rx]))
}

func (t *TypeF) Ori() {
	Result = append(Result, fmt.Sprintf("ori r%d, r%d, %d\n", t.Rx, t.Ry, t.Im16))
	processor.R[t.Rx] = processor.R[t.Ry] | uint32(t.Im16)
	Result = append(Result, fmt.Sprintf("[F] R%d = R%d | 0x%.4X = 0x%.8X\n",
		t.Rx, t.Ry, t.Im16, processor.R[t.Rx]))
}

func (t *TypeF) Xori() {
	Result = append(Result, fmt.Sprintf("xori r%d, r%d, %d\n", t.Rx, t.Ry, t.Im16))
	processor.R[t.Rx] = processor.R[t.Ry] ^ uint32(t.Im16)
	Result = append(Result, fmt.Sprintf("[F] R%d = R%d ^ 0x%.4X = 0x%.8X\n",
		t.Rx, t.Ry, t.Im16, processor.R[t.Rx]))
}

func (t *TypeF) Ldw() {
	Result = append(Result, fmt.Sprintf("ldw r%d, r%d, 0x%.4X\n", t.Rx, t.Ry, t.Im16))

	p := processor.R[t.Ry] + uint32(t.Im16)

	n, r := MathUnit.readDevice(p, processor.R[t.Rx])

	if !r {
		processor.R[t.Rx] = Mem[p]
	} else {
		if n > 0 {
			processor.R[t.Rx] = n
		}
	}

	Result = append(Result, fmt.Sprintf("[F] R%d = MEM[(R%d + 0x%.4X) << 2] = 0x%.8X\n",
		t.Rx, t.Ry, t.Im16, processor.R[t.Rx]))
}

func (t *TypeF) Stw() {
	Result = append(Result, fmt.Sprintf("stw r%d, 0x%.4X, r%d\n", t.Rx, t.Im16, t.Ry))

	p := (processor.R[t.Rx] + uint32(t.Im16))

	_, r := MathUnit.readDevice(p, processor.R[t.Ry])

	if !r {
		Mem[p] = processor.R[t.Ry]
	} else {

	}

	Result = append(Result, fmt.Sprintf("[F] MEM[(R%d + 0x%.4X) << 2] = R%d = 0x%.8X\n",
		t.Rx, t.Im16, t.Ry, processor.R[t.Ry]))
}

func (t *TypeF) Ldb() {
	Result = append(Result, fmt.Sprintf("ldb r%d, r%d, 0x%.4X\n", t.Rx, t.Ry, t.Im16))

	mask := map[uint16]uint32{0x000: 0xFF000000, 0x001: 0x00FF0000, 0x002: 0x0000FF00, 0x003: 0x000000FF}
	d := map[uint16]uint8{0x000: 24, 0x001: 16, 0x002: 8, 0x003: 0}

	e := (processor.R[t.Ry] + uint32(t.Im16))

	n, r := MathUnit.readDevice((e >> 2), processor.R[t.Rx])

	if !r {
		p := Mem[e>>2]
		i := uint16(e % 4)
		processor.R[t.Rx] = (p & mask[i]) >> d[i]
	} else {
		if n > 0 {
			processor.R[t.Rx] = n
		}
	}

	Result = append(Result, fmt.Sprintf("[F] R%d = MEM[R%d + 0x%.4X] = 0x%.2X\n",
		t.Rx, t.Ry, t.Im16, processor.R[t.Rx] /*, ((Mem[n]&0xFF000000)>>24)+uint32(t.Im16)*/))
}

func (t *TypeF) Stb() {
	Result = append(Result, fmt.Sprintf("stb r%d, 0x%.4X, r%d\n", t.Rx, t.Im16, t.Ry))

	e := (processor.R[t.Rx] + uint32(t.Im16))

	n, r := MathUnit.readDevice((e >> 2), processor.R[t.Ry])
	var mr uint32

	if !r {
		Mem[e>>2] = (processor.R[t.Ry] & 0xFF)
		mr = Mem[e>>2]
	} else {
		if n > 0 {
			processor.R[t.Ry] = n
		}

		mr = processor.R[t.Ry]
	}

	Result = append(Result, fmt.Sprintf("[F] MEM[R%d + 0x%.4X] = R%d = 0x%.2X\n",
		t.Rx, t.Im16, t.Ry, mr))
}

func (t *TypeF) Addi() {
	Result = append(Result, fmt.Sprintf("addi r%d, r%d, %d\n", t.Rx, t.Ry, t.Im16))

	sum := uint64(processor.R[t.Ry]) + uint64(t.Im16)
	processor.checkOverflow(sum)
	processor.R[t.Rx] = processor.R[t.Ry] + uint32(t.Im16)

	Result = append(Result, fmt.Sprintf("[F] FR = 0x%.8X, R%d = R%d + 0x%.4X = 0x%.8X\n",
		processor.FR, t.Rx, t.Ry, t.Im16, processor.R[t.Rx]))
}

func (t *TypeF) Subi() {
	Result = append(Result, fmt.Sprintf("subi r%d, r%d, %d\n", t.Rx, t.Ry, t.Im16))

	sum := uint64(processor.R[t.Ry]) - uint64(t.Im16)
	processor.checkOverflow(sum)
	processor.R[t.Rx] = processor.R[t.Ry] - uint32(t.Im16)

	Result = append(Result, fmt.Sprintf("[F] FR = 0x%.8X, R%d = R%d - 0x%.4X = 0x%.8X\n",
		processor.FR, t.Rx, t.Ry, t.Im16, processor.R[t.Rx]))
}

func (t *TypeF) Divi() {
	Result = append(Result, fmt.Sprintf("divi r%d, r%d, %d\n", t.Rx, t.Ry, t.Im16))
	ok := true
	if t.Im16 == 0 {
		processor.writeFlag("ZD", 1)

		ok = processor.checkSoftInt()
		if !ok {
			processor.CR = 1
		}

	} else {
		processor.writeFlag("ZD", 0)
		processor.writeFlag("OV", 0)
		if mod := processor.R[t.Ry] % uint32(t.Im16); mod > 0 {
			processor.ER = mod
		} else {
			processor.ER = 0
		}
		processor.R[t.Rx] = processor.R[t.Ry] / uint32(t.Im16)
	}
	Result = append(Result, fmt.Sprintf("[F] FR = 0x%.8X, ER = 0x%.8X, R%d = R%d / 0x%.4X = 0x%.8X\n",
		processor.FR, processor.ER, t.Rx, t.Ry, t.Im16, processor.R[t.Rx]))

	if !ok {
		Result = append(Result, fmt.Sprintf("[SOFTWARE INTERRUPTION]\n"))
	}
}

func (t *TypeF) Muli() {
	Result = append(Result, fmt.Sprintf("muli r%d, r%d, %d\n", t.Rx, t.Ry, t.Im16))

	mul := uint64(processor.R[t.Ry]) * uint64(t.Im16)
	processor.checkOverflow(mul)

	processor.R[t.Rx] = processor.R[t.Ry] * uint32(t.Im16)

	Result = append(Result, fmt.Sprintf("[F] FR = 0x%.8X, ER = 0x%.8X, R%d = R%d * 0x%.4X = 0x%.8X\n",
		processor.FR, processor.ER, t.Rx, t.Ry, t.Im16, processor.R[t.Rx]))
}

func (t *TypeF) Cmp() {
	Result = append(Result, fmt.Sprintf("cmp r%d, r%d\n", t.Rx, t.Ry))
	switch {
	case processor.R[t.Rx] == processor.R[t.Ry]:
		processor.writeFlag("LT", 0)
		processor.writeFlag("GT", 0)
		processor.writeFlag("EQ", 1)
	case processor.R[t.Rx] < processor.R[t.Ry]:
		processor.writeFlag("EQ", 0)
		processor.writeFlag("GT", 0)
		processor.writeFlag("LT", 1)
	case processor.R[t.Rx] > processor.R[t.Ry]:
		processor.writeFlag("EQ", 0)
		processor.writeFlag("LT", 0)
		processor.writeFlag("GT", 1)
	}
	Result = append(Result, fmt.Sprintf("[F] FR = 0x%.8X\n", processor.FR))
}

func (t *TypeF) Cmpi() {
	Result = append(Result, fmt.Sprintf("cmpi r%d, %d\n", t.Rx, t.Im16))
	switch {
	case processor.R[t.Rx] == uint32(t.Im16):
		processor.writeFlag("LT", 0)
		processor.writeFlag("GT", 0)
		processor.writeFlag("EQ", 1)
	case processor.R[t.Rx] < uint32(t.Im16):
		processor.writeFlag("EQ", 0)
		processor.writeFlag("GT", 0)
		processor.writeFlag("LT", 1)
	case processor.R[t.Rx] > uint32(t.Im16):
		processor.writeFlag("EQ", 0)
		processor.writeFlag("LT", 0)
		processor.writeFlag("GT", 1)
	}
	Result = append(Result, fmt.Sprintf("[F] FR = 0x%.8X\n", processor.FR))
}

func (t *TypeF) Call() {
	Result = append(Result, fmt.Sprintf("call r%d, r%d, 0x%.4X\n", t.Rx, t.Ry, t.Im16))
	newBranch = true

	processor.R[t.Rx] = (processor.PC + 1)
	if t.Rx == 0 {
		processor.R[0] = 0
	}

	processor.PC = (processor.R[t.Ry] + uint32(t.Im16))

	Result = append(Result, fmt.Sprintf("[F] R%d = (PC + 4) >> 2 = 0x%.8X, PC = (R%d + 0x%.4X) << 2 = 0x%.8X\n",
		t.Rx, processor.R[t.Rx], t.Ry, t.Im16, (processor.PC<<2)))
}

func (t *TypeF) Ret() {
	Result = append(Result, fmt.Sprintf("ret r%d\n", t.Rx))
	newBranch = true
	processor.PC = processor.R[t.Rx]
	Result = append(Result, fmt.Sprintf("[F] PC = R%d << 2 = 0x%.8X\n", t.Rx, (processor.PC<<2)))
}
