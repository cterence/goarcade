package ui

import (
	"context"

	"github.com/Zyko0/go-sdl3/sdl"
)

type UI struct {
	ReadMem func(uint16) uint8

	framebuffer [WIDTH][HEIGHT]uint8
	cancel      context.CancelFunc

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
)

var palette = [4]uint32{0xFF000000, 0xFFFFFFFF, 0xFF00FF00, 0xFFFF0000}

func (ui *UI) Init(cancel context.CancelFunc) {
	ui.framebuffer = [WIDTH][HEIGHT]uint8{}
	ui.cancel = cancel

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
		ui.texture, err = ui.renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, int(WIDTH*SCALE), int(HEIGHT*SCALE))
		if err != nil {
			panic("failed to create texture: " + err.Error())
		}
	}

	if ui.surface == nil {
		ui.surface, err = sdl.CreateSurface(int(WIDTH*SCALE), int(HEIGHT*SCALE), sdl.PIXELFORMAT_ARGB8888)
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
	ui.updateFramebuffer()
	ui.drawFramebuffer()
	ui.handleEvents()
}

func (ui *UI) updateFramebuffer() {
	for y := range HEIGHT {
		for x := range WIDTH {
			addr := VRAM_START + (x * (HEIGHT / 8)) + ((HEIGHT - y - 1) / 8)
			pixels := ui.ReadMem(addr)
			pixel := (pixels >> (7 - y%8)) & 1
			ui.framebuffer[x][y] = pixel
		}
	}
}

func (ui *UI) drawFramebuffer() {
	for y := range HEIGHT {
		for x := range WIDTH {
			rc := &sdl.Rect{
				X: int32(x * SCALE),
				Y: int32(y * SCALE),
				W: int32(SCALE),
				H: int32(SCALE),
			}

			color := palette[ui.framebuffer[x][y]]

			if err := ui.surface.FillRect(rc, color); err != nil {
				panic("failed to fill rect: " + err.Error())
			}
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
			ui.cancel()
		}
	}
}
