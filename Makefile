SERVER=root@server

all: build upload

build:
	@echo "Building"
	@env GOOS=linux GOARCH=amd64 go build -v github.com/nicklasos/ws

upload:
	@echo "Uploading"
	@scp ws $(SERVER):/var/www/websockets/ws

upload-all:
	@scp worker.conf $(SERVER):/var/www/websockets/worker.conf
	@scp restart $(SERVER):/var/www/websockets/restart
	@scp websockets.html $(SERVER):/var/www/websockets/websockets.html
	@scp .env.example $(SERVER):/var/www/websockets/.env.example
	@scp server.crt $(SERVER):/var/www/websockets/server.crt
	@scp server.key $(SERVER):/var/www/websockets/server.key

.PHONY: build
