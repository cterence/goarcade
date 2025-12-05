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
	b.WriteString(", CYC: " + strconv.FormatUint(c.sc, 10))

	return b.String()
}

func (c *CPU) Step() uint8 {
	// Fetch instruction
	inst, instLength, op1, op2, imm1, imm2, states := c.DecodeInst()
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	prevSC, prevPC := c.sc, c.pc

	if c.debug {
		// fmt.Printf("%s (%02X %02X %02X %02X) %-13s\n", c, c.ReadMem(c.pc), c.ReadMem(c.pc+1), c.ReadMem(c.pc+2), c.ReadMem(c.pc+3), inst+" "+op1+" "+op2)
		fmt.Printf("%s (%02X %02X %02X %02X)\n", c, c.ReadMem(c.pc), c.ReadMem(c.pc+1), c.ReadMem(c.pc+2), c.ReadMem(c.pc+3))
	}

	switch inst {
	case "NOP":

	case "LXI":
		c.setDoubleOp(op1, imm16)

	case "STA":
		c.WriteMem(imm16, c.a)

	case "LDA":
		c.a = c.ReadMem(imm16)

	case "LDAX":
		c.a = c.ReadMem(c.getDoubleOp(op1))

	case "JMP":
		c.jump(imm16, true)

	case "JP":
		c.jump(imm16, c.getSF() == 0)

	case "JM":
		c.jump(imm16, c.getSF() == 1)

	case "JNZ":
		c.jump(imm16, c.getZF() == 0)

	case "JZ":
		c.jump(imm16, c.getZF() == 1)

	case "JPO":
		c.jump(imm16, c.getPF() == 0)

	case "JPE":
		c.jump(imm16, c.getPF() == 1)

	case "JNC":
		c.jump(imm16, c.getCYF() == 0)

	case "JC":
		c.jump(imm16, c.getCYF() == 1)

	case "CALL":
		c.call(imm16)

	case "CNC":
		c.callCond(imm16, c.getCYF() == 0)

	case "CC":
		c.callCond(imm16, c.getCYF() == 1)

	case "CP":
		c.callCond(imm16, c.getSF() == 0)

	case "CM":
		c.callCond(imm16, c.getSF() == 1)

	case "CPO":
		c.callCond(imm16, c.getPF() == 0)

	case "CPE":
		c.callCond(imm16, c.getPF() == 1)

	case "CNZ":
		c.callCond(imm16, c.getZF() == 0)

	case "CZ":
		c.callCond(imm16, c.getZF() == 1)

	case "RET":
		c.ret()

	case "RP":
		c.retCond(c.getSF() == 0)

	case "RM":
		c.retCond(c.getSF() == 1)

	case "RNZ":
		c.retCond(c.getZF() == 0)

	case "RZ":
		c.retCond(c.getZF() == 1)

	case "RPE":
		c.retCond(c.getPF() == 1)

	case "RPO":
		c.retCond(c.getPF() == 0)

	case "RNC":
		c.retCond(c.getCYF() == 0)

	case "RC":
		c.retCond(c.getCYF() == 1)

	case "MVI":
		c.setOp(op1, imm1)

	case "MOV":
		c.setOp(op1, c.getOp(op2))

	case "INX":
		c.setDoubleOp(op1, c.getDoubleOp(op1)+1)

	case "DCX":
		c.setDoubleOp(op1, c.getDoubleOp(op1)-1)

	case "DAD":
		res := uint32(c.getHL()) + uint32(c.getDoubleOp(op1))

		c.setCYF(res > math.MaxUint16)
		c.setHL(uint16(res))
	case "CMA":
		c.a = 0xFF - c.a

	case "CMC":
		c.setCYF(c.getCYF() == 0)

	case "PUSH":
		c.push(c.getDoubleOp(op1))

	case "POP":
		c.setDoubleOp(op1, c.pop())

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

	case "IN":
		c.a = c.portIn(imm1)

	case "OUT":
		c.portOut(imm1)

	case "RST":
		c.pc = uint16(lib.Must(strconv.ParseUint(op1, 16, 16)))

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

	case "CMP", "CPI":
		var value uint8

		if inst == "CMP" {
			value = c.getOp(op1)
		} else {
			value = imm1
		}

		res := c.a - value

		c.setFlags(res&0x80 == 0x80, res == 0, c.a&0x0F >= value&0x0F, bits.OnesCount8(res)%2 == 0, c.a < value)

	case "ADD", "ADC", "ADI", "ACI":
		var value uint8

		if inst == "ADD" || inst == "ADC" {
			value = c.getOp(op1)
		} else {
			value = imm1
		}

		carry := uint8(0)

		if inst == "ADC" || inst == "ACI" {
			carry = c.getCYF()
		}

		res := c.a + value + carry

		c.setFlags(res&0x80 == 0x80, res == 0, c.a&0x0F+(value+carry)&0x0F > 0x0F, bits.OnesCount8(res)%2 == 0, uint16(c.a)+uint16(value)+uint16(carry) > 0xFF)

		c.a = res

	case "SUB", "SBB", "SUI", "SBI":
		var value uint8

		if inst == "SUB" || inst == "SBB" {
			value = c.getOp(op1)
		} else {
			value = imm1
		}

		carry := uint8(0)

		if inst == "SBB" || inst == "SBI" {
			carry = c.getCYF()
		}

		res := c.a - value - carry

		c.setFlags(res&0x80 == 0x80, res == 0, (c.a&0x0F) >= (value&0x0F+carry), bits.OnesCount8(res)%2 == 0, uint16(c.a) < uint16(value)+uint16(carry))

		c.a = res

	case "XRA", "XRI":
		var value uint8

		if inst == "XRA" {
			value = c.getOp(op1)
		} else {
			value = imm1
		}

		res := c.a ^ value

		c.setFlags(res&0x80 == 0x80, res == 0, false, bits.OnesCount8(res)%2 == 0, false)

		c.a = res

	case "ANA", "ANI":
		var value uint8

		if inst == "ANA" {
			value = c.getOp(op1)
		} else {
			value = imm1
		}

		res := c.a & value

		c.setFlags(res&0x80 == 0x80, res == 0, (c.a|value)&0x08 != 0, bits.OnesCount8(res)%2 == 0, false)

		c.a = res

	case "ORA", "ORI":
		var value uint8

		if inst == "ORA" {
			value = c.getOp(op1)
		} else {
			value = imm1
		}

		res := c.a | value

		c.setFlags(res&0x80 == 0x80, res == 0, false, bits.OnesCount8(res)%2 == 0, false)

		c.a = res

	case "INR":
		value := c.getOp(op1)
		res := value + 1
		c.setOp(op1, res)
		c.setFlags(res&0x80 == 0x80, res == 0, value&0x0F+1 > 0x0F, bits.OnesCount8(res)%2 == 0, c.getCYF() == 1)

	case "DCR":
		value := c.getOp(op1)
		res := value - 1
		c.setOp(op1, res)
		c.setFlags(res&0x80 == 0x80, res == 0, value&0x0F >= 1, bits.OnesCount8(res)%2 == 0, c.getCYF() == 1)

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

		c.setFlags(res&0x80 == 0x80, res == 0, c.a&0x0F+value&0x0F > 0x0F, bits.OnesCount8(res)%2 == 0, cy)

		c.a = res

	case "EI":
		c.interrupts = true

	case "DI":
		c.interrupts = false

	case "STAX":
		c.WriteMem(c.getDoubleOp(op1), c.a)

	case "SHLD":
		c.WriteMem(imm16, c.getL())
		c.WriteMem(imm16+1, c.getH())

	case "LHLD":
		c.setHL(uint16(c.ReadMem(imm16+1))<<8 | uint16(c.ReadMem(imm16)))

	case "STC":
		c.setCYF(true)

	default:
		log.Fatalf("unimplemented inst %s", inst)
	}

	c.sc += uint64(states)

	if prevSC >= c.sc {
		log.Fatal("state counter not incremented")
	}

	if prevPC == c.pc {
		c.pc += instLength
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

func (c *CPU) jump(imm16 uint16, cond bool) {
	if cond {
		c.pc = imm16
	}
}

func (c *CPU) call(imm16 uint16) {
	c.push(c.pc + 3)
	c.pc = imm16
}

func (c *CPU) callCond(imm16 uint16, cond bool) {
	if cond {
		c.call(imm16)
		c.sc += 6
	}
}

func (c *CPU) ret() {
	c.pc = c.pop()
}

func (c *CPU) retCond(cond bool) {
	if cond {
		c.ret()
		c.sc += 6
	}
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
