package editor

import (
	"io/fs"
	"os"
	"os/exec"

	"git.burning.moe/celediel/agedit/pkg/tmpfile"
)

type Editor struct {
	Command   string
	Args      []string
	generator tmpfile.Generator
}

// EditFile opens the specified file in the configured editor
func (e *Editor) EditFile(filename string) error {
	args := append(e.Args, filename)

	cmd := exec.Command(e.Command, args...)
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
func (e *Editor) EditTempFile(start string) ([]byte, error) {
	var (
		filename string
		bytes    []byte
		err      error
		file     *os.File
	)

	filename = e.generator.GenerateFullPath()
	if file, err = os.Create(filename); err != nil {
		return nil, err
	}

	if err = os.WriteFile(filename, []byte(start), fs.FileMode(0600)); err != nil {
		return nil, err
	}

	if err = e.EditFile(filename); err != nil {
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

// New returns an Editor configured to open files with `command` + `args`.
// The prefix and suffix will be added to the randomly generated
// filename of `length` characters.
func New(command string, args []string, prefix, suffix string, length int) Editor {
	return Editor{
		Command:   command,
		Args:      args,
		generator: tmpfile.NewGenerator(prefix, suffix, length),
	}
}
