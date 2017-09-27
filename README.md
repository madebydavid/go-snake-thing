# Golang scratch-pad

```bash
vagrant plugin install vagrant-vbguest
vagrant up
vagrant ssh
go get -u github.com/buger/goterm
go get -u github.com/go-redis/redis
go get -u github.com/golang/protobuf/protoc-gen-go
go install snakes
snakes
```


```bash
cd /vagrant/src/snakedata
protoc --go_out=. snakes.proto
```