package compatibility

import "fmt"

type GamePart struct {
	FileName     string
	StartAddr    uint16
	ExpectedSize uint16
}

type Settings struct {
	GameParts []GamePart
}

func GetGameSettings(zipFileName string) (*Settings, error) {
	s := &Settings{}

	switch zipFileName {
	case "invaders.zip":
		s.GameParts = []GamePart{
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
		}
	default:
		return nil, fmt.Errorf("no compatibility settings for game: %s", zipFileName)
	}

	return s, nil
}
