include .env
export

build: clean
	go build -o ./bin/hibari

clean:
	rm -rf ./bin
	
run-debug: build
	./bin/hibari -d

run: build
	./bin/hibari

docker: build
	podman build --no-cache -t ghcr.io/pacsui/hibari:latest .
	podman push ghcr.io/pacsui/hibari:latest
