package cpu

// Simple operands

func (c *CPU) getA() uint8 {
	return c.a
}

func (c *CPU) getF() uint8 {
	return c.f
}

func (c *CPU) getB() uint8 {
	return c.b
}

func (c *CPU) getC() uint8 {
	return c.c
}

func (c *CPU) getD() uint8 {
	return c.d
}

func (c *CPU) getE() uint8 {
	return c.e
}

func (c *CPU) getH() uint8 {
	return c.h
}

func (c *CPU) getL() uint8 {
	return c.l
}

func (c *CPU) getM() uint8 {
	return c.ReadMem(c.getHL())
}

// Set

func (c *CPU) setA(v uint8) {
	c.a = v
}

func (c *CPU) setF(v uint8) {
	c.f = v&0xD7 | 0x02
}

func (c *CPU) setB(v uint8) {
	c.b = v
}

func (c *CPU) setC(v uint8) {
	c.c = v
}

func (c *CPU) setD(v uint8) {
	c.d = v
}

func (c *CPU) setE(v uint8) {
	c.e = v
}

func (c *CPU) setH(v uint8) {
	c.h = v
}

func (c *CPU) setL(v uint8) {
	c.l = v
}

func (c *CPU) setM(v uint8) {
	c.WriteMem(c.getHL(), v)
}

// Double registers

// Get

func (c *CPU) getAF() uint16 {
	return uint16(c.a)<<8 | uint16(c.f)
}

func (c *CPU) getBC() uint16 {
	return uint16(c.b)<<8 | uint16(c.c)
}

func (c *CPU) getDE() uint16 {
	return uint16(c.d)<<8 | uint16(c.e)
}

func (c *CPU) getHL() uint16 {
	return uint16(c.h)<<8 | uint16(c.l)
}

func (c *CPU) getSP() uint16 {
	return c.sp
}

// Set

func (c *CPU) setAF(v uint16) {
	c.setA(uint8(v >> 8))
	c.setF(uint8(v))
}

func (c *CPU) setBC(v uint16) {
	c.setB(uint8(v >> 8))
	c.setC(uint8(v))
}

func (c *CPU) setDE(v uint16) {
	c.setD(uint8(v >> 8))
	c.setE(uint8(v))
}

func (c *CPU) setHL(v uint16) {
	c.setH(uint8(v >> 8))
	c.setL(uint8(v))
}

func (c *CPU) setSP(v uint16) {
	c.sp = v
}

// Flags

// Get

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

// Set

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
		return c.getA()
	case "F":
		return c.getF()
	case "B":
		return c.getB()
	case "C":
		return c.getC()
	case "D":
		return c.getD()
	case "E":
		return c.getE()
	case "H":
		return c.getH()
	case "L":
		return c.getL()
	case "M":
		return c.getM()
	default:
		panic("unsupported operand: " + op)
	}
}

func (c *CPU) setOp(op string, v uint8) {
	switch op {
	case "A":
		c.setA(v)
	case "F":
		c.setF(v)
	case "B":
		c.setB(v)
	case "C":
		c.setC(v)
	case "D":
		c.setD(v)
	case "E":
		c.setE(v)
	case "H":
		c.setH(v)
	case "L":
		c.setL(v)
	case "M":
		c.setM(v)
	default:
		panic("unsupported operand: " + op)
	}
}

func (c *CPU) getDoubleOp(op string) uint16 {
	switch op {
	case "AF":
		return c.getAF()
	case "BC":
		return c.getBC()
	case "DE":
		return c.getDE()
	case "HL":
		return c.getHL()
	case "SP":
		return c.getSP()
	default:
		panic("unsupported operand: " + op)
	}
}

func (c *CPU) setDoubleOp(op string, v uint16) {
	switch op {
	case "AF":
		c.setAF(v)
	case "BC":
		c.setBC(v)
	case "DE":
		c.setDE(v)
	case "HL":
		c.setHL(v)
	case "SP":
		c.setSP(v)
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
