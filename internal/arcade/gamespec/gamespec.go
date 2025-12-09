package gamespec

import (
	"errors"
	"fmt"
)

type ROMPart struct {
	FileName     string
	StartAddr    uint16
	ExpectedSize uint16
}

type ColorPROM struct {
	FileName     string
	ExpectedSize uint16
}

type ColorOverlay struct {
	XMin  uint16
	XMax  uint16
	YMin  uint16
	YMax  uint16
	Color uint32
}

type Specs struct {
	ROMParts      []ROMPart
	ColorOverlays []ColorOverlay
	ColorPROMs    []ColorPROM
}

const (
	MAX_X uint16 = 224
	MAX_Y uint16 = 256

	COLOR_BLACK   uint32 = 0xFF000000
	COLOR_WHITE   uint32 = 0xFFFFFFFF
	COLOR_RED     uint32 = 0xFFFF0000
	COLOR_GREEN   uint32 = 0xFF00FF00
	COLOR_MAGENTA uint32 = 0xFFFF00FF
	COLOR_YELLOW  uint32 = 0xFFFFFF00
	COLOR_CYAN    uint32 = 0xFF00FFFF
)

var gameSpecs = map[string]Specs{
	"invaders.zip": {
		ROMParts: []ROMPart{
			{
				FileName:     "invaders.h",
				StartAddr:    0x0,
				ExpectedSize: 0x800,
			},
			{
				FileName:     "invaders.g",
				StartAddr:    0x800,
				ExpectedSize: 0x800,
			},
			{
				FileName:     "invaders.f",
				StartAddr:    0x1000,
				ExpectedSize: 0x800,
			},
			{
				FileName:     "invaders.e",
				StartAddr:    0x1800,
				ExpectedSize: 0x800,
			},
		},
		ColorOverlays: []ColorOverlay{
			{
				YMin:  32,
				YMax:  63,
				Color: COLOR_RED,
			},
			{
				YMin:  180,
				YMax:  240,
				Color: COLOR_GREEN,
			},
			{
				XMin:  16,
				XMax:  128,
				YMin:  241,
				YMax:  MAX_Y,
				Color: COLOR_GREEN,
			},
		},
	},

	"invadpt2.zip": {
		ROMParts: []ROMPart{
			{
				FileName:     "pv01",
				StartAddr:    0x0,
				ExpectedSize: 0x800,
			},
			{
				FileName:     "pv02",
				StartAddr:    0x800,
				ExpectedSize: 0x800,
			},
			{
				FileName:     "pv03",
				StartAddr:    0x1000,
				ExpectedSize: 0x800,
			},
			{
				FileName:     "pv04",
				StartAddr:    0x1800,
				ExpectedSize: 0x800,
			},
			{
				FileName:     "pv05",
				StartAddr:    0x4000,
				ExpectedSize: 0x800,
			},
		},
		ColorPROMs: []ColorPROM{
			{
				FileName:     "pv06.1",
				ExpectedSize: 0x400,
			},
			{
				FileName:     "pv07.2",
				ExpectedSize: 0x400,
			},
		},
	},
}

func GetGameSettings(zipFileName string) (*Specs, error) {
	s, ok := gameSpecs[zipFileName]
	if !ok {
		return nil, fmt.Errorf("no compatibility settings for game: %s", zipFileName)
	}

	if err := validateSettings(&s); err != nil {
		return nil, fmt.Errorf("settings validation failed: %w", err)
	}

	return &s, nil
}

func validateSettings(s *Specs) error {
	if len(s.ROMParts) == 0 {
		return errors.New("missing game parts")
	}

	prevPart := s.ROMParts[0]
	for i, currentPart := range s.ROMParts {
		if i == 0 {
			continue
		}

		if currentPart.StartAddr <= prevPart.StartAddr {
			return fmt.Errorf("game parts: start address %x of part %s is higher than start address %x of part %s", prevPart.StartAddr, prevPart.FileName, currentPart.StartAddr, currentPart.FileName)
		}

		if prevPart.StartAddr+prevPart.ExpectedSize > currentPart.StartAddr {
			return fmt.Errorf("game parts: part %s (start: %x, end: %x) overlaps with part %s (start: %x, end: %x)", prevPart.FileName, prevPart.StartAddr, prevPart.StartAddr+prevPart.ExpectedSize, currentPart.FileName, currentPart.StartAddr, currentPart.StartAddr+currentPart.ExpectedSize)
		}

		prevPart = currentPart
	}

	if len(s.ColorOverlays) >= 2 {
		for x := range MAX_X {
			for y := range MAX_Y {
				var matchingCM *int

				for i, cm := range s.ColorOverlays {
					xMatch := (cm.XMin == 0 && cm.XMax == 0) || (x >= cm.XMin && x <= cm.XMax)

					yMatch := (cm.YMin == 0 && cm.YMax == 0) || (y >= cm.YMin && y <= cm.YMax)
					if xMatch && yMatch {
						if matchingCM != nil {
							return fmt.Errorf("color overlays: overlays %d and %d are overlapping at pixel x: %d, y: %d", *matchingCM, i, x, y)
						}

						matchingCM = &i
					}
				}
			}
		}
	}

	return nil
}
