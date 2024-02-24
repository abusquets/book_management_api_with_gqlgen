package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	cliApp := &cli.App{
		Before: func(cCtx *cli.Context) error {
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "test",
				Aliases: []string{"c"},
				Usage:   "test command",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("Cli is working fine!!!")
					return nil
				},
			},
		},
	}

	if cliApp.Metadata == nil {
		cliApp.Metadata = make(map[string]interface{})
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
