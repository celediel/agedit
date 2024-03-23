package main

import (
	"os"

	"git.burning.moe/celediel/agedit/internal/config"
	"git.burning.moe/celediel/agedit/pkg/editor"

	"filippo.io/age"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
)

var (
	identities              []age.Identity
	recipients              []age.Recipient
	logger                  *log.Logger
	cfg                     config.Config
	edt                     editor.Editor
	configFile              string
	input_file, output_file string
	force_overwrite         bool
)

func main() {
	app := &cli.App{
		Name:                   name,
		Usage:                  usage,
		Version:                version,
		Authors:                authors,
		Flags:                  flags,
		Before:                 before,
		Action:                 action,
		CustomAppHelpTemplate:  help_template,
		UseShortOptionHandling: true,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
