package main

import (
	"os"

	"git.burning.moe/celediel/agedit/internal/config"

	"filippo.io/age"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
)

var (
	identity   *age.X25519Identity
	logger     *log.Logger
	cfg        config.Config
	configFile string

	input_file, output_file string
	force_overwrite         bool
)

func main() {
	app := &cli.App{
		Name:                  name,
		Usage:                 usage,
		Version:               version,
		Authors:               authors,
		Flags:                 flags,
		Before:                before,
		Action:                action,
		CustomAppHelpTemplate: help_template,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
