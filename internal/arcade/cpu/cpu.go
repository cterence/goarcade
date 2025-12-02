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
	c.pc = 0x100
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
	// TODO: functions with switch and params rather than maps
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

		if c.getZF() == 0 {
			c.pc = addr

			return states
		}

	case "JPO":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getPF() == 0 {
			c.pc = addr

			return states
		}

	case "JPE":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getPF() == 1 {
			c.pc = addr

			return states
		}

	case "JZ":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getZF() == 1 {
			c.pc = addr

			return states
		}

	case "JNC":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getCYF() == 0 {
			c.pc = addr

			return states
		}

	case "JC":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getCYF() == 1 {
			c.pc = addr

			return states
		}

	case "JP":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getSF() == 0 {
			c.pc = addr

			return states
		}

	case "JM":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getSF() == 1 {
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

	case "CP":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getZF() == 0 {
			c.push(c.pc + 2)
			c.pc = addr

			return states
		}

	case "RET":
		c.pc = c.pop() + 1

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
		setDoubleRegMap[op1](getDoubleRegMap[op1]() + 1)

	case "DCX":
		setDoubleRegMap[op1](getDoubleRegMap[op1]() - 1)

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

	case "CMP", "CPI":
		var value uint8

		if inst == "CMP" {
			value = getOpMap[op1]()
		} else {
			value = uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		}

		res := c.a - value

		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(c.a&0x0F < value&0x0F)
		c.setPF(res&1 == 0)
		c.setCYF(c.a < value)

	case "CMA":
		c.a = 0xFF - c.a

	case "CMC":
		c.setCYF(c.getCYF() == 0)

	case "PUSH":
		c.push(getDoubleRegMap[op1]())

	case "POP":
		setDoubleRegMap[op1](c.pop())

	case "XCHG":
		c.h, c.l, c.d, c.e = c.d, c.e, c.h, c.l

	case "XTHL":
		sp := c.sp
		c.sp = c.getHL()
		c.setHL(sp)

	case "SPHL":
		c.sp = c.getHL()

	case "PCHL":
		c.pc = c.getHL()

	case "IN":
		value := uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		c.a = c.portIn(value)

	case "OUT":
		value := uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		c.portOut(value)

	case "RST":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		c.pc = addr

		return states

	case "RRC":
		sb := c.a & 0x1
		c.setCYF(sb == 1)
		c.a = c.a>>1 | sb<<7

	case "RLC":
		sb := c.a & 0x80 >> 7
		c.setCYF(sb == 1)
		c.a = c.a<<1 | sb

	case "ADD", "ADC", "ADI", "ACI":
		var value uint8

		if inst == "ADD" || inst == "ADC" {
			value = getOpMap[op1]()
		} else {
			value = uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		}

		res := c.a + value

		if inst == "ADC" || inst == "ACI" {
			res += c.getCYF()
		}

		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(c.a&0x0F+value&0x0F > 0x0F)
		c.setPF(res&1 == 0)
		c.setCYF(uint16(c.a)+uint16(value) > 0xFF)

		c.a = res

	case "SUB", "SBB", "SUI", "SBI":
		var value uint8

		if inst == "SUB" || inst == "SBB" {
			value = getOpMap[op1]()
		} else {
			value = uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		}

		res := c.a - value

		if inst == "SBB" || inst == "SBI" {
			res -= c.getCYF()
		}

		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(c.a&0x0F < value&0x0F)
		c.setPF(res&1 == 0)
		c.setCYF(c.a < value)

		c.a = res

	case "XRA", "XRI":
		var value uint8

		if inst == "XRA" {
			value = getOpMap[op1]()
		} else {
			value = uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		}

		res := c.a ^ value

		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(false)
		c.setPF(res&1 == 0)
		c.setCYF(false)

		c.a = res

	case "ANA", "ANI":
		var value uint8

		if inst == "ANA" {
			value = getOpMap[op1]()
		} else {
			value = uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		}

		res := c.a & value

		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(false)
		c.setPF(res&1 == 0)
		c.setCYF(false)

		c.a = res

	case "ORA", "ORI":
		var value uint8

		if inst == "ORA" {
			value = getOpMap[op1]()
		} else {
			value = uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		}

		res := c.a | value

		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(false)
		c.setPF(res&1 == 0)
		c.setCYF(false)

		c.a = res

	case "INR":
		res := getOpMap[op1]() + 1
		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(false)
		c.setPF(res&1 == 0)
		setOpMap[op1](res)

	case "DAA":
		res := c.a
		if c.a&0x0F > 0x09 || c.getACF() == 1 {
			res += 0x06
		}

		if c.a&0xF0 > 0x90 || c.getCYF() == 1 {
			res += 0x60
		}

		c.a = res

	case "EI":
		c.interrupts = true

	case "DI":
		c.interrupts = true

	case "STAX":
		c.WriteMem(getDoubleRegMap[op1](), c.a)

	case "SHLD":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		c.WriteMem(addr, c.getL())
		c.WriteMem(addr+1, c.getH())

	case "LHLD":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		c.setHL(uint16(c.ReadMem(addr+1))<<8 | uint16(c.ReadMem(addr)))

	default:
		log.Fatalf("unimplemented inst %s", inst)
	}

	c.pc += instLength

	if prevPC == c.pc {
		log.Fatal("program counter not incremented")
	}

	c.sc += states

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

func (c *CPU) portIn(_ uint8) uint8 {
	return 0x0
}

func (c *CPU) portOut(portNumber uint8) {
	switch portNumber {
	case 0:
		c.Running = false
	case 1:
		switch c.c {
		case 2:
			fmt.Printf("%c", c.e)
		case 9:
			var (
				char   byte
				offset uint16
			)

			char = c.ReadMem(c.getDE() + offset)
			for char != '$' {
				fmt.Printf("%c", char)

				offset++
				char = c.ReadMem(c.getDE() + offset)
			}

			fmt.Println()
		default:
			fmt.Printf("unimplemented out operation for port 1: %02x", c.c)
		}
	default:
		fmt.Printf("unimplemented out port: %02x", portNumber)
	}
}
