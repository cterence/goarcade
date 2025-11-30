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

	debug bool

	sc uint64
	pc uint16
	sp uint16

	a uint8
	b uint8
	c uint8
	d uint8
	e uint8
	h uint8
	l uint8

	f uint8
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

func (c *CPU) DecodeInst() (string, string, string, uint16) {
	operandOrder := []string{"B", "C", "D", "E", "H", "L", "M", "A"}
	opcode := c.ReadMem(c.pc)

	hi := opcode >> 4

	var (
		inst   string
		op1    string
		op2    string
		length uint16
	)

	switch opcode {
	default:
		length = 1
		switch opcode {
		case 0x00, 0x10, 0x20, 0x30, 0x08, 0x18, 0x28, 0x38:
			inst = "NOP"

		case 0x02, 0x12:
			inst = "STAX"
			op1 = operandOrder[(opcode-0x2)/0x8]

		case 0x03, 0x13, 0x23, 0x33:
			inst = "INX"
			op1 = operandOrder[(opcode-0x3)/0x8]
			if opcode == 0x33 {
				op1 = "SP"
			}

		case 0x04, 0x14, 0x24, 0x34, 0x0C, 0x1C, 0x2C, 0x3C:
			inst = "INR"
			op1 = operandOrder[(opcode-0x4)/0x8]

		case 0x05, 0x15, 0x25, 0x35, 0x0D, 0x1D, 0x2D, 0x3D:
			inst = "DCR"
			op1 = operandOrder[(opcode-0x4)/0x8]

		case 0x07:
			inst = "RLC"
		case 0x17:
			inst = "RAL"
		case 0x27:
			inst = "DAA"
		case 0x37:
			inst = "STC"

		case 0x09, 0x19, 0x29, 0x39:
			inst = "DAD"
			op1 = operandOrder[(opcode-0x9)/0x8]
			if opcode == 0x39 {
				op1 = "SP"
			}

		case 0x0A, 0x1A:
			inst = "LDAX"
			op1 = operandOrder[(opcode-0xA)/0x8]

		case 0x0B, 0x1B, 0x2B, 0x3B:
			inst = "DCX"
			op1 = operandOrder[(opcode-0xB)/0x8]
			if opcode == 0x3B {
				op1 = "SP"
			}

		case 0x0F:
			inst = "RRC"
		case 0x1F:
			inst = "RAR"
		case 0x2F:
			inst = "CMA"
		case 0x3F:
			inst = "CMC"

		case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F, 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F:
			inst = "MOV"
			op1 = operandOrder[(opcode-0x40)/8]
			op2 = operandOrder[hi%8]

		case 0x76:
			inst = "HALT"

		case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87:
			inst = "ADD"
			op1 = operandOrder[opcode-0x80]

		case 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F:
			inst = "ADC"
			op1 = operandOrder[opcode-0x88]

		case 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97:
			inst = "SUB"
			op1 = operandOrder[opcode-0x90]

		case 0x98, 0x99, 0x9A, 0x9B, 0x9C, 0x9D, 0x9E, 0x9F:
			inst = "SBB"
			op1 = operandOrder[opcode-0x98]

		case 0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7:
			inst = "ANA"
			op1 = operandOrder[opcode-0xA0]

		case 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF:
			inst = "XRA"
			op1 = operandOrder[opcode-0xA8]

		case 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7:
			inst = "ORA"
			op1 = operandOrder[opcode-0xB0]

		case 0xB8, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBE, 0xBF:
			inst = "CMP"
			op1 = operandOrder[opcode-0xB8]

		case 0xC0:
			inst = "RNZ"
		case 0xD0:
			inst = "RNC"
		case 0xE0:
			inst = "RPO"
		case 0xF0:
			inst = "RP"

		case 0xC1, 0xD1, 0xE1, 0xF1:
			inst = "POP"
			op1 = operandOrder[(opcode-0xC1)/0x8]
			if opcode == 0xF1 {
				op1 = "AF"
			}

		case 0xE3:
			inst = "XTHL"

		case 0xF3:
			inst = "DI"

		case 0xC5, 0xD5, 0xE5, 0xF5:
			inst = "PUSH"
			op1 = operandOrder[(opcode-0xC5)/0x8]
			if opcode == 0xF5 {
				op1 = "AF"
			}

		case 0xC7, 0xD7, 0xE7, 0xF7, 0xCF, 0xDF, 0xEF, 0xFF:
			inst = "RST"
			op1 = strconv.FormatUint(uint64(opcode-0xC7), 16)

		case 0xC8:
			inst = "RZ"
		case 0xD8:
			inst = "RC"
		case 0xE8:
			inst = "RPE"
		case 0xF8:
			inst = "RM"

		case 0xC9, 0xD9:
			inst = "RET"

		case 0xE9:
			inst = "PCHL"

		case 0xF9:
			inst = "SPHL"

		case 0xEB:
			inst = "XCHG"

		case 0xFB:
			inst = "EI"

		}
	case 0xD3, 0x06, 0x16, 0x26, 0x36, 0xC6, 0xD6, 0xE6, 0xF6, 0xDB, 0x0E, 0x1E, 0x2E, 0x3E, 0xCE, 0xDE, 0xEE, 0xFE:
		op1 = fmt.Sprintf("%02X", c.ReadMem(c.pc+1))
		length = 2

		switch opcode {
		case 0xD3:
			inst = "OUT"
		case 0x06, 0x16, 0x26, 0x36, 0x0E, 0x1E, 0x2E, 0x3E:
			inst = "MVI"
			op2 = op1
			op1 = operandOrder[(opcode-0x06)/8]
		case 0xC6:
			inst = "ADI"
		case 0xD6:
			inst = "SUI"
		case 0xE6:
			inst = "ANI"
		case 0xF6:
			inst = "ORI"
		case 0xDB:
			inst = "IN"
		case 0xCE:
			inst = "ACI"
		case 0xDE:
			inst = "SBI"
		case 0xEE:
			inst = "XRI"
		case 0xFE:
			inst = "CPI"

		}

	case 0x01, 0x11, 0x21, 0x31, 0x22, 0x32, 0xC2, 0xD2, 0xE2, 0xF2, 0xC3, 0xC4, 0xD4, 0xE4, 0xF4, 0x2A, 0x3A, 0xCA, 0xDA, 0xEA, 0xFA, 0xCB, 0xCC, 0xDC, 0xEC, 0xFC, 0xCD, 0xDD, 0xED, 0xFD, 0xDF:
		length = 3
		op1 = fmt.Sprintf("%02X%02X", c.ReadMem(c.pc+2), c.ReadMem(c.pc+1))

		switch opcode {
		case 0x01, 0x11, 0x21, 0x31:
			inst = "LXI"
			op2 = op1

			switch hi {
			case 0:
				op1 = "B"
			case 1:
				op1 = "D"
			case 2:
				op1 = "H"
			case 3:
				op1 = "SP"
			}

		case 0x22:
			inst = "SHLD"

		case 0x32:
			inst = "STA"

		case 0xC2:
			inst = "JNZ"

		case 0xD2:
			inst = "JNC"

		case 0xE2:
			inst = "JPO"

		case 0xF2:
			inst = "JP"

		case 0xC3:
			inst = "JMP"

		case 0xC4:
			inst = "CNZ"

		case 0xD4:
			inst = "CNC"

		case 0xE4:
			inst = "CPO"

		case 0xF4:
			inst = "CP"

		case 0x2A:
			inst = "LHDL"

		case 0x3A:
			inst = "LDA"

		case 0xCA:
			inst = "JZ"

		case 0xDA:
			inst = "JC"

		case 0xEA:
			inst = "JPE"

		case 0xFA:
			inst = "JM"

		case 0xCB:
			inst = "JMP"

		case 0xCC:
			inst = "CZ"

		case 0xDC:
			inst = "CC"

		case 0xEC:
			inst = "CPE"

		case 0xFC:
			inst = "CM"

		case 0xCD, 0xDD, 0xED, 0xFD:
			inst = "CALL"
		}
	}

	fullInst := fmt.Sprintf("%-4s", inst)

	if op1 != "" {
		fullInst += " " + op1
		if op2 != "" {
			fullInst += ", " + op2
		}
	}

	if c.debug {
		fmt.Printf("%02X %02X %02X %-13s%s\n", c.ReadMem(c.pc), c.ReadMem(c.pc+1), c.ReadMem(c.pc+2), fullInst, c)
	}

	return inst, op1, op2, length
}

func (c *CPU) Step() uint64 {
	// Fetch instruction
	inst, op1, op2, instLength := c.DecodeInst()
	prevSC, prevPC, states := c.sc, c.pc, uint64(0)

	regMap := map[string]*uint8{
		"B": &c.b,
		"C": &c.c,
		"D": &c.d,
		"E": &c.e,
		"H": &c.h,
		"L": &c.l,
		"A": &c.a,
		"F": &c.f,
	}

	switch inst {
	case "NOP":
		states += 4

	case "LXI":
		value := uint16(must(strconv.ParseUint(op2, 16, 16)))
		switch op1 {
		case "B":
			c.b = uint8(value >> 8)
			c.c = uint8(value)
		case "D":
			c.d = uint8(value >> 8)
			c.e = uint8(value)
		case "H":
			c.h = uint8(value >> 8)
			c.l = uint8(value)
		case "SP":
			c.sp = value
		}
		states += 10

	case "STA":
		addr := uint16(must(strconv.ParseUint(op1, 16, 16)))
		c.WriteMem(addr, c.a)
		states += 13

	case "LDA":
		addr := uint16(must(strconv.ParseUint(op1, 16, 16)))
		c.a = c.ReadMem(addr)
		states += 13

	case "LDAX":
		var value uint8

		switch op1 {
		case "BC":
			value = c.ReadMem(uint16(c.b)<<8 | uint16(c.c))
		case "DE":
			value = c.ReadMem(uint16(c.d)<<8 | uint16(c.e))
		}

		c.a = value
		states += 7

	case "JNZ":
		addr := uint16(must(strconv.ParseUint(op1, 16, 16)))
		states += 10

		if c.getZ() != 0 {
			c.pc = addr

			return states
		}

	case "JMP":
		addr := uint16(must(strconv.ParseUint(op1, 16, 16)))
		states += 10
		c.pc = addr

		return states

	case "CALL":
		addr := uint16(must(strconv.ParseUint(op1, 16, 16)))
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
		value := uint8(must(strconv.ParseUint(op2, 16, 8)))

		switch op1 {
		case "B":
			c.b = value
		case "C":
			c.c = value
		case "D":
			c.d = value
		case "E":
			c.e = value
		case "H":
			c.h = value
		case "L":
			c.l = value
		case "M":
			c.WriteMem(c.getHL(), value)
			states += 3
		case "A":
			c.a = value
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
		switch op1 {
		case "BC":
			c.c++
			if c.c == 0 {
				c.b++
			}
		case "DE":
			c.e++
			if c.e == 0 {
				c.d++
			}
		case "HL":
			c.l++
			if c.l == 0 {
				c.h++
			}
		case "SP":
			c.sp++
		}

		states += 5

	case "DCR":
		var value uint8
		switch op1 {
		case "B":
			c.b--
			value = c.b
		case "D":
			c.d--
			value = c.d
		case "H":
			c.h--
			value = c.h
		case "M":
			value = c.ReadMem(c.getHL()) - 1
			c.WriteMem(c.getHL(), value)

			states += 5
		}

		c.setS(value>>7&1 == 1)
		c.setZ(value == 0)
		c.setAC(value&0x0F != 0x0F)
		c.setP(value&1 == 0)

		states += 5

	case "DAD":
		var res uint32

		switch op1 {
		case "BC":
			res = uint32(c.getHL()) + uint32(c.getBC())
		case "DE":
			res = uint32(c.getHL()) + uint32(c.getDE())
		case "HL":
			res = uint32(c.getHL()) + uint32(c.getHL())
		case "SP":
			res = uint32(c.getHL()) + uint32(c.sp)
		}

		c.setCY(res > math.MaxUint16)
		c.setHL(uint16(res))

		states += 10

	case "ANI":
		value := uint8(must(strconv.ParseUint(op1, 16, 8)))

		res := c.a & value

		c.setS(res>>7&1 == 1)
		c.setZ(res == 0)
		c.setAC(false)
		c.setP(res&1 == 0)
		c.setCY(false)

		c.a = res

		states += 7

	case "CPI":
		value := uint8(must(strconv.ParseUint(op1, 16, 8)))
		res := c.a - value

		c.setS(res>>7&1 == 1)
		c.setZ(res == 0)
		c.setAC(c.a&0x0F < value&0x0F)
		c.setP(res&1 == 0)
		c.setCY(c.a < value)

		states += 7

	case "PUSH":
		switch op1 {
		case "B":
			c.push(c.getBC())
		case "D":
			c.push(c.getDE())
		case "H":
			c.push(c.getHL())
		case "A":
			c.push(c.getAF())
		default:
			panic("undefined operand for push: " + op1)
		}

		states += 11

	case "POP":
		switch op1 {
		case "B":
			c.setBC(c.pop())
		case "D":
			c.setDE(c.pop())
		case "H":
			c.setHL(c.pop())
		case "A":
			c.setAF(c.pop())
		}

		states += 10

	case "XCHG":
		c.h, c.l, c.d, c.e = c.d, c.e, c.h, c.l
		states += 4

	case "OUT":
		value := uint8(must(strconv.ParseUint(op1, 16, 8)))
		c.WriteMem(uint16(value)<<8|uint16(value), c.a)

		states += 10

	case "RST":
		addr := uint16(must(strconv.ParseUint(op1, 16, 16)))
		states += 11
		c.pc = addr

		return states

	case "RRC":
		c.setCY(c.a&0x1 == 1)
		c.a = c.a >> 1
		states += 4

	case "ADI":
		value := uint8(must(strconv.ParseUint(op1, 16, 8)))

		res := c.a + value

		c.setS(res>>7&1 == 1)
		c.setZ(res == 0)
		c.setAC(c.a&0x0F+value&0x0F > 0x0F)
		c.setP(res&1 == 0)
		c.setCY(uint16(c.a)+uint16(value) > 0xFF)

		c.a = res
		states += 7

	case "XRA":
		value := *regMap[op1]

		res := c.a ^ value

		c.setS(res>>7&1 == 1)
		c.setZ(res == 0)
		c.setAC(false)
		c.setP(res&1 == 0)
		c.setCY(false)

		c.a = res
		states += 4

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

func (c *CPU) setAF(af uint16) {
	c.a = uint8(af >> 8)
	c.f = uint8(af)
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

func must[T any](res T, err error) T {
	if err != nil {
		panic(err)
	}

	return res
}
