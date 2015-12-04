all: build

build: chpass_linux chpass_windows chpass_mac

test:
	go test -v ./...

clean:
	rm sourcescanner/linux_x86_64/bin/*
	rm sourcescanner/windows_x86_64/bin/*
	rm tmp/*.mark

chpass_linux:
	GOOS=linux go build -o chpassh/bin/chpassh.linux chpassh.go

chpass_windows:
	GOOS=windows go build -o chpassh/bin/chpassh.exe chpassh.go

chpass_mac:
	GOOS=darwin go build -o chpassh/bin/chpassh.mac chpassh.go

.PHONY: all build clean test chpass_linux chpass_windows chpass_mac
