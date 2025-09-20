default:
    just --list

test:
    go test -count 10 -shuffle on -race -v ./...

lint:
    golangci-lint run

fmt:
    golangci-lint fmt
