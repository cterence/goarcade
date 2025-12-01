package cpu

import (
	"fmt"
	"strconv"
)

func (c *CPU) DecodeInst() (string, string, string, uint16, uint64) {
	operandOrder := []string{"B", "C", "D", "E", "H", "L", "M", "A"}
	doubleOperandOrderAF := []string{"BC", "DE", "HL", "AF"}
	doubleOperandOrderSP := []string{"BC", "DE", "HL", "SP"}
	opcode := c.ReadMem(c.pc)

	// return values
	var (
		inst   string
		op1    string
		op2    string
		length uint16
		states uint64
	)

	switch opcode {
	default:
		length = 1

		switch opcode {
		case 0x00, 0x10, 0x20, 0x30, 0x08, 0x18, 0x28, 0x38:
			inst = "NOP"
			states = 4

		case 0x02, 0x12:
			inst = "STAX"
			op1 = operandOrder[(opcode-0x2)/0x8]
			states = 7

		case 0x03, 0x13, 0x23, 0x33:
			inst = "INX"

			op1 = doubleOperandOrderSP[(opcode-0x3)/0x10]
			if opcode == 0x33 {
				op1 = "SP"
			}

			states = 5

		case 0x04, 0x14, 0x24, 0x34, 0x0C, 0x1C, 0x2C, 0x3C:
			inst = "INR"
			op1 = operandOrder[(opcode-0x4)/0x8]

			states = 5
			if op1 == "M" {
				states = 10
			}

		case 0x05, 0x15, 0x25, 0x35, 0x0D, 0x1D, 0x2D, 0x3D:
			inst = "DCR"
			op1 = operandOrder[(opcode-0x4)/0x8]

			states = 5
			if op1 == "M" {
				states = 10
			}

		case 0x07:
			inst = "RLC"
			states = 4

		case 0x17:
			inst = "RAL"
			states = 4

		case 0x27:
			inst = "DAA"
			states = 4

		case 0x37:
			inst = "STC"
			states = 4

		case 0x09, 0x19, 0x29, 0x39:
			inst = "DAD"

			op1 = doubleOperandOrderSP[(opcode-0x9)/0x10]
			if opcode == 0x39 {
				op1 = "SP"
			}

			states = 10

		case 0x0A, 0x1A:
			inst = "LDAX"
			op1 = doubleOperandOrderAF[(opcode-0xA)/0x8]
			states = 7

		case 0x0B, 0x1B, 0x2B, 0x3B:
			inst = "DCX"

			op1 = operandOrder[(opcode-0xB)/0x8]
			if opcode == 0x3B {
				op1 = "SP"
			}

			states = 4

		case 0x0F:
			inst = "RRC"
			states = 4

		case 0x1F:
			inst = "RAR"
			states = 4

		case 0x2F:
			inst = "CMA"
			states = 4

		case 0x3F:
			inst = "CMC"
			states = 4

		case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F, 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F:
			inst = "MOV"
			op1 = operandOrder[(opcode-0x40)/8]
			op2 = operandOrder[(opcode>>4)%8]

			states = 5
			if op1 == "M" {
				states = 7
			}

		case 0x76:
			inst = "HALT"
			states = 7

		case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87:
			inst = "ADD"
			op1 = operandOrder[opcode-0x80]

			states = 4
			if op1 == "M" {
				states = 7
			}

		case 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F:
			inst = "ADC"
			op1 = operandOrder[opcode-0x88]

			states = 4
			if op1 == "M" {
				states = 7
			}

		case 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97:
			inst = "SUB"
			op1 = operandOrder[opcode-0x90]

			states = 4
			if op1 == "M" {
				states = 7
			}

		case 0x98, 0x99, 0x9A, 0x9B, 0x9C, 0x9D, 0x9E, 0x9F:
			inst = "SBB"
			op1 = operandOrder[opcode-0x98]

			states = 4
			if op1 == "M" {
				states = 7
			}

		case 0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7:
			inst = "ANA"
			op1 = operandOrder[opcode-0xA0]

			states = 4
			if op1 == "M" {
				states = 7
			}

		case 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF:
			inst = "XRA"
			op1 = operandOrder[opcode-0xA8]

			states = 4
			if op1 == "M" {
				states = 7
			}

		case 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7:
			inst = "ORA"
			op1 = operandOrder[opcode-0xB0]

			states = 4
			if op1 == "M" {
				states = 7
			}

		case 0xB8, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBE, 0xBF:
			inst = "CMP"
			op1 = operandOrder[opcode-0xB8]

			states = 4
			if op1 == "M" {
				states = 7
			}

		case 0xC0:
			inst = "RNZ"

			states = 5
			if c.getZF() == 0 {
				states = 11
			}

		case 0xD0:
			inst = "RNC"

			states = 5
			if c.getCYF() == 0 {
				states = 11
			}

		case 0xE0:
			inst = "RPO"

			states = 5
			if c.getPF() == 0 {
				states = 11
			}

		case 0xF0:
			inst = "RP"

			states = 5
			if c.getSF() == 0 {
				states = 11
			}

		case 0xC1, 0xD1, 0xE1, 0xF1:
			inst = "POP"

			op1 = doubleOperandOrderAF[(opcode-0xC1)/0x10]
			if opcode == 0xF1 {
				op1 = "AF"
			}

			states = 10

		case 0xE3:
			inst = "XTHL"
			states = 18

		case 0xF3:
			inst = "DI"
			states = 14

		case 0xC5, 0xD5, 0xE5, 0xF5:
			inst = "PUSH"

			op1 = doubleOperandOrderAF[(opcode-0xC5)/0x10]
			if opcode == 0xF5 {
				op1 = "AF"
			}

			states = 11

		case 0xC7, 0xD7, 0xE7, 0xF7, 0xCF, 0xDF, 0xEF, 0xFF:
			inst = "RST"
			op1 = strconv.FormatUint(uint64(opcode-0xC7), 16)
			states = 11

		case 0xC8:
			inst = "RZ"

			states = 5
			if c.getSF() == 1 {
				states = 11
			}

		case 0xD8:
			inst = "RC"

			states = 5
			if c.getCYF() == 1 {
				states = 11
			}

		case 0xE8:
			inst = "RPE"

			states = 5
			if c.getPF() == 1 {
				states = 11
			}

		case 0xF8:
			inst = "RM"

			states = 5
			if c.getSF() == 1 {
				states = 11
			}

		case 0xC9, 0xD9:
			inst = "RET"
			states = 10

		case 0xE9:
			inst = "PCHL"
			states = 5

		case 0xF9:
			inst = "SPHL"
			states = 5

		case 0xEB:
			inst = "XCHG"
			states = 5

		case 0xFB:
			inst = "EI"
			states = 4

		default:
			panic("undecoded opcode: " + strconv.FormatUint(uint64(opcode), 16))
		}
	case 0xD3, 0x06, 0x16, 0x26, 0x36, 0xC6, 0xD6, 0xE6, 0xF6, 0xDB, 0x0E, 0x1E, 0x2E, 0x3E, 0xCE, 0xDE, 0xEE, 0xFE:
		op1 = fmt.Sprintf("%02X", c.ReadMem(c.pc+1))
		length = 2

		switch opcode {
		case 0xD3:
			inst = "OUT"
			states = 10

		case 0x06, 0x16, 0x26, 0x36, 0x0E, 0x1E, 0x2E, 0x3E:
			inst = "MVI"
			op2 = op1
			op1 = operandOrder[(opcode-0x06)/8]

			states = 7
			if op1 == "M" {
				states = 10
			}

		case 0xC6:
			inst = "ADI"
			states = 7

		case 0xD6:
			inst = "SUI"
			states = 7

		case 0xE6:
			inst = "ANI"
			states = 7

		case 0xF6:
			inst = "ORI"
			states = 7

		case 0xDB:
			inst = "IN"
			states = 10

		case 0xCE:
			inst = "ACI"
			states = 7

		case 0xDE:
			inst = "SBI"
			states = 7

		case 0xEE:
			inst = "XRI"
			states = 7

		case 0xFE:
			inst = "CPI"
			states = 7

		default:
			panic("undecoded opcode: " + strconv.FormatUint(uint64(opcode), 16))
		}

	case 0x01, 0x11, 0x21, 0x31, 0x22, 0x32, 0xC2, 0xD2, 0xE2, 0xF2, 0xC3, 0xC4, 0xD4, 0xE4, 0xF4, 0x2A, 0x3A, 0xCA, 0xDA, 0xEA, 0xFA, 0xCB, 0xCC, 0xDC, 0xEC, 0xFC, 0xCD, 0xDD, 0xED, 0xFD, 0xDF:
		length = 3
		op1 = fmt.Sprintf("%02X%02X", c.ReadMem(c.pc+2), c.ReadMem(c.pc+1))

		switch opcode {
		case 0x01, 0x11, 0x21, 0x31:
			inst = "LXI"
			op2 = op1
			op1 = doubleOperandOrderSP[opcode>>4]
			states = 10

		case 0x22:
			inst = "SHLD"
			states = 16

		case 0x32:
			inst = "STA"
			states = 13

		case 0xC2:
			inst = "JNZ"
			states = 10

		case 0xD2:
			inst = "JNC"
			states = 10

		case 0xE2:
			inst = "JPO"
			states = 10

		case 0xF2:
			inst = "JP"
			states = 10

		case 0xC3, 0xCB:
			inst = "JMP"
			states = 10

		case 0xC4:
			inst = "CNZ"

			states = 11
			if c.getZF() == 0 {
				states = 17
			}

		case 0xD4:
			inst = "CNC"

			states = 11
			if c.getCYF() == 0 {
				states = 17
			}

		case 0xE4:
			inst = "CPO"

			states = 11
			if c.getPF() == 0 {
				states = 17
			}

		case 0xF4:
			inst = "CP"

			if c.getSF() == 0 {
				states = 17
			}

		case 0x2A:
			inst = "LHDL"
			states = 11

		case 0x3A:
			inst = "LDA"
			states = 13

		case 0xCA:
			inst = "JZ"
			states = 10

		case 0xDA:
			inst = "JC"
			states = 10

		case 0xEA:
			inst = "JPE"
			states = 10

		case 0xFA:
			inst = "JM"
			states = 10

		case 0xCC:
			inst = "CZ"

			states = 11
			if c.getZF() == 1 {
				states = 17
			}

		case 0xDC:
			inst = "CC"

			states = 11
			if c.getCYF() == 1 {
				states = 17
			}

		case 0xEC:
			inst = "CPE"

			states = 11
			if c.getPF() == 1 {
				states = 17
			}

		case 0xFC:
			inst = "CM"

			states = 11
			if c.getSF() == 1 {
				states = 17
			}

		case 0xCD, 0xDD, 0xED, 0xFD:
			inst = "CALL"
			states = 17

		default:
			panic("undecoded opcode: " + strconv.FormatUint(uint64(opcode), 16))
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

	return inst, op1, op2, length, states
}
