# agedit

open an age encrypted file in $EDITOR

## cli
`go install git.burning.moe/celediel/agedit@latest`

`agedit [flags] filename`

### flags

```   --identity value, -i value  age identity file to use
   --out value, -o value       Write to this file instead of the input file
   --log value, -l value       log level (default: "warn")
   --help, -h                  show help
   --version, -v               print the version
```

## library
`go get git.burning.moe/celediel/agedit@latest`

See `./cmd/agedit/agedit.go` for example usage.
