# Install dev

# Dev
```
go get github.com/codegangsta/gin
gin
go to http://localhost:3000 and http://localhost:3000/send
```

# Supervisor
```
sudo apt-get install supervisor
vim worker.conf
cp worker.conf /etc/supervisor/conf.d/worker.conf
chmod +x /etc/supervisor/conf.d/worker.conf
sudo supervisorctl reread
sudo supervisorctl update
sudo supervisorctl start ws-worker:*
```

# Client
```js
// id: client id
// rooms=Ukraine&rooms=Minecraft: join 2 channels
const socket = new WebSocket("{{ .Schema }}://localhost:{{.Port}}/ws?id=foo&rooms=Ukraine&rooms=Minecraft");
```


# Stats
curl https://ws/stats
```json
{
  "connections": 88,
  "users": 77,
  "users_1min": 44,
  "users_5min": 77,
  "users_15min": 77,
  "rooms": {
    "Minecraft": 48,
    "Ukraine": 88
  }
}
```