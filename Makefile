SHELL=/bin/bash

test_with_docker:
	podman run -v ./:/app -w /app public.ecr.aws/docker/library/golang:1.19.2 go test -v ./...

build_with_docker:
	podman run -e CGO_ENABLED=0 -e GOOS=linux -e GPARCH=amd64 -v ./:/app -w /app public.ecr.aws/docker/library/golang:1.19.2 go build -ldflags "-extldflags '-static'" -o ./build/gomsvc ./cmd/gomsvc

build_release_with_docker:
	podman run -e CGO_ENABLED=0 -e GOOS=linux -e GPARCH=amd64 -v ./:/app -w /app public.ecr.aws/docker/library/golang:1.19.2 go build -ldflags "-s -w -extldflags '-static'" -o ./build/gomsvc ./cmd/gomsvc

docker_build_github:
	test -n "$(tag)"
	podman build -t ghcr.io/inquizarus/gomsvc:$(tag) .

docker_push_github:
	test -n "$(tag)"
	podman push ghcr.io/inquizarus/gomsvc:$(tag)