package arcade

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cterence/space-invaders/internal/arcade/cpu"
	"github.com/cterence/space-invaders/internal/arcade/memory"
)

type arcade struct {
	cpu    *cpu.CPU
	memory *memory.Memory
}

func Run(ctx context.Context, romDirPath string) error {
	a := arcade{
		cpu:    &cpu.CPU{},
		memory: &memory.Memory{},
	}

	a.memory.Init()
	a.cpu.Init()

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

	for a.cpu.Running {
		a.cpu.Step()
	}

	return nil
}
