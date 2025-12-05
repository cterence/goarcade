package cpu

import (
	"fmt"
	"math"
	"math/bits"
	"strconv"

	"github.com/cterence/space-invaders/internal/arcade/lib"
)

func nop(*CPU, string, string) {}

func ldax(c *CPU, op1, op2 string) {
	c.a = c.ReadMem(c.getDoubleOp(op1))
}

func stax(c *CPU, op1, op2 string) {
	c.WriteMem(c.getDoubleOp(op1), c.a)
}

func inx(c *CPU, op1, op2 string) {
	c.setDoubleOp(op1, c.getDoubleOp(op1)+1)
}

func dcx(c *CPU, op1, op2 string) {
	c.setDoubleOp(op1, c.getDoubleOp(op1)-1)
}

func inr(c *CPU, op1, op2 string) {
	value := c.getOp(op1)
	res := value + 1
	c.setOp(op1, res)
	c.setFlags(res&0x80 == 0x80, res == 0, value&0x0F+1 > 0x0F, bits.OnesCount8(res)%2 == 0, c.getCYF() == 1)
}

func dcr(c *CPU, op1, op2 string) {
	value := c.getOp(op1)
	res := value - 1
	c.setOp(op1, res)
	c.setFlags(res&0x80 == 0x80, res == 0, value&0x0F >= 1, bits.OnesCount8(res)%2 == 0, c.getCYF() == 1)
}

func dad(c *CPU, op1, op2 string) {
	res := uint32(c.getHL()) + uint32(c.getDoubleOp(op1))

	c.setCYF(res > math.MaxUint16)
	c.setHL(uint16(res))
}

func lxi(c *CPU, op1, op2 string) {
	imm1, imm2 := c.ReadMem(c.pc+1), c.ReadMem(c.pc+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	c.setDoubleOp(op1, imm16)
}

func (c *CPU) pop() uint16 {
	value := uint16(c.ReadMem(c.sp+1))<<8 | uint16(c.ReadMem(c.sp))
	c.sp += 2

	return value
}

func popOp(c *CPU, op1, op2 string) {
	value := c.pop()
	c.setDoubleOp(op1, value)
}

func (c *CPU) push(addr uint16) {
	c.sp -= 2
	c.WriteMem(c.sp, uint8(addr))
	c.WriteMem(c.sp+1, uint8(addr>>8))
}

func pushOp(c *CPU, op1, op2 string) {
	c.sp -= 2
	value := c.getDoubleOp(op1)
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

func mvi(c *CPU, op1, op2 string) {
	imm1 := c.ReadMem(c.pc + 1)
	c.setOp(op1, imm1)
}

func rlc(c *CPU, op1, op2 string) {
	sb := c.a & 0x80 >> 7
	c.setCYF(sb == 1)
	c.a = c.a<<1 | sb
}
func rrc(c *CPU, op1, op2 string) {
	sb := c.a & 0x1
	c.setCYF(sb == 1)
	c.a = c.a>>1 | sb<<7
}

func ral(c *CPU, op1, op2 string) {
	sb := c.a & 0x80 >> 7
	c.a = c.a<<1 | c.getCYF()
	c.setCYF(sb == 1)
}

func rar(c *CPU, op1, op2 string) {
	sb := c.a & 0x1
	c.a = c.a>>1 | c.getCYF()<<7
	c.setCYF(sb == 1)
}

func shld(c *CPU, op1, op2 string) {
	imm1, imm2 := c.ReadMem(c.pc+1), c.ReadMem(c.pc+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	c.WriteMem(imm16, c.getL())
	c.WriteMem(imm16+1, c.getH())
}

func daa(c *CPU, op1, op2 string) {
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

func lhld(c *CPU, op1, op2 string) {
	imm1, imm2 := c.ReadMem(c.pc+1), c.ReadMem(c.pc+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	c.setHL(uint16(c.ReadMem(imm16+1))<<8 | uint16(c.ReadMem(imm16)))
}

func sta(c *CPU, op1, op2 string) {
	imm1, imm2 := c.ReadMem(c.pc+1), c.ReadMem(c.pc+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	c.WriteMem(imm16, c.a)
}

func lda(c *CPU, op1, op2 string) {
	imm1, imm2 := c.ReadMem(c.pc+1), c.ReadMem(c.pc+2)
	imm16 := uint16(imm2)<<8 | uint16(imm1)
	c.a = c.ReadMem(imm16)
}

func add(c *CPU, op1, op2 string) {
	value := c.getOp(op1)
	res := c.a + value
	c.setFlags(res&0x80 == 0x80, res == 0, c.a&0x0F+(value)&0x0F > 0x0F, bits.OnesCount8(res)%2 == 0, uint16(c.a)+uint16(value) > 0xFF)
	c.a = res
}

func adc(c *CPU, op1, op2 string) {
	value := c.getOp(op1)
	carry := c.getCYF()
	res := c.a + value + carry
	c.setFlags(res&0x80 == 0x80, res == 0, c.a&0x0F+(value)&0x0F+carry > 0x0F, bits.OnesCount8(res)%2 == 0, uint16(c.a)+uint16(value)+uint16(carry) > 0xFF)
	c.a = res
}

func sub(c *CPU, op1, op2 string) {
	value := c.getOp(op1)
	res := c.a - value
	c.setFlags(res&0x80 == 0x80, res == 0, (c.a&0x0F) >= (value&0x0F), bits.OnesCount8(res)%2 == 0, uint16(c.a) < uint16(value))
	c.a = res
}

func sbb(c *CPU, op1, op2 string) {
	value := c.getOp(op1)
	carry := c.getCYF()
	res := c.a - value - carry
	c.setFlags(res&0x80 == 0x80, res == 0, (c.a&0x0F) >= (value&0x0F)+carry, bits.OnesCount8(res)%2 == 0, uint16(c.a) < uint16(value)+uint16(carry))
	c.a = res
}

func ana(c *CPU, op1, op2 string) {
	value := c.getOp(op1)
	res := c.a & value
	c.setFlags(res&0x80 == 0x80, res == 0, (c.a|value)&0x08 != 0, bits.OnesCount8(res)%2 == 0, false)
	c.a = res
}

func xra(c *CPU, op1, op2 string) {
	value := c.getOp(op1)
	res := c.a ^ value
	c.setFlags(res&0x80 == 0x80, res == 0, false, bits.OnesCount8(res)%2 == 0, false)
	c.a = res
}

func ora(c *CPU, op1, op2 string) {
	value := c.getOp(op1)
	res := c.a | value
	c.setFlags(res&0x80 == 0x80, res == 0, false, bits.OnesCount8(res)%2 == 0, false)
	c.a = res
}

func cmp(c *CPU, op1, op2 string) {
	value := c.getOp(op1)
	res := c.a - value
	c.setFlags(res&0x80 == 0x80, res == 0, c.a&0x0F >= value&0x0F, bits.OnesCount8(res)%2 == 0, c.a < value)
}

func hlt(c *CPU, op1, op2 string) {
	// TODO: implement
}

func xthl(c *CPU, op1, op2 string) {
	lo := c.ReadMem(c.sp)
	hi := c.ReadMem(c.sp + 1)
	c.WriteMem(c.sp, c.l)
	c.WriteMem(c.sp+1, c.h)
	c.l = lo
	c.h = hi
}

func cpi(c *CPU, op1, op2 string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	res := c.a - value

	c.setFlags(res&0x80 == 0x80, res == 0, c.a&0x0F >= value&0x0F, bits.OnesCount8(res)%2 == 0, c.a < value)
}

func rst(c *CPU, op1, op2 string) {
	c.pc = uint16(lib.Must(strconv.ParseUint(op1, 16, 16)) * 8)
}

func adi(c *CPU, op1, op2 string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	res := c.a + value
	c.setFlags(res&0x80 == 0x80, res == 0, c.a&0x0F+(value)&0x0F > 0x0F, bits.OnesCount8(res)%2 == 0, uint16(c.a)+uint16(value) > 0xFF)
	c.a = res
}

func aci(c *CPU, op1, op2 string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	carry := c.getCYF()
	res := c.a + value + carry
	c.setFlags(res&0x80 == 0x80, res == 0, c.a&0x0F+(value)&0x0F+carry > 0x0F, bits.OnesCount8(res)%2 == 0, uint16(c.a)+uint16(value)+uint16(carry) > 0xFF)
	c.a = res
}

func sui(c *CPU, op1, op2 string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	res := c.a - value
	c.setFlags(res&0x80 == 0x80, res == 0, (c.a&0x0F) >= (value&0x0F), bits.OnesCount8(res)%2 == 0, uint16(c.a) < uint16(value))
	c.a = res
}

func sbi(c *CPU, op1, op2 string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	carry := c.getCYF()
	res := c.a - value - carry
	c.setFlags(res&0x80 == 0x80, res == 0, (c.a&0x0F) >= (value&0x0F)+carry, bits.OnesCount8(res)%2 == 0, uint16(c.a) < uint16(value)+uint16(carry))
	c.a = res
}

func ani(c *CPU, op1, op2 string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	res := c.a & value

	c.setFlags(res&0x80 == 0x80, res == 0, (c.a|value)&0x08 != 0, bits.OnesCount8(res)%2 == 0, false)
	c.a = res
}

func xri(c *CPU, op1, op2 string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	res := c.a ^ value
	c.setFlags(res&0x80 == 0x80, res == 0, false, bits.OnesCount8(res)%2 == 0, false)
	c.a = res
}

func ori(c *CPU, op1, op2 string) {
	imm1 := c.ReadMem(c.pc + 1)
	value := imm1
	res := c.a | value
	c.setFlags(res&0x80 == 0x80, res == 0, false, bits.OnesCount8(res)%2 == 0, false)
	c.a = res
}

func portIn(c *CPU, op1, op2 string) {}

func portOut(c *CPU, op1, op2 string) {
	portNumber := c.ReadMem(c.pc + 1)

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
