package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/cterence/space-invaders/internal/arcade"
	"github.com/urfave/cli/v3"
)

func main() {
	var (
		debug bool
	)

	cmd := &cli.Command{
		Name:      "space-invaders",
		Usage:     "Space Invaders arcade emulator",
		ArgsUsage: "[ordered rom part paths]",
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

			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "print debug logs",
				Destination: &debug,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return arcade.Run(
				ctx,
				cmd.Args().Slice(),
				arcade.WithDebug(debug),
			)
		},
		Commands: []*cli.Command{
			{
				Name:      "disassemble",
				Usage:     "disassemble a 8080 rom",
				ArgsUsage: "[ordered rom part paths]",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return arcade.Disassemble(cmd.Args().Slice())
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("runtime error: %v", err)
	}
}
