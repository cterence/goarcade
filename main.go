package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/cterence/goarcade/internal/arcade"
	"github.com/urfave/cli/v3"
)

func main() {
	var (
		debug      bool
		cpm        bool
		headless   bool
		unthrottle bool
		noAudio    bool
		soundDir   string
		saveState  string
	)

	cmd := &cli.Command{
		Name:      "goarcade",
		Usage:     "Intel 8080 arcade emulator",
		ArgsUsage: "[rom path (binary file or .zip archive)]",
		Flags: []cli.Flag{
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

			&cli.StringFlag{
				Name:        "state",
				Usage:       "save state file",
				Destination: &saveState,
			},

			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "print debug logs",
				Destination: &debug,
			},

			&cli.BoolFlag{
				Name:        "headless",
				Usage:       "run without UI window",
				Destination: &headless,
			},

			&cli.StringFlag{
				Name:        "sound-dir",
				Aliases:     []string{"s"},
				Usage:       "directory path for WAV sound files",
				Destination: &soundDir,
			},

			&cli.BoolFlag{
				Name:        "no-audio",
				Usage:       "run without audio",
				Destination: &noAudio,
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
			return arcade.Run(
				ctx,
				cmd.Args().First(),
				arcade.WithDebug(debug),
				arcade.WithCPM(cpm),
				arcade.WithHeadless(headless),
				arcade.WithNoAudio(noAudio),
				arcade.WithSoundDir(soundDir),
				arcade.WithUnthrottle(unthrottle),
				arcade.WithSaveState(saveState),
			)
		},
		Commands: []*cli.Command{
			{
				Name:      "dasm",
				Aliases:   []string{"d"},
				Usage:     "disassemble a program",
				ArgsUsage: "[rom path (binary file or .zip archive)]",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return arcade.Disassemble(cmd.Args().First())
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("runtime error: %v", err)
	}
}
