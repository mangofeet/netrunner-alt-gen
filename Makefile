
all:
	CGO_ENABLED=0 GOOS=linux go build -trimpath -v -ldflags "-X github.com/mangofeet/netrunner-alt-gen/cmd.version=$$(git describe --tags)"

docker:
	docker build -t mangofeet/netrunner-alg-gen .
