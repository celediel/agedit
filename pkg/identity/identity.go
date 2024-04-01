package identity

import (
	"io/fs"
	"os"

	"filippo.io/age"
)

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
