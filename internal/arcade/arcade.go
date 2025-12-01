package arcade

import (
	"context"
	"errors"
	"fmt"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	"github.com/cterence/space-invaders/internal/arcade/cpu"
	"github.com/cterence/space-invaders/internal/arcade/memory"
)

const (
	cpuFreq = 2_000_000
)

type arcade struct {
	cpu    *cpu.CPU
	memory *memory.Memory

	cpuOpts []cpu.Option

	cpuSC uint64
}

type Option func(*arcade)

func WithDebug(debug bool) Option {
	return func(a *arcade) {
		a.cpuOpts = append(a.cpuOpts, cpu.WithDebug(debug))
	}
}

func Run(ctx context.Context, romPaths []string, options ...Option) error {
	a := arcade{
		cpu:    &cpu.CPU{},
		memory: &memory.Memory{},
		cpuSC:  0,
	}

	for _, o := range options {
		o(&a)
	}

	a.memory.Init()
	a.cpu.Init(a.cpuOpts...)

	a.cpu.ReadMem = a.memory.Read
	a.cpu.WriteMem = a.memory.Write

	if len(romPaths) == 0 {
		return errors.New("no rom passed to emulator")
	}

	i := 0

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

	lastCPUPeriod := time.Now()
	loggedThrottled := false

	for a.cpu.Running {
		if a.cpuSC >= cpuFreq {
			if time.Since(lastCPUPeriod) <= time.Second {
				if !loggedThrottled {
					fmt.Println("============ throttled ===============")

					loggedThrottled = true
				}

				continue
			} else {
				lastCPUPeriod = time.Now()
				a.cpuSC -= cpuFreq
				loggedThrottled = false
			}
		}

		a.cpuSC += a.cpu.Step()
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
	c.Init()
	c.ReadMem = func(offset uint16) uint8 { return romBytes[i+offset] }

	for i < uint16(len(romBytes)) {
		inst, op1, op2, instLength, _ := c.DecodeInst()

		var b strings.Builder

		b.WriteString(fmt.Sprintf("%04x: ", i))

		if instLength == 1 {
			b.WriteString(fmt.Sprintf("%02x          ", romBytes[i]))
		}

		if instLength == 2 {
			b.WriteString(fmt.Sprintf("%02x %02x       ", romBytes[i], romBytes[i+1]))
		}

		if instLength == 3 {
			b.WriteString(fmt.Sprintf("%02x %02x %02x    ", romBytes[i], romBytes[i+1], romBytes[i+2]))
		}

		b.WriteString(inst)

		if op1 != "" {
			b.WriteString(" " + op1)

			if op2 != "" {
				b.WriteString(", " + op2)
			}
		}

		fmt.Println(b.String())

		i += instLength
	}

	return nil
}
