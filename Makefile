
fake:

sim: fake
	go generate ./...
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux  go build $(GOFLAGS) -o ./sim2 simulator/sim.go
