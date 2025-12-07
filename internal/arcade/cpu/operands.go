package cpu

// Flag getters.
func (c *CPU) getSF() uint8 {
	return c.F >> 7 & 1
}

func (c *CPU) getZF() uint8 {
	return c.F >> 6 & 1
}

func (c *CPU) getACF() uint8 {
	return c.F >> 4 & 1
}

func (c *CPU) getPF() uint8 {
	return c.F >> 2 & 1
}

func (c *CPU) getCYF() uint8 {
	return c.F & 1
}

// Flag setters.
func (c *CPU) setSF(b bool) {
	if b {
		c.F = c.F | 0x80
	} else {
		c.F = c.F & 0x7F
	}
}

func (c *CPU) setZF(b bool) {
	if b {
		c.F = c.F | 0x40
	} else {
		c.F = c.F & 0xBF
	}
}

func (c *CPU) setACF(b bool) {
	if b {
		c.F = c.F | 0x10
	} else {
		c.F = c.F & 0xEF
	}
}

func (c *CPU) setPF(b bool) {
	if b {
		c.F = c.F | 0x04
	} else {
		c.F = c.F & 0xFB
	}
}

func (c *CPU) setCYF(b bool) {
	if b {
		c.F = c.F | 0x01
	} else {
		c.F = c.F & 0xFE
	}
}

func (c *CPU) getOp(op string) uint8 {
	switch op {
	case "A":
		return c.A
	case "F":
		return c.F
	case "B":
		return c.B
	case "C":
		return c.C
	case "D":
		return c.D
	case "E":
		return c.E
	case "H":
		return c.H
	case "L":
		return c.L
	case "M":
		return c.Bus.Read(uint16(c.H)<<8 | uint16(c.L))
	default:
		panic("unsupported operand: " + op)
	}
}

func (c *CPU) setOp(op string, v uint8) {
	switch op {
	case "A":
		c.A = v
	case "F":
		c.F = v
	case "B":
		c.B = v
	case "C":
		c.C = v
	case "D":
		c.D = v
	case "E":
		c.E = v
	case "H":
		c.H = v
	case "L":
		c.L = v
	case "M":
		c.Bus.Write(uint16(c.H)<<8|uint16(c.L), v)
	default:
		panic("unsupported operand: " + op)
	}
}

func (c *CPU) getDoubleOp(op string) uint16 {
	switch op {
	case "AF":
		return uint16(c.A)<<8 | uint16(c.F)
	case "BC":
		return uint16(c.B)<<8 | uint16(c.C)
	case "DE":
		return uint16(c.D)<<8 | uint16(c.E)
	case "HL":
		return uint16(c.H)<<8 | uint16(c.L)
	case "SP":
		return c.SP
	default:
		panic("unsupported operand: " + op)
	}
}

func (c *CPU) setDoubleOp(op string, v uint16) {
	switch op {
	case "AF":
		c.A = uint8(v >> 8)
		c.F = uint8(v)&0xD7 | 0x02
	case "BC":
		c.B = uint8(v >> 8)
		c.C = uint8(v)
	case "DE":
		c.D = uint8(v >> 8)
		c.E = uint8(v)
	case "HL":
		c.H = uint8(v >> 8)
		c.L = uint8(v)
	case "SP":
		c.SP = v
	default:
		panic("unsupported operand: " + op)
	}
}

func (c *CPU) setFlags(s, z, ac, p, cy bool) {
	c.setSF(s)
	c.setZF(z)
	c.setACF(ac)
	c.setPF(p)
	c.setCYF(cy)
}
