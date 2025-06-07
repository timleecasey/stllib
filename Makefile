
fake:

sim: fake
	go generate ./...
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux  go build $(GOFLAGS) -o ./sim simulator/sim.go
