package arcade

import (
	"archive/zip"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/cterence/goarcade/internal/arcade/apu"
	"github.com/cterence/goarcade/internal/arcade/cpu"
	"github.com/cterence/goarcade/internal/arcade/gamespec"
	"github.com/cterence/goarcade/internal/arcade/lib"
	"github.com/cterence/goarcade/internal/arcade/memory"
	"github.com/cterence/goarcade/internal/arcade/ui"
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

	soundDir  string
	saveState string

	romPath string

	cpuOpts []cpu.Option

	cpm        bool
	headless   bool
	unthrottle bool
	noAudio    bool
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

func Run(ctx context.Context, romPath string, options ...Option) error {
	aCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if len(romPath) == 0 {
		return errors.New("no rom passed to emulator")
	}

	a := arcade{
		cpu:     &cpu.CPU{},
		memory:  &memory.Memory{},
		ui:      &ui.UI{},
		apu:     &apu.APU{},
		cancel:  cancel,
		romPath: romPath,
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

	if !a.headless {
		defer binsdl.Load().Unload()
		defer a.apu.Close()
		defer a.ui.Close()

		trapSigInt(cancel)
	}

	i := 0

	if filepath.Ext(romPath) == ".zip" {
		r, err := zip.OpenReader(romPath)
		if err != nil {
			return fmt.Errorf("failed to open zip archive: %w", err)
		}
		defer lib.DeferErr(r.Close())

		settings, err := gamespec.GetGameSettings(filepath.Base(romPath))
		if err != nil {
			return err
		}

		a.ui.ColorOverlays = settings.ColorOverlays

		a.Reset()

		for _, p := range settings.ROMParts {
			if err := a.LoadBytes(p.StartAddr, p.ExpectedSize, lib.Must(GetFileBytesFromZip(r.File, p.FileName))); err != nil {
				return fmt.Errorf("failed to load %s", p.FileName)
			}
		}

		if len(settings.ColorPROMs) > 0 {
			a.ui.ColorPROMs = make([][]uint8, len(settings.ColorPROMs))
		}

		for i, p := range settings.ColorPROMs {
			a.ui.ColorPROMs[i] = make([]uint8, p.ExpectedSize)
			copy(a.ui.ColorPROMs[i][:], lib.Must(GetFileBytesFromZip(r.File, p.FileName))[:])
		}
	} else {
		if a.cpm {
			i = 0x100
		}

		romBytes, err := os.ReadFile(romPath)
		if err != nil {
			return fmt.Errorf("failed to read rom file: %w", err)
		}

		a.Reset()

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

func GetFileBytesFromZip(files []*zip.File, fileName string) ([]uint8, error) {
	fileIndex := slices.IndexFunc(files, func(f *zip.File) bool { return f.Name == fileName })

	f, err := files[fileIndex].Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer lib.DeferErr(f.Close())

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file bytes: %w", err)
	}

	return bytes, nil
}

func (a *arcade) LoadBytes(start, expectedSize uint16, bytes []uint8) error {
	if len(bytes) != int(expectedSize) {
		return fmt.Errorf("unexpected size difference when loading bytes, expected: %d, actual: %d", expectedSize, len(bytes))
	}

	addr := start
	for _, b := range bytes {
		a.memory.Write(addr, b)
		addr++
	}

	return nil
}

func Disassemble(romPath string) error {
	if len(romPath) == 0 {
		return errors.New("no rom passed to emulator")
	}

	romBytes := []byte{}

	if filepath.Ext(romPath) == ".zip" {
		r, err := zip.OpenReader(romPath)
		if err != nil {
			return fmt.Errorf("failed to open zip archive: %w", err)
		}
		defer lib.DeferErr(r.Close())

		settings, err := gamespec.GetGameSettings(filepath.Base(romPath))
		if err != nil {
			return err
		}

		for _, p := range settings.ROMParts {
			romBytes = append(romBytes, lib.Must(GetFileBytesFromZip(r.File, p.FileName))...)
		}
	} else {
		var err error

		romBytes, err = os.ReadFile(romPath)
		if err != nil {
			return fmt.Errorf("failed to read rom file: %w", err)
		}
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

type saveState struct {
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

	s := saveState{
		CPU:    cpu,
		Memory: memory,
	}

	romDir, romFileName := filepath.Split(a.romPath)
	stateFilePath := filepath.Join(romDir, strings.ReplaceAll(romFileName, filepath.Ext(romFileName), ".state"))

	f, err := os.Create(stateFilePath)
	if err != nil {
		return err
	}
	defer lib.DeferErr(f.Close())

	enc := gob.NewEncoder(f)

	if err := enc.Encode(s); err != nil {
		return err
	}

	fmt.Println("saved state file: " + stateFilePath)

	return nil
}

func (a *arcade) LoadState() error {
	stateFilePath := a.saveState
	if stateFilePath == "" {
		romDir, romFileName := filepath.Split(a.romPath)
		stateFilePath = filepath.Join(romDir, strings.ReplaceAll(romFileName, filepath.Ext(romFileName), ".state"))
	}

	f, err := os.Open(stateFilePath)
	if err != nil {
		return err
	}
	defer lib.DeferErr(f.Close())

	var s saveState

	dec := gob.NewDecoder(f)
	if err := dec.Decode(&s); err != nil {
		return fmt.Errorf("failed to decode save state: %w", err)
	}

	if err := a.cpu.LoadState(s.CPU); err != nil {
		return fmt.Errorf("failed to load CPU state: %w", err)
	}

	for addr, b := range s.Memory {
		a.memory.Write(uint16(addr), b)
	}

	fmt.Println("loaded state file: " + stateFilePath)

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
