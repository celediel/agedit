package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"filippo.io/age"
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
	version       string = "0.1.1"
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
		&cli.StringSliceFlag{
			Name:    "identity",
			Usage:   "age `identity` (or identities) to decrypt with",
			Aliases: []string{"I"},
			Action: func(ctx *cli.Context, inputs []string) error {
				for _, input := range inputs {
					id, err := age.ParseX25519Identity(input)
					if err != nil {
						return err
					}

					identities = append(identities, id)
				}
				gave_identities = true
				return nil
			},
		},
		&cli.PathFlag{ // I dunno why PathFlag exists because cli.Path is just string
			Name:    "identity-file",
			Usage:   "read identity from `FILE`",
			Aliases: []string{"i"},
			Action: func(ctx *cli.Context, identity_file cli.Path) error {
				if identity_file != "" {
					cfg.IdentityFile = identity_file
				}
				return nil
			},
		},
		&cli.StringSliceFlag{
			Name:    "recipient",
			Usage:   "age `recipient`s to encrypt to",
			Aliases: []string{"R"},
			Action: func(ctx *cli.Context, inputs []string) error {
				for _, input := range inputs {
					logger.Debugf("parsing public key from string %s", input)
					r, err := age.ParseX25519Recipient(input)
					if err != nil {
						return err
					}
					recipients = append(recipients, r)
				}
				gave_recipients = true
				return nil
			},
		},
		&cli.PathFlag{
			Name:    "recipient-file",
			Usage:   "read recipients from `FILE`",
			Aliases: []string{"r"},
			Action: func(ctx *cli.Context, recipient_file cli.Path) error {
				if recipient_file != "" {
					cfg.RecipientFile = recipient_file
				}
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "out",
			Usage:   "write to `FILE` instead of the input file",
			Aliases: []string{"o"},
			Action: func(ctx *cli.Context, out string) error {
				output_file = out
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "editor",
			Usage:   "edit with specified `EDITOR` instead of $EDITOR",
			Aliases: []string{"e"},
			Action: func(ctx *cli.Context, editor string) error {
				cfg.Editor = editor
				return nil
			},
		},
		&cli.StringSliceFlag{
			Name:  "editor-args",
			Usage: "`arg`uments to send to the editor",
			Action: func(ctx *cli.Context, args []string) error {
				cfg.EditorArgs = args
				return nil
			},
		},
		&cli.BoolFlag{
			Name:               "force",
			Usage:              "re-encrypt the file even if no changes have been made",
			Aliases:            []string{"f"},
			DisableDefaultText: true,
			Action: func(ctx *cli.Context, b bool) error {
				force_overwrite = b
				return nil
			},
		},
		&cli.StringFlag{
			Name:  "log",
			Usage: "log `level`",
			Value: "warn",
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
	}
)

// before validates input, does some setup, and loads config info from file
func before(ctx *cli.Context) error {
	// check input
	if input_file = strings.Join(ctx.Args().Slice(), " "); input_file == "" {
		return fmt.Errorf("no file to edit, use " + name + " -h for help")
	}

	// set some defaults
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

	// setup editor with loaded config options
	edt = editor.New(cfg.Editor, cfg.EditorArgs, cfg.Prefix, cfg.Suffix, cfg.RandomLength)

	return nil
}

// action does the actual thing
func action(ctx *cli.Context) error {
	// make sure input file exists
	if _, err := os.Stat(input_file); os.IsNotExist(err) {
		return err
	}

	if output_file == "" {
		output_file = input_file
		logger.Debug("out file not specified, using input", "outfile", output_file)
	}

	// read from identity file if exists and no identities have been supplied
	if !gave_identities {
		if _, err := os.Stat(cfg.IdentityFile); os.IsNotExist(err) {
			return fmt.Errorf("identity file unset and no identities supplied, use -i to specify an idenitity file or set one in the config file, or use -I to specify an age private key")
		} else {
			f, err := os.Open(cfg.IdentityFile)
			if err != nil {
				return fmt.Errorf("couldn't open identity file: %v", err)
			}
			if ids, err := age.ParseIdentities(f); err != nil {
				return fmt.Errorf("couldn't parse identities: %v", err)
			} else {
				identities = append(identities, ids...)
			}
		}
	}

	// read from recipient file if it exists and no recipients have been supplied
	if !gave_recipients && cfg.RecipientFile != "" {
		if _, err := os.Stat(cfg.RecipientFile); os.IsNotExist(err) {
			return fmt.Errorf("recipient file doesn't exist")
		} else {
			f, err := os.Open(cfg.RecipientFile)
			if err != nil {
				return fmt.Errorf("couldn't open recipient file: %v", err)
			}
			if rs, err := age.ParseRecipients(f); err != nil {
				return fmt.Errorf("couldn't parse recipients: %v", err)
			} else {
				recipients = append(recipients, rs...)
			}
		}
	}

	// get recipients from specified identities
	for _, id := range identities {
		if actual_id, ok := id.(*age.X25519Identity); ok {
			recipients = append(recipients, actual_id.Recipient())
		}
	}

	// try to decrypt the file
	decrypted, err := encrypt.Decrypt(input_file, identities...)
	if err != nil {
		return err
	}
	logger.Debug("decrypted " + input_file + " sucessfully")

	// open decrypted data in the editor
	edited, err := edt.EditTempFile(string(decrypted))
	if err != nil {
		return err
	}
	logger.Debug("got data back from editor")

	// don't overwrite same data, unless specified
	if string(edited) == string(decrypted) && !force_overwrite {
		logger.Warn("No edits made, not writing " + output_file)
		return nil
	}

	// actually re-encrypt the data
	err = encrypt.Encrypt(edited, output_file, recipients...)
	if err != nil {
		return err
	}
	logger.Debug("re-encrypted to " + output_file)

	return nil
}
