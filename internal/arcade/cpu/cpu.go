package cpu

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"

	"github.com/cterence/goarcade/internal/arcade/config"
)

type bus interface {
	Read(addr uint16) uint8
	Write(addr uint16, value uint8)
}

type apu interface {
	PlaySound(id uint8)
	StartSoundLoop(id uint8)
	StopSoundLoop(id uint8)
}

type CPU struct {
	Bus bus
	APU apu
	state

	Running bool
}

type state struct {
	InPorts map[int][8]config.Port

	// Cycle counter
	Cyc uint64
	// Program counter
	PC uint16
	// Stack pointer
	SP uint16

	// Shift register
	SR uint16

	// IO ports
	ioPorts [8]uint8
	Debug   bool

	// Shift offset
	SO uint8

	// Registers
	A uint8
	F uint8
	B uint8
	C uint8
	D uint8
	E uint8
	H uint8
	L uint8

	// Interrupt switch
	Interrupts bool
}

type Option func(*CPU)

func WithDebug(debug bool) Option {
	return func(c *CPU) {
		c.Debug = debug
	}
}

func (c *CPU) Init(pc uint16, options ...Option) {
	c.Running = true
	c.PC = pc
	c.SP = 0
	c.SR = 0
	c.Cyc = 0
	c.SO = 0
	c.B = 0
	c.C = 0
	c.D = 0
	c.E = 0
	c.H = 0
	c.L = 0
	c.A = 0
	c.F = 2
	c.Interrupts = false

	for id, port := range c.InPorts {
		for _, portBit := range port {
			if portBit.Active {
				c.ioPorts[id] |= 1 << portBit.Bit
			} else {
				c.ioPorts[id] &= ^(1 << portBit.Bit)
			}
		}
	}

	for _, o := range options {
		o(c)
	}
}

func (c *CPU) String() string {
	var b strings.Builder

	b.WriteString("PC: " + fmt.Sprintf("%04X", c.PC))
	b.WriteString(", AF: " + fmt.Sprintf("%04X", uint16(c.A)<<8|uint16(c.F)))
	b.WriteString(", BC: " + fmt.Sprintf("%04X", uint16(c.B)<<8|uint16(c.C)))
	b.WriteString(", DE: " + fmt.Sprintf("%04X", uint16(c.D)<<8|uint16(c.E)))
	b.WriteString(", HL: " + fmt.Sprintf("%04X", uint16(c.H)<<8|uint16(c.L)))
	b.WriteString(", SP: " + fmt.Sprintf("%04X", c.SP))
	b.WriteString(", CYC: " + strconv.FormatUint(c.Cyc, 10))

	return b.String()
}

func (c *CPU) Step() uint8 {
	prevPC := c.PC

	inst := InstByOpcode[c.Bus.Read(c.PC)]

	if c.Debug {
		fmt.Printf("%s (%02X %02X %02X %02X) %-13s\n", c, c.Bus.Read(c.PC), c.Bus.Read(c.PC+1), c.Bus.Read(c.PC+2), c.Bus.Read(c.PC+3), inst.Name+" "+inst.Op1+" "+inst.Op2)
		// fmt.Printf("%s (%02X %02X %02X %02X)\n", c, c.Bus.Read(c.pc), c.Bus.Read(c.pc+1), c.Bus.Read(c.pc+2), c.Bus.Read(c.pc+3))
	}

	inst.exec(c, inst.Op1)

	if prevPC == c.PC {
		c.PC += uint16(inst.Length)
	}

	c.Cyc += uint64(inst.Cycles)

	return inst.Cycles
}

func (c *CPU) RequestInterrupt(num uint8) {
	if c.Interrupts {
		c.push(c.PC)
		c.PC = uint16(8 * num)
	}
}

func (c *CPU) SendInput(port, bit uint8, value bool) {
	if value {
		c.ioPorts[port] |= 1 << bit
	} else {
		c.ioPorts[port] &= ^(1 << bit)
	}
}

func (c *CPU) SaveState() ([]uint8, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	err := enc.Encode(c.state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode state: %w", err)
	}

	return buf.Bytes(), nil
}

func (c *CPU) LoadState(stateBytes []uint8) error {
	var buf bytes.Buffer

	_, err := buf.Write(stateBytes)
	if err != nil {
		return fmt.Errorf("failed to load state in buffer: %w", err)
	}

	enc := gob.NewDecoder(&buf)

	var s state

	err = enc.Decode(&s)
	if err != nil {
		return fmt.Errorf("failed to decode state: %w", err)
	}

	c.state = s

	return nil
}
