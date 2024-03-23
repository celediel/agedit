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

// Decrypt decrypts bytes from filename
func Decrypt(filename string, identities ...age.Identity) ([]byte, error) {
	var (
		f   *os.File
		r   io.Reader
		err error
		out = &bytes.Buffer{}
	)
	if f, err = os.Open(filename); err != nil {
		return nil, err
	}

	if r, err = age.Decrypt(f, identities...); err != nil {
		return nil, err
	}

	if _, err := io.Copy(out, r); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// NewIdentity generates a new Age identity
func NewIdentity() (*age.X25519Identity, error) {
	id, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, err
	}

	return id, nil
}

// ReadIdentityFromFile reads the identity from the supplied filename
func ReadIdentityFromFile(filename string) (*age.X25519Identity, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	id, err := age.ParseX25519Identity(string(bytes))
	if err != nil {
		return nil, err
	}

	return id, nil
}

// WriteIdentityToFile writes the supplied identity to the supplied filename
func WriteIdentityToFile(id *age.X25519Identity, filename string) error {
	err := os.WriteFile(filename, []byte(id.String()), fs.FileMode(0600))
	if err != nil {
		return err
	}

	return nil
}
