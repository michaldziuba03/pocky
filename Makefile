.PHONY: build clean

build:
	@mkdir -p build
	go build -o build/pocky main.go

clean:
	@rm -rf build

all:
	build
