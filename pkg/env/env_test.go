package env

import (
	"os"
	"testing"
)

var (
	editors = []string{"hx", "nano", "vi", "vim", "nvim", "micro", "emacs", "ed"}
)

func clearEnvForNow() {
	for _, item := range []string{"EDITOR", "VISUAL"} {
		os.Unsetenv(item)
	}
}

func TestEditorFromEnv(t *testing.T) {
	for _, item := range []string{"EDITOR", "VISUAL"} {
		clearEnvForNow()
		for _, editor := range editors {
			if err := os.Setenv(item, editor); err != nil {
				t.Fatal(err)
			}
			if got := GetEditor(); got != editor {
				t.Fatal("got", got, "but wanted", editor)
			}
		}
	}
}
