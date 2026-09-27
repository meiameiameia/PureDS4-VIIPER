set windows-shell := ["powershell.exe", "-NoProfile", "-Command"]
version := env_var_or_default("VIIPER_VERSION", "0.1.0-pureds4-dev")
ldflags := "-s -w -X main.Version=" + version + " -X github.com/Alia5/VIIPER/internal/codegen/common.Version=" + version

default:
	just --list

test:
	go test -count=1 ./...

vet:
	go vet ./...

build:
	go build -tags release -trimpath -ldflags '{{ ldflags }}' -o viiper.exe ./cmd/viiper
