.PHONY: build clean

build:
	@mkdir -p build
	go build -o build/pocky

clean:
	@rm -rf build

all:
	build
