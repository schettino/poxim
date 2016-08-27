package main

import "fmt"

// TypeS defines the struct of the regiter type U
type TypeS struct {
	fn   string
	Im26 uint32
}

func (t *TypeS) decode(i uint32) {
	t.Im26 = (0x03FFFFFF & i)
}

func (t *TypeS) Bgt() {
	Result = append(Result, fmt.Sprintf("bgt 0x%.8X\n", t.Im26))

	newBranch = true
	if processor.getFlag("GT") == 1 {
		processor.PC = t.Im26
	} else {
		processor.PC++
	}

	Result = append(Result, fmt.Sprintf("[S] PC = 0x%.8X\n", processor.PC<<2))
}

func (t *TypeS) Beq() {
	Result = append(Result, fmt.Sprintf("beq 0x%.8X\n", t.Im26))

	newBranch = true
	if processor.getFlag("EQ") == 1 {
		processor.PC = t.Im26
	} else {
		processor.PC++
	}

	Result = append(Result, fmt.Sprintf("[S] PC = 0x%.8X\n", processor.PC<<2))
}

func (t *TypeS) Blt() {
	Result = append(Result, fmt.Sprintf("blt 0x%.8X\n", t.Im26))

	newBranch = true
	if processor.getFlag("LT") == 1 {
		processor.PC = t.Im26
	} else {
		processor.PC++
	}

	Result = append(Result, fmt.Sprintf("[S] PC = 0x%.8X\n", processor.PC<<2))
}

func (t *TypeS) Bne() {
	Result = append(Result, fmt.Sprintf("bne 0x%.8X\n", t.Im26))

	newBranch = true
	//Result = append(Result, fmt.Sprintf(" !!-- EQ:%d\n", ^processor.getFlag("EQ")))
	if ^processor.getFlag("EQ") == 0xFF {
		processor.PC = t.Im26
	} else {
		processor.PC++
	}

	Result = append(Result, fmt.Sprintf("[S] PC = 0x%.8X\n", processor.PC<<2))
}

func (t *TypeS) Ble() {
	Result = append(Result, fmt.Sprintf("ble 0x%.8X\n", t.Im26))

	newBranch = true
	if processor.getFlag("EQ") == 1 || processor.getFlag("LT") == 1 {
		processor.PC = t.Im26
	} else {
		processor.PC++
	}

	Result = append(Result, fmt.Sprintf("[S] PC = 0x%.8X\n", processor.PC<<2))
}

func (t *TypeS) Bge() {
	Result = append(Result, fmt.Sprintf("bge 0x%.8X\n", t.Im26))

	newBranch = true
	if processor.getFlag("EQ") == 1 || processor.getFlag("GT") == 1 {
		processor.PC = t.Im26
	} else {
		processor.PC++
	}

	Result = append(Result, fmt.Sprintf("[S] PC = 0x%.8X\n", processor.PC<<2))
}

func (t *TypeS) Bun() {
	Result = append(Result, fmt.Sprintf("bun 0x%.8X\n", t.Im26))
	processor.PC = t.Im26
	newBranch = true
	Result = append(Result, fmt.Sprintf("[S] PC = 0x%.8X\n", t.Im26<<2))
}

func (t *TypeS) Int() {
	Result = append(Result, fmt.Sprintf("int %d\n", t.Im26))
	ok := true
	if t.Im26 == 0 {
		processor.CR = 0
		processor.PC = 0
		endState = true
	} else {
		ok = processor.checkSoftInt()
		if !ok {
			processor.CR = t.Im26
		}
	}
	Result = append(Result, fmt.Sprintf("[S] CR = 0x%.8X, PC = 0x%.8X\n", processor.CR, processor.PC<<2))

	if !ok {
		Result = append(Result, fmt.Sprintf("[SOFTWARE INTERRUPTION]\n"))
	} else {
		if len(Terminal) > 0 {
			Result = append(Result, "[TERMINAL]\n")
			Result = append(Result, string(Terminal))
			Result = append(Result, "\n")
		}
		Result = append(Result, "[END OF SIMULATION]")
	}
}
