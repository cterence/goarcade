package arcade

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cterence/space-invaders/internal/arcade/cpu"
	"github.com/cterence/space-invaders/internal/arcade/memory"
)

const (
	cpuFreq = 2000000
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

func Run(ctx context.Context, romDirPath string, options ...Option) error {
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

	// Load Space Invaders rom
	if romDirPath == "" {
		return errors.New("empty rom directory path")
	}

	i := 0

	for _, id := range []string{"h", "g", "f", "e"} {
		romBytes, err := os.ReadFile(filepath.Join(romDirPath, "invaders."+id))
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
