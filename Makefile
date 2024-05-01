
all:
	CGO_ENABLED=0 go build -ldflags "-X github.com/mangofeet/netrunner-alt-gen/cmd.version=$$(git describe --tags)"
