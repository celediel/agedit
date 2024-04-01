package decrypt

import (
	"bytes"
	"io"
	"os"

	"filippo.io/age"
)

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
