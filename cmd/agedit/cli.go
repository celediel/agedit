package main

import (
	"errors"
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
	name          = "agedit"
	usage         = "Edit age encrypted files with your $EDITOR"
	version       = "0.0.2"
	help_template = `NAME:
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
			Usage:   "age identity file to use",
			Aliases: []string{"i"},
			Action: func(ctx *cli.Context, s string) error {
				if identity_file := ctx.String("identity"); identity_file != "" {
					cfg.IdentityFile = identity_file
				}
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "out",
			Usage:   "write to this file instead of the input file",
			Aliases: []string{"o"},
			Action: func(ctx *cli.Context, s string) error {
				output_file = ctx.String("out")
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "log",
			Usage:   "log level",
			Value:   "warn",
			Aliases: []string{"l"},
			Action: func(ctx *cli.Context, s string) error {
				if lvl, err := log.ParseLevel(ctx.String("log")); err == nil {
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
			Action: func(ctx *cli.Context, s string) error {
				cfg.Editor = ctx.String("editor")
				return nil
			},
		},
	}
)

// before validates input, does some setup, and loads config info from file
func before(ctx *cli.Context) error {
	// check input
	if input_file = strings.Join(ctx.Args().Slice(), " "); input_file == "" {
		return errors.New("no file to edit, use agedit -h for help")
	}

	// do some setup
	cfg = config.Defaults
	cfg.Editor = env.GetEditor()
	cfg_dir := env.GetConfigDir("agedit")
	cfg.IdentityFile = cfg_dir + "identity.key"
	configFile = cfg_dir + "agedit.yaml"
	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
	})

	// load config from file
	_, err := os.Open(configFile)
	if err != nil && errors.Is(err, os.ErrNotExist) {
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
	if _, err := os.Open(input_file); os.IsNotExist(err) {
		return err
	}

	if output_file == "" {
		output_file = input_file
		logger.Debug("out file not specified, using input", "outfile", output_file)
	}

	if _, err := os.Open(cfg.IdentityFile); os.IsNotExist(err) {
		return errors.New("identity file unset, use -i or set one in the config file")
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

	if string(edited) == string(decrypted) {
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
