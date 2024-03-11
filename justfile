binary := "agedit"
build_dir := "bin"
cmd := "." / "cmd" / binary
output := "." / build_dir / binary

# do the thing
default: test check build install

# build binary
build:
	go build -o {{output}} {{cmd}}

# run from source
run:
	go run {{cmd}}

# build 'n run
run-binary: build
	exec {{output}}

# run with args
run-args args:
	go run {{cmd}} {{args}}

# install binary into $GOPATH
install:
	go install {{cmd}}

# clean up after yourself
clean:
	rm {{output}}

# run go tests
test:
	go test ./...

# run linter
check:
	staticcheck ./...
