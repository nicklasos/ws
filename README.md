# Install
```
go get github.com/gorilla/websocket
go get github.com/joho/godotenv
go get github.com/streadway/amqp
go get github.com/go-redis/redis
go get github.com/codegangsta/gin
```

# Supervisor
```
sudo apt-get install supervisor
vim .../worker.conf
mv .../worker.conf /etc/supervisor/conf.d/worker.conf
chmod +x /etc/supervisor/conf.d/worker.conf
sudo supervisorctl reread
sudo supervisorctl update
sudo supervisorctl start ws-worker:*
```