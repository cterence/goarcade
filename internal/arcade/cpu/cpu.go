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
	b.WriteString(" S:" + fmt.Sprintf("%01b", c.getS()))
	b.WriteString(" Z:" + fmt.Sprintf("%01b", c.getZ()))
	b.WriteString(" AC:" + fmt.Sprintf("%01b", c.getAC()))
	b.WriteString(" P:" + fmt.Sprintf("%01b", c.getP()))
	b.WriteString(" CY:" + fmt.Sprintf("%01b", c.getCY()))

	return b.String()
}

func (c *CPU) Step() uint64 {
	// Fetch instruction
	inst, op1, op2, instLength := c.DecodeInst()
	prevSC, prevPC, states := c.sc, c.pc, uint64(0)

	// convenience maps
	regMap := map[string]*uint8{
		"A": &c.a,
		"F": &c.f,
		"B": &c.b,
		"C": &c.c,
		"D": &c.d,
		"E": &c.e,
		"H": &c.h,
		"L": &c.l,
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
		states += 4

	case "LXI":
		value := uint16(lib.Must(strconv.ParseUint(op2, 16, 16)))
		setDoubleRegMap[op1](value)

		states += 10

	case "STA":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		c.WriteMem(addr, c.a)
		states += 13

	case "LDA":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		c.a = c.ReadMem(addr)
		states += 13

	case "LDAX":
		c.a = c.ReadMem(getDoubleRegMap[op1]())
		states += 7

	case "JNZ":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		states += 10

		if c.getZ() != 0 {
			c.pc = addr

			return states
		}

	case "JMP":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		states += 10
		c.pc = addr

		return states

	case "CALL":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		states += 17
		c.push(c.pc + 2)
		c.pc = addr

		return states

	case "RET":
		c.pc = c.pop()
		states += 10

		return states

	case "RP":
		states += 5
		if c.getS() == 0 {
			states += 6
			c.pc = c.pop()

			return states
		}

	case "MVI":
		value := uint8(lib.Must(strconv.ParseUint(op2, 16, 8)))
		if op1 == "M" {
			c.WriteMem(c.getHL(), value)
			states += 3
		} else {
			*regMap[op1] = value
		}

		states += 7

	case "MOV":
		var src uint8
		if op2 == "M" {
			src = c.ReadMem(c.getHL())
			states += 2
		} else {
			src = *regMap[op2]
		}

		if op1 == "M" {
			c.WriteMem(c.getHL(), src)
			states += 2
		} else {
			*regMap[op1] = src
		}
		states += 5

	case "INX":
		v := getDoubleRegMap[op1]() + 1
		setDoubleRegMap[op1](v)
		states += 5

	case "DCR":
		var value uint8

		if op1 == "M" {
			value = c.ReadMem(c.getHL())
			c.WriteMem(c.getHL(), value)
			states += 5
		} else {
			value = *regMap[op1]
		}

		c.setS(value>>7&1 == 1)
		c.setZ(value == 0)
		c.setAC(value&0x0F != 0x0F)
		c.setP(value&1 == 0)

		states += 5

	case "DAD":
		res := uint32(c.getHL()) + uint32(getDoubleRegMap[op1]())

		c.setCY(res > math.MaxUint16)
		c.setHL(uint16(res))

		states += 10

	case "ANI":
		value := uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))

		res := c.a & value

		c.setS(res&0x80 == 0x80)
		c.setZ(res == 0)
		c.setAC(false)
		c.setP(res&1 == 0)
		c.setCY(false)

		c.a = res

		states += 7

	case "CPI":
		value := uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		res := c.a - value

		c.setS(res&0x80 == 0x80)
		c.setZ(res == 0)
		c.setAC(c.a&0x0F < value&0x0F)
		c.setP(res&1 == 0)
		c.setCY(c.a < value)

		states += 7

	case "PUSH":
		c.push(getDoubleRegMap[op1]())
		states += 11

	case "POP":
		setDoubleRegMap[op1](c.pop())
		states += 10

	case "XCHG":
		c.h, c.l, c.d, c.e = c.d, c.e, c.h, c.l
		states += 4

	case "OUT":
		value := uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		c.WriteMem(uint16(value)<<8|uint16(value), c.a)

		states += 10

	case "RST":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		states += 11
		c.pc = addr

		return states

	case "RRC":
		c.setCY(c.a&0x1 == 1)
		c.a = c.a >> 1
		states += 4

	case "ADI":
		value := uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))

		res := c.a + value

		c.setS(res&0x80 == 0x80)
		c.setZ(res == 0)
		c.setAC(c.a&0x0F+value&0x0F > 0x0F)
		c.setP(res&1 == 0)
		c.setCY(uint16(c.a)+uint16(value) > 0xFF)

		c.a = res
		states += 7

	case "XRA":
		var value uint8

		if op1 == "M" {
			value = c.ReadMem(c.getHL())
		} else {
			value = *regMap[op1]
		}

		res := c.a ^ value

		c.setS(res&0x80 == 0x80)
		c.setZ(res == 0)
		c.setAC(false)
		c.setP(res&1 == 0)
		c.setCY(false)

		c.a = res
		states += 4

	case "ANA":
		var value uint8

		if op1 == "M" {
			value = c.ReadMem(c.getHL())
		} else {
			value = *regMap[op1]
		}

		res := c.a & value

		c.setS(res&0x80 == 0x80)
		c.setZ(res == 0)
		c.setAC(false)
		c.setP(res&1 == 0)
		c.setCY(false)

		c.a = res
		states += 4

	case "ORA":
		var value uint8

		if op1 == "M" {
			value = c.ReadMem(c.getHL())
		} else {
			value = *regMap[op1]
		}

		res := c.a | value

		c.setS(res&0x80 == 0x80)
		c.setZ(res == 0)
		c.setAC(false)
		c.setP(res&1 == 0)
		c.setCY(false)

		c.a = res
		states += 4

	case "EI":
		c.interrupts = true
		states += 4
		c.pc++

	case "DI":
		c.interrupts = true
		states += 4
		c.pc++

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

func (c *CPU) getBC() uint16 {
	return uint16(c.b)<<8 | uint16(c.c)
}

func (c *CPU) getDE() uint16 {
	return uint16(c.d)<<8 | uint16(c.e)
}

func (c *CPU) getHL() uint16 {
	return uint16(c.h)<<8 | uint16(c.l)
}

func (c *CPU) getAF() uint16 {
	return uint16(c.a)<<8 | uint16(c.f)
}

func (c *CPU) getSP() uint16 {
	return c.sp
}

func (c *CPU) setBC(v uint16) {
	c.b = uint8(v >> 8)
	c.c = uint8(v)
}

func (c *CPU) setDE(v uint16) {
	c.b = uint8(v >> 8)
	c.c = uint8(v)
}

func (c *CPU) setHL(v uint16) {
	c.b = uint8(v >> 8)
	c.c = uint8(v)
}

func (c *CPU) setAF(v uint16) {
	c.a = uint8(v >> 8)
	c.f = uint8(v)
}

func (c *CPU) setSP(v uint16) {
	c.sp = v
}

func (c *CPU) getS() uint8 {
	return c.f >> 7 & 1
}

func (c *CPU) getZ() uint8 {
	return c.f >> 6 & 1
}

func (c *CPU) getAC() uint8 {
	return c.f >> 4 & 1
}

func (c *CPU) getP() uint8 {
	return c.f >> 2 & 1
}

func (c *CPU) getCY() uint8 {
	return c.f & 1
}

func (c *CPU) setS(b bool) {
	if b {
		c.f = c.f | 0x80
	} else {
		c.f = c.f & 0x7F
	}
}

func (c *CPU) setZ(b bool) {
	if b {
		c.f = c.f | 0x40
	} else {
		c.f = c.f & 0xBF
	}
}

func (c *CPU) setAC(b bool) {
	if b {
		c.f = c.f | 0x10
	} else {
		c.f = c.f & 0xEF
	}
}

func (c *CPU) setP(b bool) {
	if b {
		c.f = c.f | 0x04
	} else {
		c.f = c.f & 0xFB
	}
}

func (c *CPU) setCY(b bool) {
	if b {
		c.f = c.f | 0x01
	} else {
		c.f = c.f & 0xFE
	}
}
