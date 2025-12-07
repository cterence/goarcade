package cpu

import (
	"fmt"
	"math"
	"math/bits"
	"strconv"

	"github.com/cterence/space-invaders/internal/arcade/lib"
)

func nop(*CPU, string) {}

func ldax(c *CPU, op string) {
	c.A = c.Bus.Read(c.getDoubleOp(op))
}

func stax(c *CPU, op string) {
	c.Bus.Write(c.getDoubleOp(op), c.A)
}

func inx(c *CPU, op string) {
	c.setDoubleOp(op, c.getDoubleOp(op)+1)
}

func dcx(c *CPU, op string) {
	c.setDoubleOp(op, c.getDoubleOp(op)-1)
}

func inr(c *CPU, op string) {
	value := c.getOp(op)
	res := value + 1
	c.setOp(op, res)
	c.setFlags(res&0x80 == 0x80, res == 0, value&0x0F+1 > 0x0F, bits.OnesCount8(res)%2 == 0, c.getCYF() == 1)
}

func dcr(c *CPU, op string) {
	value := c.getOp(op)
	res := value - 1
	c.setOp(op, res)
	c.setFlags(res&0x80 == 0x80, res == 0, value&0x0F >= 1, bits.OnesCount8(res)%2 == 0, c.getCYF() == 1)
}

func dad(c *CPU, op string) {
	res := uint32(uint16(c.H)<<8|uint16(c.L)) + uint32(c.getDoubleOp(op))

	c.setCYF(res > math.MaxUint16)
	c.H = uint8(res >> 8)
	c.L = uint8(res)
}

func lxi(c *CPU, op string) {
	imm1, imm2 := c.Bus.Read(c.PC+1), c.Bus.Read(c.PC+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	c.setDoubleOp(op, imm16)
}

func (c *CPU) pop() uint16 {
	value := uint16(c.Bus.Read(c.SP+1))<<8 | uint16(c.Bus.Read(c.SP))
	c.SP += 2

	return value
}

func popOp(c *CPU, op string) {
	value := c.pop()
	c.setDoubleOp(op, value)
}

func (c *CPU) push(addr uint16) {
	c.SP -= 2
	c.Bus.Write(c.SP, uint8(addr))
	c.Bus.Write(c.SP+1, uint8(addr>>8))
}

func pushOp(c *CPU, op string) {
	c.SP -= 2
	value := c.getDoubleOp(op)
	c.Bus.Write(c.SP, uint8(value))
	c.Bus.Write(c.SP+1, uint8(value>>8))
}

func (c *CPU) jumpCond(cond bool) {
	imm1, imm2 := c.Bus.Read(c.PC+1), c.Bus.Read(c.PC+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)

	if cond {
		c.PC = imm16
	}
}

func (c *CPU) call() {
	imm1, imm2 := c.Bus.Read(c.PC+1), c.Bus.Read(c.PC+2)
	addr := uint16(imm2)<<8 | uint16(imm1)

	c.push(c.PC + 3)
	c.PC = addr
}

func (c *CPU) callCond(cond bool) {
	if cond {
		c.call()
		c.Cyc += 6
	}
}

func (c *CPU) ret() {
	c.PC = c.pop()
}

func (c *CPU) retCond(cond bool) {
	if cond {
		c.ret()
		c.Cyc += 6
	}
}

func mvi(c *CPU, op string) {
	imm1 := c.Bus.Read(c.PC + 1)
	c.setOp(op, imm1)
}

func rlc(c *CPU, _ string) {
	sb := c.A & 0x80 >> 7
	c.setCYF(sb == 1)
	c.A = c.A<<1 | sb
}
func rrc(c *CPU, _ string) {
	sb := c.A & 0x1
	c.setCYF(sb == 1)
	c.A = c.A>>1 | sb<<7
}

func ral(c *CPU, _ string) {
	sb := c.A & 0x80 >> 7
	c.A = c.A<<1 | c.getCYF()
	c.setCYF(sb == 1)
}

func rar(c *CPU, _ string) {
	sb := c.A & 0x1
	c.A = c.A>>1 | c.getCYF()<<7
	c.setCYF(sb == 1)
}

func shld(c *CPU, _ string) {
	imm1, imm2 := c.Bus.Read(c.PC+1), c.Bus.Read(c.PC+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	c.Bus.Write(imm16, c.L)
	c.Bus.Write(imm16+1, c.H)
}

func daa(c *CPU, _ string) {
	cy := c.getCYF() == 1
	value := uint8(0)

	if c.A&0x0F > 0x09 || c.getACF() == 1 {
		value += 0x06
	}

	if (c.A+value)&0xF0 > 0x90 || cy || c.A > 0x99 {
		value += 0x60
		cy = true
	}

	res := c.A + value

	c.setFlags(res&0x80 == 0x80, res == 0, c.A&0x0F+value&0x0F > 0x0F, bits.OnesCount8(res)%2 == 0, cy)

	c.A = res
}

func lhld(c *CPU, _ string) {
	imm1, imm2 := c.Bus.Read(c.PC+1), c.Bus.Read(c.PC+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	res := uint16(c.Bus.Read(imm16+1))<<8 | uint16(c.Bus.Read(imm16))
	c.H = uint8(res >> 8)
	c.L = uint8(res)
}

func sta(c *CPU, _ string) {
	imm1, imm2 := c.Bus.Read(c.PC+1), c.Bus.Read(c.PC+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	c.Bus.Write(imm16, c.A)
}

func lda(c *CPU, _ string) {
	imm1, imm2 := c.Bus.Read(c.PC+1), c.Bus.Read(c.PC+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	c.A = c.Bus.Read(imm16)
}

func add(c *CPU, op string) {
	value := c.getOp(op)
	res := c.A + value
	c.setFlags(res&0x80 == 0x80, res == 0, c.A&0x0F+(value)&0x0F > 0x0F, bits.OnesCount8(res)%2 == 0, uint16(c.A)+uint16(value) > 0xFF)
	c.A = res
}

func adc(c *CPU, op string) {
	value := c.getOp(op)
	carry := c.getCYF()
	res := c.A + value + carry
	c.setFlags(res&0x80 == 0x80, res == 0, c.A&0x0F+(value)&0x0F+carry > 0x0F, bits.OnesCount8(res)%2 == 0, uint16(c.A)+uint16(value)+uint16(carry) > 0xFF)
	c.A = res
}

func sub(c *CPU, op string) {
	value := c.getOp(op)
	res := c.A - value
	c.setFlags(res&0x80 == 0x80, res == 0, (c.A&0x0F) >= (value&0x0F), bits.OnesCount8(res)%2 == 0, uint16(c.A) < uint16(value))
	c.A = res
}

func sbb(c *CPU, op string) {
	value := c.getOp(op)
	carry := c.getCYF()
	res := c.A - value - carry
	c.setFlags(res&0x80 == 0x80, res == 0, (c.A&0x0F) >= (value&0x0F)+carry, bits.OnesCount8(res)%2 == 0, uint16(c.A) < uint16(value)+uint16(carry))
	c.A = res
}

func ana(c *CPU, op string) {
	value := c.getOp(op)
	res := c.A & value
	c.setFlags(res&0x80 == 0x80, res == 0, (c.A|value)&0x08 != 0, bits.OnesCount8(res)%2 == 0, false)
	c.A = res
}

func xra(c *CPU, op string) {
	value := c.getOp(op)
	res := c.A ^ value
	c.setFlags(res&0x80 == 0x80, res == 0, false, bits.OnesCount8(res)%2 == 0, false)
	c.A = res
}

func ora(c *CPU, op string) {
	value := c.getOp(op)
	res := c.A | value
	c.setFlags(res&0x80 == 0x80, res == 0, false, bits.OnesCount8(res)%2 == 0, false)
	c.A = res
}

func cmp(c *CPU, op string) {
	value := c.getOp(op)
	res := c.A - value
	c.setFlags(res&0x80 == 0x80, res == 0, c.A&0x0F >= value&0x0F, bits.OnesCount8(res)%2 == 0, c.A < value)
}

func hlt(c *CPU, _ string) {
	// TODO: implement
}

func xthl(c *CPU, _ string) {
	lo := c.Bus.Read(c.SP)
	hi := c.Bus.Read(c.SP + 1)
	c.Bus.Write(c.SP, c.L)
	c.Bus.Write(c.SP+1, c.H)
	c.L = lo
	c.H = hi
}

func cpi(c *CPU, _ string) {
	imm1 := c.Bus.Read(c.PC + 1)
	value := imm1
	res := c.A - value

	c.setFlags(res&0x80 == 0x80, res == 0, c.A&0x0F >= value&0x0F, bits.OnesCount8(res)%2 == 0, c.A < value)
}

func rst(c *CPU, op string) {
	c.push(c.PC + 1)
	c.PC = uint16(lib.Must(strconv.ParseUint(op, 16, 16)) * 8)
}

func adi(c *CPU, _ string) {
	imm1 := c.Bus.Read(c.PC + 1)
	value := imm1
	res := c.A + value
	c.setFlags(res&0x80 == 0x80, res == 0, c.A&0x0F+(value)&0x0F > 0x0F, bits.OnesCount8(res)%2 == 0, uint16(c.A)+uint16(value) > 0xFF)
	c.A = res
}

func aci(c *CPU, _ string) {
	imm1 := c.Bus.Read(c.PC + 1)
	value := imm1
	carry := c.getCYF()
	res := c.A + value + carry
	c.setFlags(res&0x80 == 0x80, res == 0, c.A&0x0F+(value)&0x0F+carry > 0x0F, bits.OnesCount8(res)%2 == 0, uint16(c.A)+uint16(value)+uint16(carry) > 0xFF)
	c.A = res
}

func sui(c *CPU, _ string) {
	imm1 := c.Bus.Read(c.PC + 1)
	value := imm1
	res := c.A - value
	c.setFlags(res&0x80 == 0x80, res == 0, (c.A&0x0F) >= (value&0x0F), bits.OnesCount8(res)%2 == 0, uint16(c.A) < uint16(value))
	c.A = res
}

func sbi(c *CPU, _ string) {
	imm1 := c.Bus.Read(c.PC + 1)
	value := imm1
	carry := c.getCYF()
	res := c.A - value - carry
	c.setFlags(res&0x80 == 0x80, res == 0, (c.A&0x0F) >= (value&0x0F)+carry, bits.OnesCount8(res)%2 == 0, uint16(c.A) < uint16(value)+uint16(carry))
	c.A = res
}

func ani(c *CPU, _ string) {
	imm1 := c.Bus.Read(c.PC + 1)
	value := imm1
	res := c.A & value

	c.setFlags(res&0x80 == 0x80, res == 0, (c.A|value)&0x08 != 0, bits.OnesCount8(res)%2 == 0, false)
	c.A = res
}

func xri(c *CPU, _ string) {
	imm1 := c.Bus.Read(c.PC + 1)
	value := imm1
	res := c.A ^ value
	c.setFlags(res&0x80 == 0x80, res == 0, false, bits.OnesCount8(res)%2 == 0, false)
	c.A = res
}

func ori(c *CPU, _ string) {
	imm1 := c.Bus.Read(c.PC + 1)
	value := imm1
	res := c.A | value
	c.setFlags(res&0x80 == 0x80, res == 0, false, bits.OnesCount8(res)%2 == 0, false)
	c.A = res
}

func portIn(c *CPU, _ string) {
	portNumber := c.Bus.Read(c.PC + 1)

	switch portNumber {
	case 3:
		c.A = uint8(c.SR >> (8 - c.SO))
	default:
		c.A = c.IOPorts[portNumber]
	}
}

func portOut(c *CPU, _ string) {
	portNumber := c.Bus.Read(c.PC + 1)

	switch portNumber {
	case 0:
		c.Running = false
	case 1:
		switch c.C {
		case 2:
			fmt.Print(string(c.E))
		case 9:
			addr := uint16(c.D)<<8 | uint16(c.E)
			for offset := uint16(0); ; offset++ {
				b := c.Bus.Read(addr + offset)
				if b == '$' {
					break
				}

				fmt.Print(string(b))
			}

		default:
			fmt.Printf("unimplemented out operation for port 1: %02x\n", c.C)
		}
	case 2:
		c.SO = c.A & 0x7
	case 3:
		rising := c.A & ^c.IOPorts[portNumber]
		falling := c.IOPorts[portNumber] & ^c.A

		switch {
		case rising&1 == 1:
			c.APU.StartSoundLoop(0)
		case falling&1 == 1:
			c.APU.StopSoundLoop(0)
		case rising>>1&1 == 1:
			c.APU.PlaySound(1)
		case rising>>2&1 == 1:
			c.APU.PlaySound(2)
		case rising>>3&1 == 1:
			c.APU.PlaySound(3)
		}

		c.IOPorts[portNumber] = c.A
	case 4:
		c.SR = uint16(c.A)<<8 | c.SR>>8
	case 5:
		rising := c.A & ^c.IOPorts[portNumber]

		switch {
		case rising&1 == 1:
			c.APU.PlaySound(4)
		case rising>>1&1 == 1:
			c.APU.PlaySound(5)
		case rising>>2&1 == 1:
			c.APU.PlaySound(6)
		case rising>>3&1 == 1:
			c.APU.PlaySound(7)
		case rising>>4&1 == 1:
			c.APU.PlaySound(8)
		}

		c.IOPorts[portNumber] = c.A
	case 6: // NOP for watchdog
	default:
		fmt.Printf("unimplemented out port: %02x\n", portNumber)
	}
}
