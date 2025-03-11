package main

import (
	"FrenchConnections/internal"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:        "server-fc",
		Usage:       "French-Connection Backend Server",
		Description: "Server that runs the logic for the french connection adaption",
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("Hello friend!")

			if internal.DEBUG {
				slog.SetLogLoggerLevel(slog.LevelDebug)
				slog.Debug(fmt.Sprintf("DEBUG -> %t\n", internal.DEBUG))
			}

			slog.Debug(fmt.Sprintf("API_DOMAIN -> %s\n", internal.API_DOMAIN))
			slog.Debug(fmt.Sprintf("API_IP -> %s\n", internal.API_IP))
			slog.Debug(fmt.Sprintf("API_PORT -> %v\n", internal.API_PORT))

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "api-ip",
				Aliases:     []string{"ip"},
				Usage:       "The IP to expose the API against",
				Value:       "127.0.0.1",
				Destination: &internal.API_IP,
			},
			&cli.StringFlag{
				Name:        "api-domain",
				Usage:       "The domain on which the API is exposed",
				Value:       "localhost",
				Destination: &internal.API_DOMAIN,
			},
			&cli.IntFlag{
				Name:        "api-port",
				Aliases:     []string{"p"},
				Usage:       "The port to expose the API against",
				Value:       8080,
				Destination: &internal.API_PORT,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Value:       false,
				Destination: &internal.DEBUG,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
