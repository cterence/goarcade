package cpu

// Flag getters.
func (c *CPU) getSF() uint8 {
	return c.f >> 7 & 1
}

func (c *CPU) getZF() uint8 {
	return c.f >> 6 & 1
}

func (c *CPU) getACF() uint8 {
	return c.f >> 4 & 1
}

func (c *CPU) getPF() uint8 {
	return c.f >> 2 & 1
}

func (c *CPU) getCYF() uint8 {
	return c.f & 1
}

// Flag setters.
func (c *CPU) setSF(b bool) {
	if b {
		c.f = c.f | 0x80
	} else {
		c.f = c.f & 0x7F
	}
}

func (c *CPU) setZF(b bool) {
	if b {
		c.f = c.f | 0x40
	} else {
		c.f = c.f & 0xBF
	}
}

func (c *CPU) setACF(b bool) {
	if b {
		c.f = c.f | 0x10
	} else {
		c.f = c.f & 0xEF
	}
}

func (c *CPU) setPF(b bool) {
	if b {
		c.f = c.f | 0x04
	} else {
		c.f = c.f & 0xFB
	}
}

func (c *CPU) setCYF(b bool) {
	if b {
		c.f = c.f | 0x01
	} else {
		c.f = c.f & 0xFE
	}
}

func (c *CPU) getOp(op string) uint8 {
	switch op {
	case "A":
		return c.a
	case "F":
		return c.f
	case "B":
		return c.b
	case "C":
		return c.c
	case "D":
		return c.d
	case "E":
		return c.e
	case "H":
		return c.h
	case "L":
		return c.l
	case "M":
		return c.ReadMem(uint16(c.h)<<8 | uint16(c.l))
	default:
		panic("unsupported operand: " + op)
	}
}

func (c *CPU) setOp(op string, v uint8) {
	switch op {
	case "A":
		c.a = v
	case "F":
		c.f = v
	case "B":
		c.b = v
	case "C":
		c.c = v
	case "D":
		c.d = v
	case "E":
		c.e = v
	case "H":
		c.h = v
	case "L":
		c.l = v
	case "M":
		c.WriteMem(uint16(c.h)<<8|uint16(c.l), v)
	default:
		panic("unsupported operand: " + op)
	}
}

func (c *CPU) getDoubleOp(op string) uint16 {
	switch op {
	case "AF":
		return uint16(c.a)<<8 | uint16(c.f)
	case "BC":
		return uint16(c.b)<<8 | uint16(c.c)
	case "DE":
		return uint16(c.d)<<8 | uint16(c.e)
	case "HL":
		return uint16(c.h)<<8 | uint16(c.l)
	case "SP":
		return c.sp
	default:
		panic("unsupported operand: " + op)
	}
}

func (c *CPU) setDoubleOp(op string, v uint16) {
	switch op {
	case "AF":
		c.a = uint8(v >> 8)
		c.f = uint8(v)&0xD7 | 0x02
	case "BC":
		c.b = uint8(v >> 8)
		c.c = uint8(v)
	case "DE":
		c.d = uint8(v >> 8)
		c.e = uint8(v)
	case "HL":
		c.h = uint8(v >> 8)
		c.l = uint8(v)
	case "SP":
		c.sp = v
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
