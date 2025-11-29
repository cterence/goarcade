package main

import (
	"context"
	"log"
	"os"

	"github.com/cterence/space-invaders/internal/arcade"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:      "space-invaders",
		Usage:     "Space Invaders arcade emulator",
		ArgsUsage: "[rom directory path]",
		Action: func(ctx context.Context, c *cli.Command) error {
			return arcade.Run(ctx, c.Args().First())
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("runtime error: %v", err)
	}
}
