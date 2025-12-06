package cpu

import (
	"fmt"
	"strconv"
	"strings"
)

type CPU struct {
	ReadMem  func(uint16) uint8
	WriteMem func(uint16, uint8)

	Running bool

	debug bool

	sc uint64
	pc uint16
	sp uint16

	a uint8
	f uint8
	b uint8
	c uint8
	d uint8
	e uint8
	h uint8
	l uint8

	interrupts bool
}

type Option func(*CPU)

func WithDebug(debug bool) Option {
	return func(c *CPU) {
		c.debug = debug
	}
}

func (c *CPU) Init(pc uint16, options ...Option) {
	c.Running = true
	c.pc = pc
	c.sp = 0
	c.b = 0
	c.c = 0
	c.d = 0
	c.e = 0
	c.h = 0
	c.l = 0
	c.a = 0
	c.f = 2
	c.interrupts = false

	for _, o := range options {
		o(c)
	}
}

func (c *CPU) String() string {
	var b strings.Builder

	b.WriteString("PC: " + fmt.Sprintf("%04X", c.pc))
	b.WriteString(", AF: " + fmt.Sprintf("%04X", uint16(c.a)<<8|uint16(c.f)))
	b.WriteString(", BC: " + fmt.Sprintf("%04X", uint16(c.b)<<8|uint16(c.c)))
	b.WriteString(", DE: " + fmt.Sprintf("%04X", uint16(c.d)<<8|uint16(c.e)))
	b.WriteString(", HL: " + fmt.Sprintf("%04X", uint16(c.h)<<8|uint16(c.l)))
	b.WriteString(", SP: " + fmt.Sprintf("%04X", c.sp))
	b.WriteString(", CYC: " + strconv.FormatUint(c.sc, 10))

	return b.String()
}

func (c *CPU) Step() uint8 {
	prevPC := c.pc

	inst := InstByOpcode[c.ReadMem(c.pc)]

	if c.debug {
		fmt.Printf("%s (%02X %02X %02X %02X) %-13s\n", c, c.ReadMem(c.pc), c.ReadMem(c.pc+1), c.ReadMem(c.pc+2), c.ReadMem(c.pc+3), inst.Name+" "+inst.Op1+" "+inst.Op2)
		// fmt.Printf("%s (%02X %02X %02X %02X)\n", c, c.ReadMem(c.pc), c.ReadMem(c.pc+1), c.ReadMem(c.pc+2), c.ReadMem(c.pc+3))
	}

	inst.exec(c, inst.Op1)

	if prevPC == c.pc {
		c.pc += uint16(inst.Length)
	}

	c.sc += uint64(inst.States)

	return inst.States
}
