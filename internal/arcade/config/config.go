package config

import (
	"errors"
	"fmt"

	"github.com/goccy/go-yaml"
)

type ROMPart struct {
	FileName     string `yaml:"fileName"`
	StartAddr    uint16 `yaml:"startAddr"`
	ExpectedSize uint16 `yaml:"expectedSize"`
}

type ColorPROM struct {
	FileName     string `yaml:"fileName"`
	ExpectedSize uint16 `yaml:"expectedSize"`
}

type ColorOverlay struct {
	XMin  uint16 `yaml:"xMin"`
	XMax  uint16 `yaml:"xMax"`
	YMin  uint16 `yaml:"yMin"`
	YMax  uint16 `yaml:"yMax"`
	Color uint32 `yaml:"color"`
}

type activePosition string

const (
	POSITION_LOW  activePosition = "low"
	POSITION_HIGH activePosition = "high"
)

type Port struct {
	Bit    uint8 `yaml:"bit"`
	Active bool  `yaml:"active"`
	// ActivePosition activePosition `yaml:"activePosition"`
}

type GameSpec struct {
	InPorts       map[int][8]Port `yaml:"inPorts"`
	ROMParts      []ROMPart       `yaml:"romParts"`
	ColorOverlays []ColorOverlay  `yaml:"colorOverlays"`
	ColorPROMs    []ColorPROM     `yaml:"colorPROMs"`
}

type Config struct {
	GameSpecs map[string]GameSpec `yaml:"gameSpecs"`
}

const (
	MAX_X uint16 = 224
	MAX_Y uint16 = 256
)

func LoadConfig(configBytes []uint8, gameName string) (*GameSpec, error) {
	var config Config

	err := yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gamespecs: %w", err)
	}

	s, ok := config.GameSpecs[gameName]
	if !ok {
		return nil, fmt.Errorf("no specs for game: %s", gameName)
	}

	if err := validateConfig(&s); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &s, nil
}

func validateConfig(s *GameSpec) error {
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
