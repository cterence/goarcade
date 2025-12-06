package arcade

import (
	"context"
	"errors"
	"fmt"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/cterence/space-invaders/internal/arcade/cpu"
	"github.com/cterence/space-invaders/internal/arcade/memory"
	"github.com/cterence/space-invaders/internal/arcade/ui"
)

const (
	CPU_FREQ = 2_000_000
	UI_FREQ  = 60
)

type arcade struct {
	cpu    *cpu.CPU
	memory *memory.Memory
	ui     *ui.UI

	cpuOpts []cpu.Option

	cpuSC uint64

	stop       uint64
	cpm        bool
	headless   bool
	unthrottle bool
}

type Option func(*arcade)

func WithDebug(debug bool) Option {
	return func(a *arcade) {
		a.cpuOpts = append(a.cpuOpts, cpu.WithDebug(debug))
	}
}

func WithStop(stop uint64) Option {
	return func(a *arcade) {
		a.stop = stop
	}
}

func WithCPM(cpm bool) Option {
	return func(a *arcade) {
		a.cpm = cpm
	}
}

func WithHeadless(headless bool) Option {
	return func(a *arcade) {
		a.headless = headless
	}
}

func WithUnthrottle(unthrottle bool) Option {
	return func(a *arcade) {
		a.unthrottle = unthrottle
	}
}

func Run(ctx context.Context, romPaths []string, options ...Option) error {
	a := arcade{
		cpu:    &cpu.CPU{},
		memory: &memory.Memory{},
		ui:     &ui.UI{},
		cpuSC:  0,
	}

	for _, o := range options {
		o(&a)
	}

	cpuPC := uint16(0)

	if a.cpm {
		cpuPC = 0x100
	}

	a.cpu.Init(cpuPC, a.cpuOpts...)
	a.memory.Init()

	if !a.headless {
		defer binsdl.Load().Unload()

		a.ui.Init()
		defer a.ui.Close()
	}

	a.cpu.ReadMem = a.memory.Read
	a.cpu.WriteMem = a.memory.Write
	a.ui.ReadMem = a.memory.Read

	if len(romPaths) == 0 {
		return errors.New("no rom passed to emulator")
	}

	i := 0

	if a.cpm {
		i = 0x100
	}

	for _, p := range romPaths {
		romBytes, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("failed to read rom file: %w", err)
		}

		for _, b := range romBytes {
			a.memory.Write(uint16(i), b)
			i++
		}
	}

	if a.cpm {
		// inject "out 0,a" at 0x0000 (signal to stop the test)
		a.memory.Write(0x0000, 0xD3)
		a.memory.Write(0x0001, 0x00)

		// inject "out 1,a" at 0x0005 (signal to output some characters)
		a.memory.Write(0x0005, 0xD3)
		a.memory.Write(0x0006, 0x01)
		a.memory.Write(0x0007, 0xC9)
	}

	lastCPUPeriod := time.Now().Add(-time.Second)
	lastUIStep := time.Now().Add(-time.Second)

	sc := uint64(0)
	for a.cpu.Running && (a.stop == 0 || a.cpuSC < a.stop) {
		if a.unthrottle || sc < CPU_FREQ {
			sc += uint64(a.cpu.Step())
			a.cpuSC += sc
		}

		if !a.unthrottle && time.Since(lastCPUPeriod) >= time.Second {
			lastCPUPeriod = time.Now()
			sc = 0
		}

		if !a.headless && time.Since(lastUIStep) >= time.Second/UI_FREQ {
			a.ui.Step()

			lastUIStep = time.Now()
		}
	}

	return nil
}

func Disassemble(romPaths []string) error {
	if len(romPaths) == 0 {
		return errors.New("no rom passed to emulator")
	}

	romBytes := []byte{}

	for _, p := range romPaths {
		romPartBytes, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("failed to read rom file: %w", err)
		}

		romBytes = append(romBytes, romPartBytes...)
	}

	var i uint16

	c := cpu.CPU{}
	c.Init(0)
	c.ReadMem = func(_ uint16) uint8 { return romBytes[i] }

	for i < uint16(len(romBytes)) {
		inst := cpu.InstByOpcode[romBytes[i]]

		var b strings.Builder

		b.WriteString(fmt.Sprintf("%04x: ", i))

		if inst.Length == 1 {
			b.WriteString(fmt.Sprintf("%02x          ", romBytes[i]))
		}

		if inst.Length == 2 {
			b.WriteString(fmt.Sprintf("%02x %02x       ", romBytes[i], romBytes[i+1]))
		}

		if inst.Length == 3 {
			b.WriteString(fmt.Sprintf("%02x %02x %02x    ", romBytes[i], romBytes[i+1], romBytes[i+2]))
		}

		b.WriteString(inst.Name + " " + inst.Op1 + " " + inst.Op2)

		fmt.Println(b.String())

		i += uint16(inst.Length)
	}

	return nil
}
