package main

import (
	"fmt"
	"math"
)

// MathDevice defines the struct of the Math Device
type MathDevice struct {
	x, y, z, op, s, clock uint32
	inProgress            bool
}

/*
* Math Unity methods
 */
func (m *MathDevice) reset() {
	m.op = 0
	m.inProgress = false
	m.clock = 0
}

func (m *MathDevice) operate() {
	switch m.op {
	case 0:
	case 1:
		m.z = uint32(math.Floor(math.Sqrt(float64((m.x * m.x) + (m.y * m.y)))))
		m.s = 0
		m.inProgress = true
	case 2:
		m.z = uint32(math.Pow(float64(m.x), float64(m.y)))
		m.s = 0
		m.inProgress = true
	case 3:
		m.z = (m.x + m.y) / 2
		m.s = 0
		m.inProgress = true
	default:
		m.s = 1
		m.inProgress = true
	}
}

func (m *MathDevice) readDevice(p, v uint32) (n uint32, r bool) {
	if p > uint32(len(Mem)) {
		p = p << 2
		switch p {
		case 0x00008000:
			m.x = v
		case 0x00008004:
			m.y = v
		case 0x00008008:
			n = m.z
		case 0x0000800C:
			m.op = v
			m.operate()

			if m.s == 0 {
				n = m.op
			} else {
				n = 0x20
			}

		case 0x00008888:
			Terminal = append(Terminal, byte(v))
		default:
			fmt.Printf("Memory error")
		}
		r = true
		return
	}
	r = false
	return
}
