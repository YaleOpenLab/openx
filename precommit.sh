# Pre-commit checks, can add them in a git hook if desired
go fmt ./...
go vet ./...
./misspell -w -q ./...
./staticcheck github.com/YaleOpenLab/openx/...
