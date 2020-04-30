run:
	-@go run main.go

clean:
	-@rm app
	-@rm server

build:
	go build -o app

build-linux:
	env GOOS=linux GOARCH=am6d4 go build -o server

build-docker: build-linux
	docker build -t nasa-api .

run-docker:
	docker run --rm --name="nasa-api" -p "127.0.0.1:8080:8080" nasa-api

stop-docker:
	docker rm -f nasa-api


test:
	go test -v ./...