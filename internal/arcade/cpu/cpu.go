package cpu

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/cterence/space-invaders/internal/arcade/lib"
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

func (c *CPU) Init(options ...Option) {
	c.Running = true
	c.pc = 0
	c.sp = 0
	c.b = 0
	c.c = 0
	c.d = 0
	c.e = 0
	c.h = 0
	c.l = 0
	c.a = 0
	c.f = 0
	c.interrupts = false

	for _, o := range options {
		o(c)
	}
}

func (c *CPU) String() string {
	var b strings.Builder

	b.WriteString(" PC:" + fmt.Sprintf("%-4X", c.pc))
	b.WriteString(" SC:" + fmt.Sprintf("%-5d", c.sc))
	b.WriteString(" SP:" + fmt.Sprintf("%-4X", c.sp))
	b.WriteString(" A:" + fmt.Sprintf("%-2X", c.a))
	b.WriteString(" B:" + fmt.Sprintf("%-2X", c.b))
	b.WriteString(" C:" + fmt.Sprintf("%-2X", c.c))
	b.WriteString(" D:" + fmt.Sprintf("%-2X", c.d))
	b.WriteString(" E:" + fmt.Sprintf("%-2X", c.e))
	b.WriteString(" H:" + fmt.Sprintf("%-2X", c.h))
	b.WriteString(" L:" + fmt.Sprintf("%-2X", c.l))
	b.WriteString(" S:" + fmt.Sprintf("%01b", c.getSF()))
	b.WriteString(" Z:" + fmt.Sprintf("%01b", c.getZF()))
	b.WriteString(" AC:" + fmt.Sprintf("%01b", c.getACF()))
	b.WriteString(" P:" + fmt.Sprintf("%01b", c.getPF()))
	b.WriteString(" CY:" + fmt.Sprintf("%01b", c.getCYF()))

	return b.String()
}

func (c *CPU) Step() uint64 {
	// Fetch instruction
	inst, op1, op2, instLength, states := c.DecodeInst()
	prevSC, prevPC := c.sc, c.pc

	// convenience maps
	getOpMap := map[string]func() uint8{
		"A": c.getA,
		"F": c.getF,
		"B": c.getB,
		"C": c.getC,
		"D": c.getD,
		"E": c.getE,
		"H": c.getH,
		"L": c.getL,
		"M": c.getM,
	}

	setOpMap := map[string]func(uint8){
		"A": c.setA,
		"F": c.setF,
		"B": c.setB,
		"C": c.setC,
		"D": c.setD,
		"E": c.setE,
		"H": c.setH,
		"L": c.setL,
		"M": c.setM,
	}

	getDoubleRegMap := map[string]func() uint16{
		"AF": c.getAF,
		"BC": c.getBC,
		"DE": c.getDE,
		"HL": c.getHL,
		"SP": c.getSP,
	}

	setDoubleRegMap := map[string]func(uint16){
		"AF": c.setAF,
		"BC": c.setBC,
		"DE": c.setDE,
		"HL": c.setHL,
		"SP": c.setSP,
	}

	switch inst {
	case "NOP":

	case "LXI":
		value := uint16(lib.Must(strconv.ParseUint(op2, 16, 16)))
		setDoubleRegMap[op1](value)

	case "STA":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		c.WriteMem(addr, c.a)

	case "LDA":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		c.a = c.ReadMem(addr)

	case "LDAX":
		c.a = c.ReadMem(getDoubleRegMap[op1]())

	case "JNZ":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getZF() != 0 {
			c.pc = addr

			return states
		}

	case "JMP":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		c.pc = addr

		return states

	case "CALL":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		c.push(c.pc + 2)
		c.pc = addr

		return states

	case "RET":
		c.pc = c.pop()

		return states

	case "RP":
		if c.getSF() == 0 {
			c.pc = c.pop()

			return states
		}

	case "MVI":
		value := uint8(lib.Must(strconv.ParseUint(op2, 16, 8)))
		setOpMap[op1](value)

	case "MOV":
		setOpMap[op1](getOpMap[op2]())

	case "INX":
		v := getDoubleRegMap[op1]() + 1
		setDoubleRegMap[op1](v)

	case "DCR":
		value := getOpMap[op1]() - 1
		setOpMap[op1](value)

		c.setSF(value>>7&1 == 1)
		c.setZF(value == 0)
		c.setACF(value&0x0F != 0x0F)
		c.setPF(value&1 == 0)

	case "DAD":
		res := uint32(c.getHL()) + uint32(getDoubleRegMap[op1]())

		c.setCYF(res > math.MaxUint16)
		c.setHL(uint16(res))

	case "ANI":
		value := uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		res := c.a & value

		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(false)
		c.setPF(res&1 == 0)
		c.setCYF(false)

		c.a = res

	case "CPI":
		value := uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		res := c.a - value

		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(c.a&0x0F < value&0x0F)
		c.setPF(res&1 == 0)
		c.setCYF(c.a < value)

	case "PUSH":
		c.push(getDoubleRegMap[op1]())

	case "POP":
		setDoubleRegMap[op1](c.pop())

	case "XCHG":
		c.h, c.l, c.d, c.e = c.d, c.e, c.h, c.l

	case "OUT":
		value := uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		c.WriteMem(uint16(value)<<8|uint16(value), c.a)

	case "RST":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		c.pc = addr

		return states

	case "RRC":
		c.setCYF(c.a&0x1 == 1)
		c.a = c.a >> 1

	case "ADI":
		value := uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))

		res := c.a + value

		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(c.a&0x0F+value&0x0F > 0x0F)
		c.setPF(res&1 == 0)
		c.setCYF(uint16(c.a)+uint16(value) > 0xFF)

		c.a = res

	case "XRA":
		value := getOpMap[op1]()
		res := c.a ^ value

		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(false)
		c.setPF(res&1 == 0)
		c.setCYF(false)

		c.a = res

	case "ANA":
		value := getOpMap[op1]()
		res := c.a & value

		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(false)
		c.setPF(res&1 == 0)
		c.setCYF(false)

		c.a = res

	case "ORA":
		value := getOpMap[op1]()
		res := c.a | value

		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(false)
		c.setPF(res&1 == 0)
		c.setCYF(false)

		c.a = res

	case "EI":
		c.interrupts = true

	case "DI":
		c.interrupts = true

	default:
		log.Fatalf("unimplemented inst %s", inst)
	}

	c.pc += instLength

	if prevPC == c.pc {
		log.Fatal("program counter not incremented")
	}

	c.sc += uint64(states)

	if prevSC >= c.sc {
		log.Fatal("state counter not incremented")
	}

	return states
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
