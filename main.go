package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cterence/goarcade/internal/arcade"
	"github.com/urfave/cli/v3"
)

func readFiles(romPath, configPath, soundDir string) ([]uint8, []uint8, [][]uint8, error) {
	romBytes, err := os.ReadFile(romPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read rom file: %w", err)
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var soundListBytes [][]uint8

	if soundDir != "" {
		soundFiles, err := os.ReadDir(soundDir)
		if err != nil {
			panic("failed to read sound directory: " + err.Error())
		}

		soundListBytes = make([][]uint8, len(soundFiles))

		for i, f := range soundFiles {
			if filepath.Ext(f.Name()) == ".wav" {
				soundData, err := os.ReadFile(filepath.Join(soundDir, f.Name()))
				if err != nil {
					panic("failed to load WAV file: " + err.Error())
				}

				soundListBytes[i] = make([]uint8, len(soundData))
				copy(soundListBytes[i], soundData)
			}
		}
	}

	return romBytes, configBytes, soundListBytes, err
}

func main() {
	var (
		debug         bool
		cpm           bool
		headless      bool
		unthrottle    bool
		mute          bool
		soundDir      string
		saveStatePath string
		configPath    string
	)

	cmd := &cli.Command{
		Name:      "goarcade",
		Usage:     "Intel 8080 arcade emulator",
		ArgsUsage: "[rom path (binary file or .zip archive)]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "config file path",
				Value:       "./config.yaml",
				TakesFile:   true,
				Destination: &configPath,
			},

			&cli.StringFlag{
				Name:        "state",
				Aliases:     []string{"s"},
				Usage:       "save state file path",
				TakesFile:   true,
				Destination: &saveStatePath,
			},

			&cli.StringFlag{
				Name:        "sound-dir",
				Aliases:     []string{"sd"},
				Usage:       "directory path for WAV sound files",
				TakesFile:   true,
				Destination: &soundDir,
			},

			&cli.BoolFlag{
				Name:    "pprof",
				Aliases: []string{"p"},
				Usage:   "run pprof webserver on localhost:6060",
				Action: func(_ context.Context, _ *cli.Command, _ bool) error {
					go func() {
						log.Println(http.ListenAndServe("localhost:6060", nil))
					}()

					return nil
				},
			},

			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "print debug logs",
				Destination: &debug,
			},

			&cli.BoolFlag{
				Name:        "headless",
				Aliases:     []string{"hl"},
				Usage:       "run without UI window",
				Destination: &headless,
			},
			&cli.BoolFlag{
				Name:        "mute",
				Aliases:     []string{"m"},
				Usage:       "run without audio",
				Destination: &mute,
			},

			&cli.BoolFlag{
				Name:        "cpm",
				Usage:       "run in CP/M compatibility mode (for CPU tests)",
				Destination: &cpm,
			},

			&cli.BoolFlag{
				Name:        "unthrottle",
				Aliases:     []string{"u"},
				Usage:       "do not throttle cpu at 2MHz",
				Destination: &unthrottle,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			romPath := cmd.Args().First()

			if romPath == "" {
				fmt.Printf("error: no rom path given\n\n")
				return cli.ShowSubcommandHelp(cmd)
			}

			romBytes, configBytes, soundListBytes, err := readFiles(romPath, configPath, soundDir)
			if err != nil {
				return err
			}

			return arcade.Run(
				ctx,
				romBytes,
				configBytes,
				soundListBytes,
				romPath,
				arcade.WithDebug(debug),
				arcade.WithCPM(cpm),
				arcade.WithHeadless(headless),
				arcade.WithMute(mute),
				arcade.WithUnthrottle(unthrottle),
				arcade.WithSaveState(saveStatePath),
			)
		},
		Commands: []*cli.Command{
			{
				Name:      "dasm",
				Aliases:   []string{"d"},
				Usage:     "disassemble a program",
				ArgsUsage: "[rom path (binary file or .zip archive)]",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					romPath := cmd.Args().First()

					if romPath == "" {
						fmt.Printf("error: no rom path given\n\n")
						return cli.ShowSubcommandHelp(cmd)
					}

					romBytes, configBytes, _, err := readFiles(romPath, configPath, "")
					if err != nil {
						return err
					}

					return arcade.Disassemble(romBytes, configBytes, romPath)
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("runtime error: %v", err)
	}
}
