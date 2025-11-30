package cpu

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

type CPU struct {
	ReadMem  func(uint16) uint8
	WriteMem func(uint16, uint8)

	Running bool

	inst string
	sc   uint64

	pc uint16
	sp uint16

	b uint8
	c uint8
	d uint8
	e uint8
	h uint8
	l uint8

	acc   uint8
	act   uint8
	flags uint8
	tmp   uint8
}

func (c *CPU) Init() {
	c.Running = true
	c.pc = 0
	c.sp = 0
	c.b = 0
	c.c = 0
	c.d = 0
	c.e = 0
	c.h = 0
	c.l = 0
	c.acc = 0
	c.act = 0
	c.flags = 2
	c.tmp = 0
}

func (c *CPU) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%-12s", c.inst))
	b.WriteString(" SC:" + fmt.Sprintf("%-5s", strconv.FormatUint(uint64(c.sc), 10)))
	b.WriteString(" PC:" + fmt.Sprintf("%-4s", strconv.FormatUint(uint64(c.pc), 16)))
	b.WriteString(" SP:" + fmt.Sprintf("%-4s", strconv.FormatUint(uint64(c.sp), 16)))
	b.WriteString(" B:" + fmt.Sprintf("%-2s", strconv.FormatUint(uint64(c.b), 16)))
	b.WriteString(" C:" + fmt.Sprintf("%-2s", strconv.FormatUint(uint64(c.c), 16)))
	b.WriteString(" D:" + fmt.Sprintf("%-2s", strconv.FormatUint(uint64(c.d), 16)))
	b.WriteString(" E:" + fmt.Sprintf("%-2s", strconv.FormatUint(uint64(c.e), 16)))
	b.WriteString(" H:" + fmt.Sprintf("%-2s", strconv.FormatUint(uint64(c.h), 16)))
	b.WriteString(" L:" + fmt.Sprintf("%-2s", strconv.FormatUint(uint64(c.l), 16)))
	b.WriteString(" ACC:" + fmt.Sprintf("%-2s", strconv.FormatUint(uint64(c.acc), 16)))
	b.WriteString(" ACT:" + fmt.Sprintf("%-2s", strconv.FormatUint(uint64(c.act), 16)))
	b.WriteString(" TMP:" + fmt.Sprintf("%-2s", strconv.FormatUint(uint64(c.tmp), 16)))
	b.WriteString(" FLG:" + fmt.Sprintf("%08b", c.flags))

	return b.String()
}

func (c *CPU) Step() {
	// Fetch instruction
	opcode := c.ReadMem(c.pc)
	inst := ""

	lo := opcode & 0x0F
	hi := opcode >> 4

	prevSC, prevPC := c.sc, c.pc

	switch opcode {
	case 0x00, 0x10, 0x20, 0x30, 0x08, 0x18, 0x28, 0x38:
		inst = "NOP"
		c.sc += 4
		c.pc++

	case 0x01, 0x11, 0x21, 0x31:
		value := c.nextU16()
		switch hi {
		case 0:
			inst = "LXI BC, " + strconv.FormatUint(uint64(value), 16)
			c.b = uint8(value >> 8)
			c.c = uint8(value)
		case 1:
			inst = "LXI DE, " + strconv.FormatUint(uint64(value), 16)
			c.d = uint8(value >> 8)
			c.e = uint8(value)
		case 2:
			inst = "LXI HL, " + strconv.FormatUint(uint64(value), 16)
			c.h = uint8(value >> 8)
			c.l = uint8(value)
		case 3:
			inst = "LXI SP, " + strconv.FormatUint(uint64(value), 16)
			c.sp = value
		}
		c.sc += 10
		c.pc += 3

	case 0x0A, 0x1A:
		var value uint8

		switch hi {
		case 0:
			value = c.ReadMem(uint16(c.b)<<8 | uint16(c.c))
			inst = "LDAX BC"
		case 1:
			value = c.ReadMem(uint16(c.d)<<8 | uint16(c.e))
			inst = "LDAX DE"
		}

		c.acc = value
		c.sc += 7
		c.pc++

	case 0xC2:
		addr := c.nextU16()
		inst = "JNZ " + strconv.FormatUint(uint64(addr), 16)

		if c.getZ() != 0 {
			c.pc = addr
		} else {
			c.pc += 3
		}

		c.sc += 10

	case 0xC3:
		addr := c.nextU16()
		inst = "JMP " + strconv.FormatUint(uint64(addr), 16)
		c.sc += 10
		c.pc = addr

	case 0xCD, 0xDD, 0xED, 0xFD:
		addr := c.nextU16()
		inst = "CALL " + strconv.FormatUint(uint64(addr), 16)
		c.sc += 17
		c.push(c.pc)
		c.pc = addr

	case 0xC9, 0xD9:
		inst = "RET"
		c.pc = c.pop()
		c.pc++
		c.sc += 10

	case 0x06, 0x16, 0x26, 0x36, 0x0E, 0x1E, 0x2E, 0x3E:
		inst = "MVI "

		value := c.nextU8()

		switch hi {
		case 0:
			if lo == 0x6 {
				inst += "B, " + strconv.FormatUint(uint64(value), 16)
				c.b = value
			} else {
				inst += "C, " + strconv.FormatUint(uint64(value), 16)
				c.c = value
			}
		case 1:
			if lo == 0x6 {
				inst += "D, " + strconv.FormatUint(uint64(value), 16)
				c.d = value
			} else {
				inst += "E, " + strconv.FormatUint(uint64(value), 16)
				c.e = value
			}
		case 2:
			if lo == 0x6 {
				inst += "H, " + strconv.FormatUint(uint64(value), 16)
				c.h = value
			} else {
				inst += "L, " + strconv.FormatUint(uint64(value), 16)
				c.l = value
			}
		case 3:
			if lo == 0x6 {
				inst += "M, " + strconv.FormatUint(uint64(value), 16)
				c.WriteMem(c.getHL(), value)
				c.sc += 3
			} else {
				inst += "A, " + strconv.FormatUint(uint64(value), 16)
				c.acc = value
			}
		}

		c.sc += 7
		c.pc += 2

	case 0x4C, 0x5C, 0x6C, 0x7C:
		inst = "MOV "

		value := c.h
		switch hi {
		case 4:
			inst += "C, "
			c.c = value
		case 5:
			inst += "E, "
			c.e = value
		case 6:
			inst += "L, "
			c.l = value
		case 7:
			inst += "A, "
			c.acc = value
		}

		inst += "H"

		c.sc += 5
		c.pc++

	case 0x4F, 0x5F, 0x6F, 0x7F:
		inst = "MOV "

		value := c.acc
		switch hi {
		case 4:
			inst += "C, "
			c.c = value
		case 5:
			inst += "E, "
			c.e = value
		case 6:
			inst += "L, "
			c.l = value
		case 7:
			inst += "A, "
			c.acc = value
		}

		inst += "A"

		c.sc += 5
		c.pc++

	case 0x47, 0x57, 0x67, 0x77:
		inst = "MOV "

		value := c.acc
		switch hi {
		case 4:
			inst += "B, "
			c.b = value
		case 5:
			inst += "D, "
			c.d = value
		case 6:
			inst += "H, "
			c.h = value
		case 7:
			inst += "M, "
			c.WriteMem(c.getHL(), value)
			c.sc += 2
		}

		inst += "A"

		c.sc += 5
		c.pc++

	case 0x03, 0x13, 0x23, 0x33:
		inst = "INX "

		switch hi {
		case 0:
			inst += "BC"
			c.c++
			if c.c == 0 {
				c.b++
			}
		case 1:
			inst += "DE"
			c.e++
			if c.e == 0 {
				c.d++
			}
		case 2:
			inst += "HL"
			c.l++
			if c.l == 0 {
				c.h++
			}
		case 3:
			inst += "SP"
			c.sp++
		}

		c.sc += 5
		c.pc++

	case 0x05, 0x15, 0x25, 0x35:
		inst = "DCR "
		var value uint8
		switch hi {
		case 0:
			inst += "B"
			c.b--
			value = c.b
		case 1:
			inst += "D"
			c.d--
			value = c.d
		case 2:
			inst += "H"
			c.h--
			value = c.h
		case 3:
			inst += "M"
			value = c.ReadMem(c.getHL()) - 1
			c.WriteMem(c.getHL(), value)

			c.sc += 5
		}

		c.setS(value>>7&1 == 1)
		c.setZ(value == 0)
		c.setA(value&0x0F != 0x0F)
		c.setP(value&1 == 0)

		c.sc += 5
		c.pc++

	case 0x09, 0x19, 0x29, 0x39:
		var res uint32

		inst = "DAD "

		switch hi {
		case 0:
			inst += "BC"
			res = uint32(c.getHL()) + uint32(c.getBC())
		case 1:
			inst += "DE"
			res = uint32(c.getHL()) + uint32(c.getDE())
		case 2:
			inst += "HL"
			res = uint32(c.getHL()) + uint32(c.getHL())
		case 3:
			inst += "SP"
			res = uint32(c.getHL()) + uint32(c.sp)
		}

		c.setC(res > math.MaxUint16)
		c.setHL(uint16(res))

		c.sc += 10
		c.pc++

	case 0xE6:
		value := c.nextU8()
		inst = "ANI " + strconv.FormatUint(uint64(value), 16)

		res := c.acc & value

		c.setS(res>>7&1 == 1)
		c.setZ(res == 0)
		// c.setA(c.acc&0x0F+value&0x0F > 0x0F)
		c.setA(false)
		c.setP(res&1 == 0)
		// c.setC(uint16(c.acc)+uint16(value) > 0xFF)
		c.setC(false)

		c.acc = res

		c.sc += 7
		c.pc += 2

	case 0xFE:
		value := c.nextU8()
		inst = "CPI " + strconv.FormatUint(uint64(value), 16)
		res := c.acc - value

		c.setS(res>>7&1 == 1)
		c.setZ(res == 0)
		c.setA(c.acc&0x0F < value&0x0F)
		c.setP(res&1 == 0)
		c.setC(c.acc < value)

		c.sc += 7
		c.pc += 2

	case 0xC5, 0xD5, 0xE5, 0xF5:
		inst = "PUSH "

		switch hi {
		case 0xC:
			inst += "BC"
			c.push(c.getBC())
		case 0xD:
			inst += "DE"
			c.push(c.getDE())
		case 0xE:
			inst += "HL"
			c.push(c.getHL())
		case 0xF:
			inst += "PSW"
			c.push(c.getPSW())
		}

		c.sc += 11
		c.pc++

	case 0xC1, 0xD1, 0xE1, 0xF1:
		inst = "POP "

		switch hi {
		case 0xC:
			inst += "BC"
			c.setBC(c.pop())
		case 0xD:
			inst += "DE"
			c.setDE(c.pop())
		case 0xE:
			inst += "HL"
			c.setHL(c.pop())
		case 0xF:
			inst += "PSW"
			c.setPSW(c.pop())
		}

		c.sc += 10
		c.pc++

	case 0xEB:
		inst = "XCHG HL, DE"
		c.h, c.l, c.d, c.e = c.d, c.e, c.h, c.l
		c.pc++
		c.sc += 4

	default:
		log.Fatalf("unimplemented opcode %02x", opcode)
	}

	c.inst = inst

	log.Println(c)

	if prevPC == c.pc {
		log.Fatal("program counter not incremented")
	}

	if prevSC == c.sc {
		log.Fatal("state counter not incremented")
	}
}

func (c *CPU) nextU8() uint8 {
	return c.ReadMem(c.pc + 1)
}

func (c *CPU) nextU16() uint16 {
	return uint16(c.ReadMem(c.pc+2))<<8 | uint16(c.ReadMem(c.pc+1))
}

func (c *CPU) push(value uint16) {
	c.sp -= 2
	c.WriteMem(c.sp, uint8(value))
	c.WriteMem(c.sp+1, uint8(value>>8))
}

func (c *CPU) pop() uint16 {
	value := uint16(c.ReadMem(c.sp+1))<<8 | uint16(c.ReadMem(c.sp))
	c.sp += 2

	return value
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

func (c *CPU) getPSW() uint16 {
	return uint16(c.acc)<<8 | uint16(c.flags)
}

func (c *CPU) setBC(bc uint16) {
	c.b = uint8(bc >> 8)
	c.c = uint8(bc)
}

func (c *CPU) setDE(de uint16) {
	c.b = uint8(de >> 8)
	c.c = uint8(de)
}

func (c *CPU) setHL(hl uint16) {
	c.b = uint8(hl >> 8)
	c.c = uint8(hl)
}

func (c *CPU) setPSW(hl uint16) {
	c.acc = uint8(hl >> 8)
	c.flags = uint8(hl)
}

func (c *CPU) getS() uint8 {
	return c.flags >> 7 & 1
}

func (c *CPU) getZ() uint8 {
	return c.flags >> 6 & 1
}

func (c *CPU) getA() uint8 {
	return c.flags >> 4 & 1
}

func (c *CPU) getP() uint8 {
	return c.flags >> 2 & 1
}

func (c *CPU) getC() uint8 {
	return c.flags & 1
}

func (c *CPU) setS(b bool) {
	if b {
		c.flags = c.flags | 0x80
	} else {
		c.flags = c.flags & 0x7F
	}
}

func (c *CPU) setZ(b bool) {
	if b {
		c.flags = c.flags | 0x40
	} else {
		c.flags = c.flags & 0xBF
	}
}

func (c *CPU) setA(b bool) {
	if b {
		c.flags = c.flags | 0x10
	} else {
		c.flags = c.flags & 0xEF
	}
}

func (c *CPU) setP(b bool) {
	if b {
		c.flags = c.flags | 0x04
	} else {
		c.flags = c.flags & 0xFB
	}
}

func (c *CPU) setC(b bool) {
	if b {
		c.flags = c.flags | 0x01
	} else {
		c.flags = c.flags & 0xFE
	}
}
