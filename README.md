# agedit

open an age encrypted file in $EDITOR

## cli
`go install git.burning.moe/celediel/agedit/cmd/agedit@latest`

`agedit [flags] filename`

### flags

```text
   --identity value, -i value  age identity file to use
   --out value, -o value       write to this file instead of the input file
   --log value, -l value       log level (default: "warn")
   --editor value, -e value    specify the editor to use
   --help, -h                  show help
   --version, -v               print the version
```

## library
`go get git.burning.moe/celediel/agedit@latest`

See `./cmd/agedit` for example usage.
