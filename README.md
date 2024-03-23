# agedit

open an [age](https://github.com/FiloSottile/age) encrypted file in $EDITOR

## cli
`go install git.burning.moe/celediel/agedit/cmd/agedit@latest`

`agedit [flags] filename`

### flags

```text
   --identity identity, -I identity [ --identity identity, -I identity ]        age identity (or identities) to decrypt with
   --identity-file FILE, -i FILE                                                read identity from FILE
   --recipient recipient, -R recipient [ --recipient recipient, -R recipient ]  age recipients to encrypt to
   --recipient-file FILE, -r FILE                                               read recipients from FILE
   --out FILE, -o FILE                                                          write to FILE instead of the input file
   --editor EDITOR, -e EDITOR                                                   edit with specified EDITOR instead of $EDITOR
   --editor-args arg [ --editor-args arg ]                                      arguments to send to the editor
   --force, -f                                                                  re-encrypt the file even if no changes have been made. (default: false)
   --log level                                                                  log level (default: "warn")
   --help, -h                                                                   show help
   --version, -v                                                                print the version
```

## library
`go get git.burning.moe/celediel/agedit@latest`

See `./cmd/agedit` for example usage.

## TODO
- support for password encrypted key files
