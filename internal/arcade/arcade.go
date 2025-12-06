package arcade

import (
	"context"
	"errors"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/cterence/space-invaders/internal/arcade/apu"
	"github.com/cterence/space-invaders/internal/arcade/cpu"
	"github.com/cterence/space-invaders/internal/arcade/memory"
	"github.com/cterence/space-invaders/internal/arcade/ui"
)

const (
	CPU_TPS           = 2_000_000
	FPS               = 60
	CPU_TPS_PER_FRAME = CPU_TPS / FPS
)

type arcade struct {
	cpu    *cpu.CPU
	memory *memory.Memory
	ui     *ui.UI
	apu    *apu.APU

	cpuOpts []cpu.Option

	cpm        bool
	headless   bool
	unthrottle bool
	noAudio    bool
	soundDir   string
}

type Option func(*arcade)

func WithDebug(debug bool) Option {
	return func(a *arcade) {
		a.cpuOpts = append(a.cpuOpts, cpu.WithDebug(debug))
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

func WithSoundDir(soundDir string) Option {
	return func(a *arcade) {
		a.soundDir = soundDir
	}
}

func WithNoAudio(noAudio bool) Option {
	return func(a *arcade) {
		a.noAudio = noAudio
	}
}

func WithUnthrottle(unthrottle bool) Option {
	return func(a *arcade) {
		a.unthrottle = unthrottle
	}
}

func Run(ctx context.Context, romPaths []string, options ...Option) error {
	aCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	a := arcade{
		cpu:    &cpu.CPU{},
		memory: &memory.Memory{},
		ui:     &ui.UI{},
		apu:    &apu.APU{},
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

		a.ui.Init(cancel)
		defer a.ui.Close()

		trapSigInt(cancel)

		if !a.noAudio {
			a.apu.Init(a.soundDir)

			defer a.apu.Close()
		}
	}

	a.cpu.ReadMem = a.memory.Read
	a.cpu.WriteMem = a.memory.Write
	a.cpu.PlaySound = a.apu.PlaySound
	a.cpu.StartSoundLoop = a.apu.StartLoop
	a.cpu.StopSoundLoop = a.apu.StopLoop
	a.ui.ReadMem = a.memory.Read
	a.ui.RequestInterrupt = a.cpu.RequestInterrupt
	a.ui.SendInput = a.cpu.SendInput

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

	cpuCycles := uint64(0)

	frameTicker := time.NewTicker(time.Second / FPS)
	defer frameTicker.Stop()

	for {
		if !a.cpu.Running {
			return nil
		}

		if a.unthrottle {
			cpuCycles += uint64(a.cpu.Step())

			if cpuCycles >= CPU_TPS_PER_FRAME {
				cpuCycles = 0

				if !a.headless {
					a.ui.Step()
				}
			}
		} else {
			select {
			case <-frameTicker.C:
				for cpuCycles < CPU_TPS_PER_FRAME {
					cpuCycles += uint64(a.cpu.Step())
				}

				cpuCycles = 0

				if !a.headless {
					a.ui.Step()
				}
			case <-aCtx.Done():
				return nil
			}
		}
	}
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

func trapSigInt(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		cancel()
	}()
}
