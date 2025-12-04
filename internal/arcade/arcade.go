package arcade

import (
	"context"
	"errors"
	"fmt"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/cterence/space-invaders/internal/arcade/cpu"
	"github.com/cterence/space-invaders/internal/arcade/memory"
)

const (
	CPU_FREQ = 2_000_000
)

type arcade struct {
	cpu    *cpu.CPU
	memory *memory.Memory

	cpuOpts []cpu.Option

	cpuSC uint64
	stop  uint64
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

	i := 0x100

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

	// inject "out 0,a" at 0x0000 (signal to stop the test)
	a.memory.Write(0x0000, 0xD3)
	a.memory.Write(0x0001, 0x00)

	// inject "out 1,a" at 0x0005 (signal to output some characters)
	a.memory.Write(0x0005, 0xD3)
	a.memory.Write(0x0006, 0x01)
	a.memory.Write(0x0007, 0xC9)

	// lastCPUPeriod := time.Now()
	// loggedThrottled := false
	milestone := uint64(1000000)
	sc := uint64(0)

	for a.cpu.Running && (a.stop == 0 || a.cpuSC < a.stop) {
		// if sc >= CPU_FREQ {
		// 	if time.Since(lastCPUPeriod) <= time.Second {
		// 		if !loggedThrottled {
		// 			fmt.Println("============ throttled ===============")

		// 			loggedThrottled = true
		// 		}

		// 		continue
		// 	} else {
		// 		lastCPUPeriod = time.Now()
		// 		loggedThrottled = false
		// 		sc = 0
		// 	}
		// }
		a.cpuSC += uint64(a.cpu.Step())

		sc += a.cpuSC
		if a.cpuSC > milestone {
			// fmt.Printf("States milestone reached: %d\n", milestone)
			milestone += 1000000
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
	c.Init()
	c.ReadMem = func(offset uint16) uint8 { return romBytes[i+offset] }

	for i < uint16(len(romBytes)) {
		inst, _, _, instLength, _ := c.DecodeInst()

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

		fmt.Println(b.String())

		i += instLength
	}

	return nil
}
