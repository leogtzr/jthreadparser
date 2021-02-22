.DEFAULT_GOAL := test

stringer:
	stringer -type=SynchronizerState

clean:
	go clean
	rm --force cp.out

test:
	go test

check:
	go test

cover:
	go test -coverprofile cp.out
	go tool cover -html=cp.out

lint:
	golangci-lint run --enable-all

