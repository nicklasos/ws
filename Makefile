SERVER=root@giveaway-ws
#SERVER=root@giveaway-on-cases-ws
#SERVER=root@pff-ws
#SERVER=root@ge-ws

all: build upload

build:
	@echo "Building"
	@env GOOS=linux GOARCH=amd64 go build -v github.com/nicklasos/ws

upload: build
	@echo "Uploading"
	@scp ws $(SERVER):/var/www/websockets/ws

upload-all: upload
	@echo "Uploading assets"
	@scp worker.conf $(SERVER):/var/www/websockets/worker.conf
	@scp restart $(SERVER):/var/www/websockets/restart
	@scp Makefile $(SERVER):/var/www/websockets/Makefile
	@scp websockets.html $(SERVER):/var/www/websockets/websockets.html
	@scp .env.example $(SERVER):/var/www/websockets/.env.example
	@scp certs/cf-ge.crt $(SERVER):/var/www/websockets/cf-ge.crt
	@scp certs/cf-ge.crt $(SERVER):/var/www/websockets/cf-ge.key

reload:
	@rm ws_copy
	@cp ws ws_copy
	@rm websockets_copy
	@cp websockets websockets_copy
	@mv ws websockets
	@sudo supervisorctl restart ws-worker:*

.PHONY: build
