package editor

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"

	"git.burning.moe/celediel/agedit/pkg/tmpfile"
)

// EditFile opens the specified file in the configured editor
func EditFile(editor, filename string) error {
	if editor == "" {
		return errors.New("editor not set")
	}

	// TODO: handle editors that require arguments
	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// EditTempFile creates a temporary file with a random name, opens it in the
// editor, and returns the byte slice of its contents.
func EditTempFile(editor, start, prefix, suffix string, filename_length int) ([]byte, error) {
	var (
		filename string
		bytes    []byte
		err      error
		file     *os.File
	)

	// generator := tmpfile.NewGenerator("agedit_", ".txt", 13)
	generator := tmpfile.NewGenerator(prefix, suffix, filename_length)

	filename = generator.GenerateFullPath()
	if file, err = os.Create(filename); err != nil {
		return nil, err
	}

	if err = os.WriteFile(filename, []byte(start), fs.FileMode(0600)); err != nil {
		return nil, err
	}

	if err = EditFile(editor, filename); err != nil {
		return nil, err
	}

	if bytes, err = os.ReadFile(filename); err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	if err = os.Remove(filename); err != nil {
		return nil, err
	}

	return bytes, nil
}
