package ui

import (
	"fmt"
	"unsafe"

	"github.com/Zyko0/go-sdl3/sdl"
)

type bus interface {
	Read(addr uint16) uint8
}

type cpu interface {
	RequestInterrupt(id uint8)
	SendInput(port uint8, bit uint8, value bool)
}

type apu interface {
	TogglePauseAudio(paused bool)
}

type arcade interface {
	Reset()
	SaveState() error
	LoadState() error
	Shutdown()
}

type UI struct {
	Arcade arcade
	Bus    bus
	CPU    cpu
	APU    apu
	Paused bool

	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
	surface  *sdl.Surface
}

const (
	VRAM_START uint16 = 0x2400
	WIDTH      uint16 = 224
	HEIGHT     uint16 = 256
	SCALE      uint16 = 3

	COLOR_BLACK uint32 = 0xFF000000
	COLOR_WHITE uint32 = 0xFFFFFFFF
	COLOR_RED   uint32 = 0xFFFF0000
	COLOR_GREEN uint32 = 0xFF00FF00
)

func (ui *UI) Init() {
	ui.Paused = false

	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		panic("failed to init sdl: " + err.Error())
	}

	if ui.window == nil && ui.renderer == nil {
		ui.window, ui.renderer, err = sdl.CreateWindowAndRenderer("Space Invaders", int(WIDTH*SCALE), int(HEIGHT*SCALE), sdl.WINDOW_RESIZABLE)
		if err != nil {
			panic("failed to create window and renderer: " + err.Error())
		}
	}

	if ui.texture == nil {
		ui.texture, err = ui.renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, int(WIDTH), int(HEIGHT))
		if err != nil {
			panic("failed to create texture: " + err.Error())
		}

		ui.texture.SetScaleMode(sdl.SCALEMODE_NEAREST)
	}

	if ui.surface == nil {
		ui.surface, err = sdl.CreateSurface(int(WIDTH), int(HEIGHT), sdl.PIXELFORMAT_ARGB8888)
		if err != nil {
			panic("failed to create surface: " + err.Error())
		}
	}
}

func (ui *UI) Close() {
	ui.surface.Destroy()
	ui.texture.Destroy()
	ui.renderer.Destroy()
	ui.window.Destroy()
}

func (ui *UI) Step() {
	ui.drawVRAM()
	ui.handleEvents()
}

func (ui *UI) drawVRAM() {
	pixels := ui.surface.Pixels()
	pitch := int(ui.surface.Pitch) / 4
	pixelData := unsafe.Slice((*uint32)(unsafe.Pointer(&pixels[0])), len(pixels)/4)

	for y := range HEIGHT {
		rowStart := int(y) * pitch
		for x := range WIDTH {
			addr := VRAM_START + (x * (HEIGHT / 8)) + ((HEIGHT - y - 1) / 8)
			pixels := ui.Bus.Read(addr)
			pixel := (pixels >> (7 - y%8)) & 1

			color := COLOR_BLACK
			if pixel == 1 {
				color = COLOR_WHITE

				if y > 180 {
					if y < 240 || (x > 16 && x < 128) {
						color = COLOR_GREEN
					}
				} else if y >= 32 && y < 64 {
					color = COLOR_RED
				}
			}

			pixelData[rowStart+int(x)] = color
		}

		if y == HEIGHT/2 && !ui.Paused {
			ui.CPU.RequestInterrupt(1)
		}

		if y == HEIGHT-1 && !ui.Paused {
			ui.CPU.RequestInterrupt(2)
		}
	}

	if err := ui.texture.Update(nil, ui.surface.Pixels(), ui.surface.Pitch); err != nil {
		panic("failed to update texture: " + err.Error())
	}

	if err := ui.renderer.Clear(); err != nil {
		panic("failed to clear renderer: " + err.Error())
	}

	if err := ui.renderer.RenderTexture(ui.texture, nil, nil); err != nil {
		panic("failed to render texture: " + err.Error())
	}

	if err := ui.renderer.Present(); err != nil {
		panic("failed to present UI: " + err.Error())
	}
}

func (ui *UI) handleEvents() {
	var event sdl.Event

	for sdl.PollEvent(&event) {
		switch event.Type {
		case sdl.EVENT_QUIT, sdl.EVENT_WINDOW_DESTROYED:
			ui.Arcade.Shutdown()

		case sdl.EVENT_KEY_DOWN, sdl.EVENT_KEY_UP:
			pressed := event.Type == sdl.EVENT_KEY_DOWN

			switch event.KeyboardEvent().Key {
			case sdl.K_R: // Reset
				if !pressed {
					ui.Arcade.Reset()
				}
			case sdl.K_P: // Pause
				if !pressed {
					ui.Paused = !ui.Paused
					ui.APU.TogglePauseAudio(ui.Paused)
				}
			case sdl.K_9:
				if !pressed {
					err := ui.Arcade.LoadState()
					if err != nil {
						fmt.Println("failed to save state:", err.Error())
					}
				}
			case sdl.K_0:
				if !pressed {
					err := ui.Arcade.SaveState()
					if err != nil {
						fmt.Println("failed to save state:", err.Error())
					}
				}

			// Menu
			case sdl.K_C: // Add coin
				ui.CPU.SendInput(1, 0, pressed)
			case sdl.K_1: // Select 1 player
				ui.CPU.SendInput(1, 2, pressed)
			case sdl.K_2: // Select 2 players
				ui.CPU.SendInput(1, 1, pressed)

			// P1 controls
			case sdl.K_A: // Left
				ui.CPU.SendInput(1, 5, pressed)
			case sdl.K_D: // Right
				ui.CPU.SendInput(1, 6, pressed)
			case sdl.K_W: // Shoot
				ui.CPU.SendInput(1, 4, pressed)

			// P2 controls
			case sdl.K_LEFT: // Left
				ui.CPU.SendInput(2, 5, pressed)
			case sdl.K_RIGHT: // Right
				ui.CPU.SendInput(2, 6, pressed)
			case sdl.K_UP: // Shoot
				ui.CPU.SendInput(2, 4, pressed)
			}
		}
	}
}
