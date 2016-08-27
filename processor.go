package main

// processor struct
type pro struct {
	ER, FR, IR, PC, CR, IPC uint32     // single-purpose registers
	R                       [32]uint32 // Avaible registers
}

/*
* processor Methods
 */
func (p *pro) getFlag(f string) (n uint8) {
	var mask uint32 = 0x00000020
	o := map[string]uint8{"IE": 0, "OV": 1, "ZD": 2, "GT": 3, "LT": 4, "EQ": 5}
	n = uint8(((mask >> o[f]) & p.FR) >> (5 - o[f]))
	return
}

func (p *pro) writeFlag(flag string, i uint32) uint32 {
	o := map[string]uint32{"IE": 5, "OV": 4, "ZD": 3, "GT": 2, "LT": 1, "EQ": 0}
	if i == 1 {
		processor.FR |= 1 << o[flag]
	} else {
		processor.FR &= ^(1 << o[flag])
	}
	return processor.FR
}

// Software interruption
func (p *pro) checkSoftInt() (ok bool) {
	if n := p.getFlag("IE"); n == 1 {
		processor.IPC = processor.PC + 1
		processor.PC = 3
		newBranch = true
		return false
	}
	return true
}

// Interruption routine
func (p *pro) makeInt() {
	p.IPC = processor.PC // PC already incremented when hit here
	p.PC = 1
	newBranch = true
	processor.CR = 0xFF
}

func (p *pro) checkOverflow(n uint64) {
	if n > 0xFFFFFFFF {
		p.writeFlag("OV", 1)
		processor.ER = uint32(n >> 32)
	} else {
		p.writeFlag("OV", 0)
		processor.ER = 0
	}
}

func (p *pro) getExtension(n uint64) (max, min uint32) {
	min = uint32(n)
	max = uint32(n >> 32)
	return
}
