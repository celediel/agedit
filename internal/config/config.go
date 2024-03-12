package config

type Config struct {
	IdentityFile string `json:"identityfile" yaml:"identityfile" toml:"identityfile"`
	Editor       string `json:"editor" yaml:"editor" toml:"editor"`
	Prefix       string `json:"randomfileprefix" yaml:"randomfileprefix" toml:"randomfileprefix"`
	Suffix       string `json:"randomfilesuffix" yaml:"randomfilesuffix" toml:"randomfilesuffix"`
	RandomLength int    `json:"randomfilenamelength" yaml:"randomfilenamelength" toml:"randomfilenamelength"`
}

var Defaults = Config{
	Prefix:       "agedit_",
	Suffix:       ".txt",
	RandomLength: 13,
}
