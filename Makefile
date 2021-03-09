.PHONY: make

make: clean build/goqa

build/goqa:
	@go generate
	@cd cmd && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags prod -ldflags "-w -extldflags '-static'" -o ../build/goqa

clean:
	@rm -f build/goqa