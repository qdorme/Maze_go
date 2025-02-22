ifeq ($(OS),Windows_NT)
	SHELL := pwsh.exe
	.SHELLFLAGS := -NoProfile -Command
endif

clean:
	go clean

build: clean
	go mod tidy
	go build -o ./maze.exe cmd/main.go

install: build
	Remove-Item -Force -Recurse -ErrorAction Ignore ${CURDIR}/build; $$null
	New-Item ${CURDIR}/build -ItemType Directory
	Move-Item -Path ${CURDIR}/maze.exe -Destination ${CURDIR}/build/maze.exe
	Copy-Item -Path ${CURDIR}/public -Destination ${CURDIR}/build/public -Recurse

Run: install
	./build/maze.exe
