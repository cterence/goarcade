package apu

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Zyko0/go-sdl3/sdl"
)

type APU struct {
	device   sdl.AudioDeviceID
	streams  []*sdl.AudioStream
	sounds   [][]uint8
	looping  []bool
	loopStop []chan struct{}
}

func (a *APU) Init(soundDir string) {
	if soundDir == "" {
		fmt.Println("warning: sound files not loaded, audio disabled")

		return
	}

	err := sdl.Init(sdl.INIT_AUDIO)
	if err != nil {
		panic("failed to init sdl audio: " + err.Error())
	}

	spec := &sdl.AudioSpec{
		Format:   sdl.AUDIO_U8,
		Channels: 1,
		Freq:     11025,
	}

	a.device, err = sdl.AUDIO_DEVICE_DEFAULT_PLAYBACK.OpenAudioDevice(spec)
	if err != nil {
		panic("failed to get default playback audio device: " + err.Error())
	}

	wavFiles, err := os.ReadDir(soundDir)
	if err != nil {
		panic("failed to read sound directory: " + err.Error())
	}

	a.streams = make([]*sdl.AudioStream, len(wavFiles))
	a.sounds = make([][]uint8, len(wavFiles))
	a.looping = make([]bool, len(wavFiles))
	a.loopStop = make([]chan struct{}, len(wavFiles))

	for i, f := range wavFiles {
		if filepath.Ext(f.Name()) == ".wav" {
			soundData, err := sdl.LoadWAV(filepath.Join(soundDir, f.Name()), spec)
			if err != nil {
				panic("failed to load WAV file: " + err.Error())
			}

			// Downscale volume at 33% to allow 3 sounds to play simultaneously without audio clipping
			a.sounds[i] = scaleVolume(soundData, 0.33)

			if a.streams[i] == nil {
				a.streams[i], err = sdl.CreateAudioStream(spec, spec)
				if err != nil {
					panic("failed to create audio stream: " + err.Error())
				}
			}

			if a.streams[i].Device() == 0 {
				if err := a.device.BindAudioStream(a.streams[i]); err != nil {
					panic("failed to bind audio stream to device: " + err.Error())
				}
			}
		}
	}

	a.TogglePauseAudio(false)
}

func (a *APU) TogglePauseAudio(pause bool) {
	if a.streams == nil {
		return
	}

	if pause {
		if err := a.device.Pause(); err != nil {
			panic("failed to pause default audio device: " + err.Error())
		}
	} else {
		if err := a.device.Resume(); err != nil {
			panic("failed to resume default audio device: " + err.Error())
		}
	}
}

func (a *APU) Close() {
	for _, s := range a.streams {
		s.Destroy()
	}
}

func (a *APU) PlaySound(soundIndex uint8) {
	if len(a.streams) == 0 {
		return
	}

	available, err := a.streams[soundIndex].Available()
	if err != nil {
		panic("failed to get available audio stream: " + err.Error())
	}

	sound := a.sounds[soundIndex]

	if available < int32(len(sound)) {
		err := a.streams[soundIndex].PutData(sound)
		if err != nil {
			panic("failed to put sound data to stream: " + err.Error())
		}
	}
}

func (a *APU) StartSoundLoop(soundIndex uint8) {
	if len(a.streams) == 0 || a.looping[soundIndex] {
		return
	}

	a.looping[soundIndex] = true
	a.loopStop[soundIndex] = make(chan struct{})

	go func() {
		stream := a.streams[soundIndex]
		sound := a.sounds[soundIndex]

		for {
			select {
			case <-a.loopStop[soundIndex]:
				stream.Clear()

				return
			default:
				queued, err := stream.Queued()
				if err != nil {
					panic("failed to get queued bytes: " + err.Error())
				}

				if queued < int32(len(sound)) {
					stream.PutData(sound)
				}

				time.Sleep(10 * time.Millisecond)
			}
		}
	}()
}

func (a *APU) StopSoundLoop(soundIndex uint8) {
	if len(a.streams) == 0 || !a.looping[soundIndex] {
		return
	}

	close(a.loopStop[soundIndex])
	a.looping[soundIndex] = false
}

func scaleVolume(data []byte, scale float64) []byte {
	scaled := make([]byte, len(data))
	for i, sample := range data {
		// Convert U8 (0-255, center at 128) to signed, scale, convert back
		signed := float64(sample) - 128
		signed *= scale
		scaled[i] = byte(signed + 128)
	}

	return scaled
}
