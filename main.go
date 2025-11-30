package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cterence/space-invaders/internal/arcade"
	"github.com/cterence/space-invaders/internal/arcade/cpu"
	"github.com/urfave/cli/v3"
)

func main() {
	var (
		debug bool
	)

	cmd := &cli.Command{
		Name:      "space-invaders",
		Usage:     "Space Invaders arcade emulator",
		ArgsUsage: "[rom directory path]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Destination: &debug,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			return arcade.Run(
				ctx,
				c.Args().First(),
				arcade.WithDebug(debug),
			)
		},
		Commands: []*cli.Command{
			{
				Name:      "disassemble",
				Usage:     "disassemble a 8080 rom",
				ArgsUsage: "[rom directory path]",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					romDirPath := cmd.Args().First()

					if romDirPath == "" {
						return errors.New("empty rom directory path")
					}

					romBytes := []byte{}

					for _, id := range []string{"h", "g", "f", "e"} {
						romPartBytes, err := os.ReadFile(filepath.Join(romDirPath, "invaders."+id))
						if err != nil {
							return fmt.Errorf("failed to read rom file: %w", err)
						}
						romBytes = append(romBytes, romPartBytes...)
					}

					var i uint16

					c := cpu.CPU{}
					c.Init()
					c.ReadMem = func(offset uint16) uint8 { return romBytes[i+offset] }

					for i < uint16(len(romBytes)) {
						inst, op1, op2, instLength := c.DecodeInst()
						var b strings.Builder

						b.WriteString(fmt.Sprintf("%04x: ", i))
						if instLength == 1 {
							b.WriteString(fmt.Sprintf("%02x          ", romBytes[i]))
						}
						if instLength == 2 {
							b.WriteString(fmt.Sprintf("%02x %02x       ", romBytes[i], romBytes[i+1]))
						}
						if instLength == 3 {
							b.WriteString(fmt.Sprintf("%02x %02x %02x    ", romBytes[i], romBytes[i+1], romBytes[i+2]))
						}
						b.WriteString(inst)
						if op1 != "" {
							b.WriteString(" " + op1)
							if op2 != "" {
								b.WriteString(", " + op2)
							}
						}
						fmt.Println(b.String())
						i += instLength
					}

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("runtime error: %v", err)
	}
}
