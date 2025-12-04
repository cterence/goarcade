package cpu

import (
	"fmt"
	"log"
	"math"
	"math/bits"
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

	b.WriteString("PC: " + fmt.Sprintf("%04X", c.pc))
	b.WriteString(", AF: " + fmt.Sprintf("%04X", c.getAF()))
	b.WriteString(", BC: " + fmt.Sprintf("%04X", c.getBC()))
	b.WriteString(", DE: " + fmt.Sprintf("%04X", c.getDE()))
	b.WriteString(", HL: " + fmt.Sprintf("%04X", c.getHL()))
	b.WriteString(", SP: " + fmt.Sprintf("%04X", c.sp))
	b.WriteString(", CYC: " + fmt.Sprintf("%-6d", c.sc))

	return b.String()
}

func (c *CPU) Step() uint64 {
	// Fetch instruction
	inst, op1, op2, instLength, states := c.DecodeInst()
	prevSC, prevPC := c.sc, c.pc

	if c.debug {
		// fmt.Printf("%s (%02X %02X %02X %02X) %-13s\n", c, c.ReadMem(c.pc), c.ReadMem(c.pc+1), c.ReadMem(c.pc+2), c.ReadMem(c.pc+3), inst+" "+op1+" "+op2)
		fmt.Printf("%s (%02X %02X %02X %02X)\n", c, c.ReadMem(c.pc), c.ReadMem(c.pc+1), c.ReadMem(c.pc+2), c.ReadMem(c.pc+3))
	}

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
			c.sc += states

			return states
		}

	case "JPO":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getPF() == 0 {
			c.pc = addr
			c.sc += states

			return states
		}

	case "JPE":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getPF() == 1 {
			c.pc = addr
			c.sc += states

			return states
		}

	case "JZ":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getZF() == 1 {
			c.pc = addr
			c.sc += states

			return states
		}

	case "JNC":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getCYF() == 0 {
			c.pc = addr
			c.sc += states

			return states
		}

	case "JC":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getCYF() == 1 {
			c.pc = addr
			c.sc += states

			return states
		}

	case "JP":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getSF() == 0 {
			c.pc = addr
			c.sc += states

			return states
		}

	case "JM":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getSF() == 1 {
			c.pc = addr
			c.sc += states

			return states
		}

	case "JMP":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		c.pc = addr
		c.sc += states

		return states

	case "CALL":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		c.push(c.pc + 3)
		c.pc = addr
		c.sc += states

		return states

	case "CNC":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getCYF() == 0 {
			c.push(c.pc + 3)
			c.pc = addr
			c.sc += states

			return states
		}

	case "CC":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getCYF() == 1 {
			c.push(c.pc + 3)
			c.pc = addr
			c.sc += states

			return states
		}

	case "CP":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getSF() == 0 {
			c.push(c.pc + 3)
			c.pc = addr
			c.sc += states

			return states
		}

	case "CM":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getSF() == 1 {
			c.push(c.pc + 3)
			c.pc = addr
			c.sc += states

			return states
		}

	case "CPO":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getPF() == 0 {
			c.push(c.pc + 3)
			c.pc = addr
			c.sc += states

			return states
		}

	case "CPE":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getPF() == 1 {
			c.push(c.pc + 3)
			c.pc = addr
			c.sc += states

			return states
		}

	case "CNZ":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getZF() == 0 {
			c.push(c.pc + 3)
			c.pc = addr
			c.sc += states

			return states
		}

	case "CZ":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

		if c.getZF() == 1 {
			c.push(c.pc + 3)
			c.pc = addr
			c.sc += states

			return states
		}

	case "RET":
		c.pc = c.pop()
		c.sc += states

		return states

	case "RP":
		if c.getSF() == 0 {
			c.pc = c.pop()
			c.sc += states

			return states
		}

	case "RM":
		if c.getSF() == 1 {
			c.pc = c.pop()
			c.sc += states

			return states
		}

	case "RNC":
		if c.getCYF() == 0 {
			c.pc = c.pop()
			c.sc += states

			return states
		}

	case "RC":
		if c.getCYF() == 1 {
			c.pc = c.pop()
			c.sc += states

			return states
		}

	case "RNZ":
		if c.getZF() == 0 {
			c.pc = c.pop()
			c.sc += states

			return states
		}

	case "RZ":
		if c.getZF() == 1 {
			c.pc = c.pop()
			c.sc += states

			return states
		}

	case "RPE":
		if c.getPF() == 1 {
			c.pc = c.pop()
			c.sc += states

			return states
		}

	case "RPO":
		if c.getPF() == 0 {
			c.pc = c.pop()
			c.sc += states

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
		c.setPF(bits.OnesCount8(value)%2 == 0)

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
		c.setACF(c.a&0x0F >= value&0x0F)
		c.setPF(bits.OnesCount8(res)%2 == 0)
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
		lo := c.ReadMem(c.sp)
		hi := c.ReadMem(c.sp + 1)
		c.WriteMem(c.sp, c.l)
		c.WriteMem(c.sp+1, c.h)
		c.l = lo
		c.h = hi

	case "SPHL":
		c.sp = c.getHL()

	case "PCHL":
		c.pc = c.getHL()
		c.sc += states

		return states

	case "IN":
		value := uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		c.a = c.portIn(value)

	case "OUT":
		value := uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		c.portOut(value)

	case "RST":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		c.pc = addr
		c.sc += states

		return states

	case "RRC":
		sb := c.a & 0x1
		c.setCYF(sb == 1)
		c.a = c.a>>1 | sb<<7

	case "RLC":
		sb := c.a & 0x80 >> 7
		c.setCYF(sb == 1)
		c.a = c.a<<1 | sb

	case "RAR":
		sb := c.a & 0x1
		c.a = c.a>>1 | c.getCYF()<<7
		c.setCYF(sb == 1)

	case "RAL":
		sb := c.a & 0x80 >> 7
		c.a = c.a<<1 | c.getCYF()
		c.setCYF(sb == 1)

	case "ADD", "ADC", "ADI", "ACI":
		var value uint8

		if inst == "ADD" || inst == "ADC" {
			value = getOpMap[op1]()
		} else {
			value = uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		}

		carry := uint8(0)

		if inst == "ADC" || inst == "ACI" {
			carry = c.getCYF()
		}

		res := c.a + value + carry

		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(c.a&0x0F+(value+carry)&0x0F > 0x0F)
		c.setPF(bits.OnesCount8(res)%2 == 0)
		c.setCYF(uint16(c.a)+uint16(value+carry) > 0xFF)

		c.a = res

	case "SUB", "SBB", "SUI", "SBI":
		var value uint8

		if inst == "SUB" || inst == "SBB" {
			value = getOpMap[op1]()
		} else {
			value = uint8(lib.Must(strconv.ParseUint(op1, 16, 8)))
		}

		carry := uint8(0)

		if inst == "SBB" || inst == "SBI" {
			carry = c.getCYF()
		}

		res := c.a - value - carry

		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF((c.a & 0x0F) >= (value&0x0F + carry))
		c.setPF(bits.OnesCount8(res)%2 == 0)
		c.setCYF(uint16(c.a) < uint16(value)+uint16(carry))

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
		c.setPF(bits.OnesCount8(res)%2 == 0)
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
		c.setACF((c.a|value)&0x08 != 0)
		c.setPF(bits.OnesCount8(res)%2 == 0)
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
		c.setPF(bits.OnesCount8(res)%2 == 0)
		c.setCYF(false)

		c.a = res

	case "INR":
		res := getOpMap[op1]() + 1
		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(false)
		c.setPF(bits.OnesCount8(res)%2 == 0)
		setOpMap[op1](res)

	case "DAA":
		cy := c.getCYF() == 1
		value := uint8(0)

		if c.a&0x0F > 0x09 || c.getACF() == 1 {
			value += 0x06
		}

		if (c.a+value)&0xF0 > 0x90 || cy || c.a > 0x99 {
			value += 0x60
			cy = true
		}

		res := c.a + value

		c.setSF(res&0x80 == 0x80)
		c.setZF(res == 0)
		c.setACF(c.a&0x0F+value&0x0F > 0x0F)
		c.setPF(bits.OnesCount8(res)%2 == 0)
		c.setCYF(cy)

		c.a = res

	case "EI":
		c.interrupts = true

	case "DI":
		c.interrupts = false

	case "STAX":
		c.WriteMem(getDoubleRegMap[op1](), c.a)

	case "SHLD":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		c.WriteMem(addr, c.getL())
		c.WriteMem(addr+1, c.getH())

	case "LHLD":
		addr := uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))
		c.setHL(uint16(c.ReadMem(addr+1))<<8 | uint16(c.ReadMem(addr)))

	case "STC":
		c.setCYF(true)

	default:
		log.Fatalf("unimplemented inst %s", inst)
	}

	c.sc += states

	if prevSC >= c.sc {
		log.Fatal("state counter not incremented")
	}

	c.pc += instLength

	if prevPC == c.pc {
		log.Fatal("program counter not incremented")
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

		fmt.Println()
	case 1:
		switch c.c {
		case 2:
			fmt.Print(string(c.e))
		case 9:
			addr := c.getDE()
			for offset := uint16(0); ; offset++ {
				b := c.ReadMem(addr + offset)
				if b == '$' {
					break
				}

				fmt.Print(string(b))
			}

		default:
			fmt.Printf("unimplemented out operation for port 1: %02x", c.c)
		}
	default:
		fmt.Printf("unimplemented out port: %02x", portNumber)
	}
}
