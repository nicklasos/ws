# Install dev
```
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

# Stats
curl https://ws/stats
```json
{
  "connections": 88,
  "users": 77,
  "users_1min": 44,
  "users_5min": 77,
  "users_15min": 77
}
```