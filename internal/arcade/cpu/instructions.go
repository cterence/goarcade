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
	c.a = c.ReadMem(c.getDoubleOp(op))
}

func stax(c *CPU, op string) {
	c.WriteMem(c.getDoubleOp(op), c.a)
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
	res := uint32(uint16(c.h)<<8|uint16(c.l)) + uint32(c.getDoubleOp(op))

	c.setCYF(res > math.MaxUint16)
	c.h = uint8(res >> 8)
	c.l = uint8(res)
}

func lxi(c *CPU, op string) {
	imm1, imm2 := c.ReadMem(c.pc+1), c.ReadMem(c.pc+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	c.setDoubleOp(op, imm16)
}

func (c *CPU) pop() uint16 {
	value := uint16(c.ReadMem(c.sp+1))<<8 | uint16(c.ReadMem(c.sp))
	c.sp += 2

	return value
}

func popOp(c *CPU, op string) {
	value := c.pop()
	c.setDoubleOp(op, value)
}

func (c *CPU) push(addr uint16) {
	c.sp -= 2
	c.WriteMem(c.sp, uint8(addr))
	c.WriteMem(c.sp+1, uint8(addr>>8))
}

func pushOp(c *CPU, op string) {
	c.sp -= 2
	value := c.getDoubleOp(op)
	c.WriteMem(c.sp, uint8(value))
	c.WriteMem(c.sp+1, uint8(value>>8))
}

func (c *CPU) jumpCond(cond bool) {
	imm1, imm2 := c.ReadMem(c.pc+1), c.ReadMem(c.pc+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)

	if cond {
		c.pc = imm16
	}
}

func (c *CPU) call() {
	imm1, imm2 := c.ReadMem(c.pc+1), c.ReadMem(c.pc+2)
	addr := uint16(imm2)<<8 | uint16(imm1)

	c.push(c.pc + 3)
	c.pc = addr
}

func (c *CPU) callCond(cond bool) {
	if cond {
		c.call()
		c.cyc += 6
	}
}

func (c *CPU) ret() {
	c.pc = c.pop()
}

func (c *CPU) retCond(cond bool) {
	if cond {
		c.ret()
		c.cyc += 6
	}
}

func mvi(c *CPU, op string) {
	imm1 := c.ReadMem(c.pc + 1)
	c.setOp(op, imm1)
}

func rlc(c *CPU, _ string) {
	sb := c.a & 0x80 >> 7
	c.setCYF(sb == 1)
	c.a = c.a<<1 | sb
}
func rrc(c *CPU, _ string) {
	sb := c.a & 0x1
	c.setCYF(sb == 1)
	c.a = c.a>>1 | sb<<7
}

func ral(c *CPU, _ string) {
	sb := c.a & 0x80 >> 7
	c.a = c.a<<1 | c.getCYF()
	c.setCYF(sb == 1)
}

func rar(c *CPU, _ string) {
	sb := c.a & 0x1
	c.a = c.a>>1 | c.getCYF()<<7
	c.setCYF(sb == 1)
}

func shld(c *CPU, _ string) {
	imm1, imm2 := c.ReadMem(c.pc+1), c.ReadMem(c.pc+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	c.WriteMem(imm16, c.l)
	c.WriteMem(imm16+1, c.h)
}

func daa(c *CPU, _ string) {
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
}

func lhld(c *CPU, _ string) {
	imm1, imm2 := c.ReadMem(c.pc+1), c.ReadMem(c.pc+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	res := uint16(c.ReadMem(imm16+1))<<8 | uint16(c.ReadMem(imm16))
	c.h = uint8(res >> 8)
	c.l = uint8(res)
}

func sta(c *CPU, _ string) {
	imm1, imm2 := c.ReadMem(c.pc+1), c.ReadMem(c.pc+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	c.WriteMem(imm16, c.a)
}

func lda(c *CPU, _ string) {
	imm1, imm2 := c.ReadMem(c.pc+1), c.ReadMem(c.pc+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	c.a = c.ReadMem(imm16)
}

func add(c *CPU, op string) {
	value := c.getOp(op)
	res := c.a + value
	c.setFlags(res&0x80 == 0x80, res == 0, c.a&0x0F+(value)&0x0F > 0x0F, bits.OnesCount8(res)%2 == 0, uint16(c.a)+uint16(value) > 0xFF)
	c.a = res
}

func adc(c *CPU, op string) {
	value := c.getOp(op)
	carry := c.getCYF()
	res := c.a + value + carry
	c.setFlags(res&0x80 == 0x80, res == 0, c.a&0x0F+(value)&0x0F+carry > 0x0F, bits.OnesCount8(res)%2 == 0, uint16(c.a)+uint16(value)+uint16(carry) > 0xFF)
	c.a = res
}

func sub(c *CPU, op string) {
	value := c.getOp(op)
	res := c.a - value
	c.setFlags(res&0x80 == 0x80, res == 0, (c.a&0x0F) >= (value&0x0F), bits.OnesCount8(res)%2 == 0, uint16(c.a) < uint16(value))
	c.a = res
}

func sbb(c *CPU, op string) {
	value := c.getOp(op)
	carry := c.getCYF()
	res := c.a - value - carry
	c.setFlags(res&0x80 == 0x80, res == 0, (c.a&0x0F) >= (value&0x0F)+carry, bits.OnesCount8(res)%2 == 0, uint16(c.a) < uint16(value)+uint16(carry))
	c.a = res
}

func ana(c *CPU, op string) {
	value := c.getOp(op)
	res := c.a & value
	c.setFlags(res&0x80 == 0x80, res == 0, (c.a|value)&0x08 != 0, bits.OnesCount8(res)%2 == 0, false)
	c.a = res
}

func xra(c *CPU, op string) {
	value := c.getOp(op)
	res := c.a ^ value
	c.setFlags(res&0x80 == 0x80, res == 0, false, bits.OnesCount8(res)%2 == 0, false)
	c.a = res
}

func ora(c *CPU, op string) {
	value := c.getOp(op)
	res := c.a | value
	c.setFlags(res&0x80 == 0x80, res == 0, false, bits.OnesCount8(res)%2 == 0, false)
	c.a = res
}

func cmp(c *CPU, op string) {
	value := c.getOp(op)
	res := c.a - value
	c.setFlags(res&0x80 == 0x80, res == 0, c.a&0x0F >= value&0x0F, bits.OnesCount8(res)%2 == 0, c.a < value)
}

func hlt(c *CPU, _ string) {
	// TODO: implement
}

func xthl(c *CPU, _ string) {
	lo := c.ReadMem(c.sp)
	hi := c.ReadMem(c.sp + 1)
	c.WriteMem(c.sp, c.l)
	c.WriteMem(c.sp+1, c.h)
	c.l = lo
	c.h = hi
}

func cpi(c *CPU, _ string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	res := c.a - value

	c.setFlags(res&0x80 == 0x80, res == 0, c.a&0x0F >= value&0x0F, bits.OnesCount8(res)%2 == 0, c.a < value)
}

func rst(c *CPU, op string) {
	c.push(c.pc + 1)
	c.pc = uint16(lib.Must(strconv.ParseUint(op, 16, 16)) * 8)
}

func adi(c *CPU, _ string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	res := c.a + value
	c.setFlags(res&0x80 == 0x80, res == 0, c.a&0x0F+(value)&0x0F > 0x0F, bits.OnesCount8(res)%2 == 0, uint16(c.a)+uint16(value) > 0xFF)
	c.a = res
}

func aci(c *CPU, _ string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	carry := c.getCYF()
	res := c.a + value + carry
	c.setFlags(res&0x80 == 0x80, res == 0, c.a&0x0F+(value)&0x0F+carry > 0x0F, bits.OnesCount8(res)%2 == 0, uint16(c.a)+uint16(value)+uint16(carry) > 0xFF)
	c.a = res
}

func sui(c *CPU, _ string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	res := c.a - value
	c.setFlags(res&0x80 == 0x80, res == 0, (c.a&0x0F) >= (value&0x0F), bits.OnesCount8(res)%2 == 0, uint16(c.a) < uint16(value))
	c.a = res
}

func sbi(c *CPU, _ string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	carry := c.getCYF()
	res := c.a - value - carry
	c.setFlags(res&0x80 == 0x80, res == 0, (c.a&0x0F) >= (value&0x0F)+carry, bits.OnesCount8(res)%2 == 0, uint16(c.a) < uint16(value)+uint16(carry))
	c.a = res
}

func ani(c *CPU, _ string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	res := c.a & value

	c.setFlags(res&0x80 == 0x80, res == 0, (c.a|value)&0x08 != 0, bits.OnesCount8(res)%2 == 0, false)
	c.a = res
}

func xri(c *CPU, _ string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	res := c.a ^ value
	c.setFlags(res&0x80 == 0x80, res == 0, false, bits.OnesCount8(res)%2 == 0, false)
	c.a = res
}

func ori(c *CPU, _ string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	res := c.a | value
	c.setFlags(res&0x80 == 0x80, res == 0, false, bits.OnesCount8(res)%2 == 0, false)
	c.a = res
}

func portIn(c *CPU, _ string) {
	portNumber := c.ReadMem(c.pc + 1)

	switch portNumber {
	case 3:
		c.a = uint8(c.sr >> (8 - c.so))
	default:
		c.a = c.ioPorts[portNumber]
	}
}

func portOut(c *CPU, _ string) {
	portNumber := c.ReadMem(c.pc + 1)

	switch portNumber {
	case 0:
		c.Running = false
	case 1:
		switch c.c {
		case 2:
			fmt.Print(string(c.e))
		case 9:
			addr := uint16(c.d)<<8 | uint16(c.e)
			for offset := uint16(0); ; offset++ {
				b := c.ReadMem(addr + offset)
				if b == '$' {
					break
				}

				fmt.Print(string(b))
			}

		default:
			fmt.Printf("unimplemented out operation for port 1: %02x\n", c.c)
		}
	case 2:
		c.so = c.a & 0x7
	case 3:
		// TODO: sound
	case 4:
		c.sr = uint16(c.a)<<8 | c.sr>>8
	case 5:
		// TODO: sound
	case 6: // NOP for watchdog
	default:
		fmt.Printf("unimplemented out port: %02x\n", portNumber)
	}
}
