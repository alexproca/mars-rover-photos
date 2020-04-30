run:
	@go run main.go

clean:
	-@rm app
	-@rm server

build:
	go build -o app

build-linux:
	-@rm server
	env GOOS=linux GOARCH=amd64 go build -o server

build-docker: build-linux
	docker build -t nasa-api .

run-docker: build-docker stop-docker
	docker run --rm --name="nasa-api" -p "127.0.0.1:8081:8080" -d nasa-api

stop-docker:
	-docker rm -f nasa-api


test:
	go test -v ./...

deploy-heroku: build-linux
	heroku container:push web --app nasa-rover-photos
	heroku container:release web --app nasa-rover-photos

