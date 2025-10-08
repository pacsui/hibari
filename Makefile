build: clean
	go build -o ./bin/hibari

clean:
	rm -rf ./bin
	
run: build
	./bin/hibari -d
docker: build
	podman build --no-cache -t registry.gitlab.com/itspacchu/containerdump:threadingbot_rc .
	podman push registry.gitlab.com/itspacchu/containerdump:threadingbot_rc
