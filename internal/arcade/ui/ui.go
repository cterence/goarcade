package ui

import (
	"encoding/binary"
	"fmt"

	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/cterence/goarcade/internal/arcade/config"
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

	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
	surface  *sdl.Surface

	ColorOverlays []config.ColorOverlay
	ColorPROM     []uint8

	colors      [WIDTH][HEIGHT]uint32
	framebuffer [WIDTH * HEIGHT * PIXEL_BYTES]uint8

	Paused     bool
	startYDraw int
}

const (
	VRAM_START uint16 = 0x2400
	VRAM_SIZE  uint16 = 0x1C00

	WIDTH       = 224
	HEIGHT      = 256
	PIXEL_BYTES = 4
	SCALE       = 3

	COLOR_BLACK uint32 = 0xFF000000
	COLOR_WHITE uint32 = 0xFFFFFFFF
)

func (ui *UI) Init() {
	ui.Paused = false
	ui.startYDraw = 0
	ui.computeColorLUT()

	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		panic("failed to init sdl: " + err.Error())
	}

	if ui.window == nil && ui.renderer == nil {
		ui.window, ui.renderer, err = sdl.CreateWindowAndRenderer("goarcade", WIDTH*SCALE, HEIGHT*SCALE, sdl.WINDOW_RESIZABLE)
		if err != nil {
			panic("failed to create window and renderer: " + err.Error())
		}
	}

	if ui.texture == nil {
		ui.texture, err = ui.renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, WIDTH, HEIGHT)
		if err != nil {
			panic("failed to create texture: " + err.Error())
		}

		if err := ui.texture.SetScaleMode(sdl.SCALEMODE_NEAREST); err != nil {
			panic("failed to set texture scale mode: " + err.Error())
		}
	}

	if ui.surface == nil {
		ui.surface, err = sdl.CreateSurface(WIDTH, HEIGHT, sdl.PIXELFORMAT_ARGB8888)
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
	y := ui.startYDraw

	// Draw half the frame so that the UI requests half and full VBLANK interrupts at the correct CPU timing
	for y < ui.startYDraw+HEIGHT/2 {
		row := ui.framebuffer[y*int(ui.surface.Pitch) : y*int(ui.surface.Pitch)+WIDTH*PIXEL_BYTES]
		for x := range WIDTH {
			addr := VRAM_START + uint16(x*(HEIGHT/8)) + uint16((HEIGHT-y-1)/8)
			vramPixels := ui.Bus.Read(addr)
			vramPixel := (vramPixels >> (7 - y%8)) & 1
			color := COLOR_BLACK

			if vramPixel == 1 {
				color = ui.getColor(x, y)
			}

			binary.LittleEndian.PutUint32(row[x*PIXEL_BYTES:], color)
		}

		if y == (HEIGHT/2)-1 && !ui.Paused {
			ui.CPU.RequestInterrupt(1)
		}

		if y == HEIGHT-1 && !ui.Paused {
			ui.CPU.RequestInterrupt(2)
		}

		y++
	}

	ui.startYDraw = y % HEIGHT

	if err := ui.texture.Update(nil, ui.framebuffer[:], ui.surface.Pitch); err != nil {
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

					if ui.Paused {
						fmt.Println("arcade paused")
					} else {
						fmt.Println("arcade resumed")
					}
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
			case sdl.K_5: // Add coin
				ui.CPU.SendInput(1, 0, pressed)
			case sdl.K_1: // Select 1 player
				ui.CPU.SendInput(1, 2, pressed)
			case sdl.K_2: // Select 2 players
				ui.CPU.SendInput(1, 1, pressed)

			// Player controls
			case sdl.K_LEFT: // Left
				ui.CPU.SendInput(1, 5, pressed)
				ui.CPU.SendInput(2, 5, pressed)
			case sdl.K_RIGHT: // Right
				ui.CPU.SendInput(1, 6, pressed)
				ui.CPU.SendInput(2, 6, pressed)
			case sdl.K_LCTRL: // Shoot
				ui.CPU.SendInput(1, 4, pressed)
				ui.CPU.SendInput(2, 4, pressed)
			}
		}
	}
}

func (ui *UI) computeColorLUT() {
	for x := range WIDTH {
		for y := range HEIGHT {
			color := COLOR_WHITE

			for _, cm := range ui.ColorOverlays {
				xMatch := (cm.XMin == 0 && cm.XMax == 0) || (x >= int(cm.XMin) && x <= int(cm.XMax))

				yMatch := (cm.YMin == 0 && cm.YMax == 0) || (y >= int(cm.YMin) && y <= int(cm.YMax))
				if xMatch && yMatch {
					color = cm.Color

					break
				}
			}

			ui.colors[x][y] = color
		}
	}
}

func (ui *UI) getColor(x, y int) uint32 {
	if len(ui.ColorPROM) == 0 {
		return ui.colors[x][y]
	}

	// Convert horizontal coords to vertical (flip 90 degree clockwise)
	origX := HEIGHT - y - 1
	origY := x
	offs := origY*32 + (origX >> 3)
	colorAddress := uint16((((offs >> 8) << 5) | (offs & 0x1F)) + 0x80)
	colorBits := ui.ColorPROM[colorAddress] & 0x07

	// RED BLUE GREEN
	var r, b, g uint32
	if colorBits&0x01 != 0 {
		r = 0xFF
	}

	if colorBits&0x02 != 0 {
		b = 0xFF
	}

	if colorBits&0x04 != 0 {
		g = 0xFF
	}

	return 0xFF000000 | (r << 16) | (g << 8) | b
}
