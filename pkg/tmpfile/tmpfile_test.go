package tmpfile

import (
	"io/fs"
	"os"
	"testing"
)

var generator = NewGenerator("test_", ".txt", 18)

// TestCanCreateTmpFile tests if temporary files can be created and removed successfully
func TestCanCreateTmpFile(t *testing.T) {
	b := []byte{104, 101, 108, 108, 111, 32, 116, 104, 101, 114, 101}

	for range 1000 {
		outfile := generator.GenerateFullPath()
		err := os.WriteFile(outfile, b, fs.FileMode(0600))
		if err != nil {
			t.Fatal(err)
		}

		if _, err = os.Stat(outfile); err != nil && os.IsNotExist(err) {
			t.Fatal(err)
		}

		if err = os.Remove(outfile); err != nil {
			t.Fatal(err)
		}
	}
}

// TestUniqueTmpFile generates a large number of random names to make sure they're all unique
func TestUniqueTmpFile(t *testing.T) {
	var generated_names = map[string]string{}

	for range 100000 {
		name := generator.GenerateName()
		if val, ok := generated_names[name]; ok {
			t.Fatal("Non unique name", val)
		}
		generated_names[name] = name
	}
}
