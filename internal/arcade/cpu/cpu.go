package cpu

import (
	"log"
	"strconv"
	"strings"
)

type CPU struct {
	ReadMem  func(uint16) uint8
	WriteMem func(uint16, uint8)

	Running bool

	cycles uint64

	pc uint16
	sp uint16
	b  uint8
	c  uint8
	d  uint8
	e  uint8
	h  uint8
	l  uint8
	w  uint8
	z  uint8

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
	c.w = 0
	c.z = 0
	c.acc = 0
	c.act = 0
	c.flags = 0
	c.tmp = 0
}

func (c *CPU) String() string {
	var b strings.Builder

	b.WriteString("CYCLES:" + strconv.FormatUint(uint64(c.cycles), 10))
	b.WriteString(" PC:" + strconv.FormatUint(uint64(c.pc), 16))
	b.WriteString(" SP:" + strconv.FormatUint(uint64(c.sp), 16))
	b.WriteString(" B:" + strconv.FormatUint(uint64(c.b), 16))
	b.WriteString(" C:" + strconv.FormatUint(uint64(c.c), 16))
	b.WriteString(" D:" + strconv.FormatUint(uint64(c.d), 16))
	b.WriteString(" E:" + strconv.FormatUint(uint64(c.e), 16))
	b.WriteString(" H:" + strconv.FormatUint(uint64(c.h), 16))
	b.WriteString(" L:" + strconv.FormatUint(uint64(c.l), 16))
	b.WriteString(" W:" + strconv.FormatUint(uint64(c.w), 16))
	b.WriteString(" Z:" + strconv.FormatUint(uint64(c.z), 16))
	b.WriteString(" ACC:" + strconv.FormatUint(uint64(c.acc), 16))
	b.WriteString(" ACT:" + strconv.FormatUint(uint64(c.act), 16))
	b.WriteString(" FLG:" + strconv.FormatUint(uint64(c.flags), 16))
	b.WriteString(" TMP:" + strconv.FormatUint(uint64(c.tmp), 16))

	return b.String()
}

func (c *CPU) Step() {
	// Fetch instruction
	opcode := c.ReadMem(c.pc)
	inst := ""

	// lo := opcode & 0x0F
	hi := opcode >> 4

	switch opcode {
	case 0x00, 0x10, 0x20, 0x30:
		inst = "NOP"
		c.cycles += 4
		c.pc++

	case 0x06:
		value := c.nextU8()
		inst = "MVI B, " + strconv.FormatUint(uint64(value), 16)
		c.b = value
		c.cycles += 7
		c.pc += 2

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
		c.cycles += 10
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
		c.cycles += 7
		c.pc++

	case 0xC3:
		addr := c.nextU16()
		inst = "JMP " + strconv.FormatUint(uint64(addr), 16)
		c.cycles += 10
		c.pc = addr

	case 0xCD, 0xDD, 0xED, 0xFD:
		addr := c.nextU16()
		inst = "CALL " + strconv.FormatUint(uint64(addr), 16)
		c.cycles += 17
		c.push(c.pc)
		c.pc = addr

	case 0x47, 0x57, 0x67, 0x77:
		inst = "MOV "

		value := c.acc
		switch hi {
		case 4:
			inst += "B, A"
			c.b = value
		case 5:
			inst += "D, A"
			c.d = value
		case 6:
			inst += "H, A"
			c.h = value
		case 7:
			inst += "M, A"
			c.WriteMem(uint16(c.h)<<8|uint16(c.l), value)
			c.cycles += 2
		}

		c.cycles += 5
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

		c.cycles += 5
		c.pc++

	default:
		panic("unimplemented opcode " + strconv.FormatUint(uint64(opcode), 16))
	}

	log.Println(inst, "-", c)
}

func (c *CPU) nextU8() uint8 {
	return c.ReadMem(c.pc + 1)
}

func (c *CPU) nextU16() uint16 {
	return uint16(c.ReadMem(c.pc+2))<<8 | uint16(c.ReadMem(c.pc+1))
}

func (c *CPU) push(value uint16) {
	c.WriteMem(c.sp, uint8(value))
	c.WriteMem(c.sp+1, uint8(value>>8))
	c.sp += 2
}

func (c *CPU) pop() uint16 {
	c.sp -= 2

	return uint16(c.ReadMem(c.sp+1))<<8 | uint16(c.ReadMem(c.sp))
}
