package tmpfile

import (
	"math/rand"
	"os"
	"strings"

	"git.burning.moe/celediel/agedit/pkg/env"
)

const chars string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

type Generator struct {
	Prefix, Suffix string
	Length         int
}

// GenerateName generates a random temporary filename like agedit_geef0XYC30RGV
func (g *Generator) GenerateName() string {
	return g.Prefix + randomString(chars, g.Length) + g.Suffix
}

// GenerateFullPath generates a random temporary filename and appends it to the OS's temporary directory
func (g *Generator) GenerateFullPath() string {
	return env.GetTempDirectory() + string(os.PathSeparator) + g.GenerateName()
}

// NewGenerator returns a new Generator
func NewGenerator(prefix, suffix string, length int) Generator {
	return Generator{
		Prefix: prefix,
		Suffix: suffix,
		Length: length,
	}
}

func randomString(set string, length int) string {
	out := strings.Builder{}
	for i := 0; i < length; i++ {
		out.WriteByte(randomChar(set))
	}
	return out.String()
}

func randomChar(set string) byte {
	return set[rand.Intn(len(set))]
}
