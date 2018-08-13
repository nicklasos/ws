SERVER=root@server

all: build upload

build:
	@echo "Building"
	@env GOOS=linux GOARCH=amd64 go build -v github.com/nicklasos/ws

upload:
	@echo "Uploading"
	@scp ws $(SERVER):/var/www/websockets/ws

upload-all: upload
	@echo "Uploading assets"
	@scp worker.conf $(SERVER):/var/www/websockets/worker.conf
	@scp restart $(SERVER):/var/www/websockets/restart
	@scp websockets.html $(SERVER):/var/www/websockets/websockets.html
	@scp .env.example $(SERVER):/var/www/websockets/.env.example
	@scp cf.crt $(SERVER):/var/www/websockets/server.crt
	@scp cf.key $(SERVER):/var/www/websockets/server.key

reload:
	@cp ws websockets
    @sudo supervisorctl restart ws-worker:*

.PHONY: build

