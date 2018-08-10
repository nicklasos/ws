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