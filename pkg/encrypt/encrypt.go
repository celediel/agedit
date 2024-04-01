package encrypt

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"

	"filippo.io/age"
)

// Encrypt encrypts bytes into filename
func Encrypt(data []byte, filename string, recipients ...age.Recipient) error {
	var (
		w   io.WriteCloser
		out = &bytes.Buffer{}
		err error
	)

	if len(recipients) == 0 {
		return errors.New("no recepients? who's trying to encrypt?")
	}

	if w, err = age.Encrypt(out, recipients...); err != nil {
		return err
	}

	io.WriteString(w, string(data))
	if err = w.Close(); err != nil {
		return err
	}

	os.Truncate(filename, 0) // in case it exists already
	if err = os.WriteFile(filename, out.Bytes(), fs.FileMode(0600)); err != nil {
		return err
	}

	return nil
}
