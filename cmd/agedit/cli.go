package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"git.burning.moe/celediel/agedit/internal/config"
	"git.burning.moe/celediel/agedit/pkg/editor"
	"git.burning.moe/celediel/agedit/pkg/encrypt"
	"git.burning.moe/celediel/agedit/pkg/env"

	"github.com/charmbracelet/log"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/urfave/cli/v2"
)

const (
	name          string = "agedit"
	usage         string = "Edit age encrypted files with your $EDITOR"
	version       string = "0.0.2"
	help_template string = `NAME:
   {{.Name}} {{if .Version}}v{{.Version}}{{end}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[flags]{{end}} [filename]
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
FLAGS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}
`
)

var (
	authors = []*cli.Author{{
		Name:  "Lilian Jónsdóttir",
		Email: "lilian.jonsdottir@gmail.com",
	}}

	flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "identity",
			Usage:   "age identity file to decrypt with",
			Aliases: []string{"i"},
			Action: func(ctx *cli.Context, identity_file string) error {
				if identity_file != "" {
					cfg.IdentityFile = identity_file
				}
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "out",
			Usage:   "write to this file instead of the input file",
			Aliases: []string{"o"},
			Action: func(ctx *cli.Context, out string) error {
				output_file = out
				return nil
			},
		},
		&cli.BoolFlag{
			Name:    "force",
			Usage:   "Re-encrypt the file even if no changes have been made.",
			Aliases: []string{"f"},
			Action: func(ctx *cli.Context, b bool) error {
				force_overwrite = b
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "log",
			Usage:   "log level",
			Value:   "warn",
			Aliases: []string{"l"},
			Action: func(ctx *cli.Context, s string) error {
				if lvl, err := log.ParseLevel(s); err == nil {
					logger.SetLevel(lvl)
					// Some extra info for debug level
					if logger.GetLevel() == log.DebugLevel {
						logger.SetReportCaller(true)
					}
				} else {
					logger.SetLevel(log.WarnLevel)
				}
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "editor",
			Usage:   "specify the editor to use",
			Aliases: []string{"e"},
			Action: func(ctx *cli.Context, editor string) error {
				cfg.Editor = editor
				return nil
			},
		},
	}
)

// before validates input, does some setup, and loads config info from file
func before(ctx *cli.Context) error {
	// check input
	if input_file = strings.Join(ctx.Args().Slice(), " "); input_file == "" {
		return fmt.Errorf("no file to edit, use " + name + " -h for help")
	}

	// do some setup
	cfg = config.Defaults
	cfg.Editor = env.GetEditor()
	cfg_dir := env.GetConfigDir(name)
	cfg.IdentityFile = cfg_dir + "identity.key"
	configFile = cfg_dir + name + ".yaml"
	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
	})

	// load config from file
	if _, err := os.Stat(configFile); err != nil && errors.Is(err, os.ErrNotExist) {
		// or not
		logger.Debug("couldn't load config file", "file", configFile)
	} else {
		err = cleanenv.ReadConfig(configFile, &cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

// action does the actual thing
func action(ctx *cli.Context) error {
	if _, err := os.Stat(input_file); os.IsNotExist(err) {
		return err
	}

	if output_file == "" {
		output_file = input_file
		logger.Debug("out file not specified, using input", "outfile", output_file)
	}

	if _, err := os.Stat(cfg.IdentityFile); os.IsNotExist(err) {
		return fmt.Errorf("identity file unset, use -i or set one in the config file")
	}

	if id, err := encrypt.ReadIdentityFromFile(cfg.IdentityFile); err != nil {
		return err
	} else {
		identity = id
	}
	logger.Debug("read identity from file", "id", identity.Recipient())

	decrypted, err := encrypt.Decrypt(input_file, identity)
	if err != nil {
		return err
	}
	logger.Debug("decrypted " + input_file + " sucessfully")

	edited, err := editor.EditTempFile(cfg.Editor, string(decrypted), cfg.Prefix, cfg.Suffix, cfg.RandomLength)
	if err != nil {
		return err
	}
	logger.Debug("got data back from editor")

	// don't overwrite same data, unless specified
	if string(edited) == string(decrypted) && !force_overwrite {
		logger.Warn("No edits made, not writing " + output_file)
		return nil
	}

	err = encrypt.Encrypt(edited, output_file, identity)
	if err != nil {
		return err
	}
	logger.Debug("re-encrypted to " + output_file)

	return nil
}
