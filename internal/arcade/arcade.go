package arcade

import (
	"context"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
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
	cancel context.CancelFunc

	cpuOpts []cpu.Option

	cpm        bool
	headless   bool
	unthrottle bool
	noAudio    bool
	soundDir   string
	saveState  string

	romHash string
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

func WithSaveState(saveState string) Option {
	return func(a *arcade) {
		a.saveState = saveState
	}
}

func (a *arcade) Reset() {
	cpuPC := uint16(0)

	if a.cpm {
		cpuPC = 0x100
	}

	a.cpu.Init(cpuPC, a.cpuOpts...)

	if !a.headless {
		a.ui.Init()

		if !a.noAudio {
			a.apu.Init(a.soundDir)
		}
	}
}

func Run(ctx context.Context, romPaths []string, options ...Option) error {
	aCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	a := arcade{
		cpu:     &cpu.CPU{},
		memory:  &memory.Memory{},
		ui:      &ui.UI{},
		apu:     &apu.APU{},
		cancel:  cancel,
		romHash: romHash(romPaths)[:8],
	}

	a.cpu.Bus = a.memory
	a.cpu.APU = a.apu

	a.ui.Arcade = &a
	a.ui.Bus = a.memory
	a.ui.CPU = a.cpu
	a.ui.APU = a.apu

	for _, o := range options {
		o(&a)
	}

	if len(romPaths) == 0 {
		return errors.New("no rom passed to emulator")
	}

	if !a.headless {
		defer binsdl.Load().Unload()
		defer a.apu.Close()
		defer a.ui.Close()

		trapSigInt(cancel)
	}

	a.Reset()

	i := 0

	if a.cpm {
		i = 0x100
	}

	for _, p := range romPaths {
		// Must not write ROM to 0x2000 - 0x3FFF
		if i == 0x2000 {
			i = 0x4000
		}

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

	if a.saveState != "" {
		if err := a.LoadState(); err != nil {
			return err
		}
	}

	cpuCycles := uint64(0)

	frameTicker := time.NewTicker(time.Second / FPS)
	defer frameTicker.Stop()

	if a.unthrottle {
		for a.cpu.Running {
			cpuCycles += uint64(a.cpu.Step())

			if cpuCycles >= CPU_TPS_PER_FRAME {
				cpuCycles = 0

				if !a.headless {
					a.ui.Step()
				}
			}
		}
	}

	for a.cpu.Running {
		select {
		case <-frameTicker.C:
			for !a.ui.Paused && cpuCycles < CPU_TPS_PER_FRAME {
				cpuCycles += uint64(a.cpu.Step())
			}

			if cpuCycles >= CPU_TPS_PER_FRAME {
				cpuCycles = 0
			}

			if !a.headless {
				a.ui.Step()
			}
		case <-aCtx.Done():
			return nil
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

type state struct {
	CPU    []uint8
	Memory []uint8
}

func (a *arcade) SaveState() error {
	cpu, err := a.cpu.SaveState()
	if err != nil {
		return err
	}

	memory := make([]uint8, 0x10000)

	for addr := range memory {
		memory[addr] = a.memory.Read(uint16(addr))
	}

	s := state{
		CPU:    cpu,
		Memory: memory,
	}

	f, err := os.Create(a.romHash + ".bin")
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)

	return enc.Encode(s)
}

func (a *arcade) LoadState() error {
	fileName := a.saveState
	if fileName == "" {
		fileName = a.romHash + ".bin"
	}

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	var s state

	dec := gob.NewDecoder(f)
	dec.Decode(&s)

	if err := a.cpu.LoadState(s.CPU); err != nil {
		return fmt.Errorf("failed to load CPU state: %w", err)
	}

	for addr, b := range s.Memory {
		a.memory.Write(uint16(addr), b)
	}

	return nil
}

func (a *arcade) Shutdown() {
	a.cancel()
}

func trapSigInt(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		cancel()
	}()
}

func romHash(romPaths []string) string {
	h := sha256.New()

	for _, f := range romPaths {
		h.Write([]uint8(filepath.Base(f)))
	}

	return hex.EncodeToString(h.Sum(nil))
}
