Install Golang: go1.12
More doc:
https://github.com/gin-gonic/gin#gin-v1-stable
https://github.com/json-iterator/go
http://jsoniter.com/migrate-from-go-std.html
https://github.com/go-redis/redis

## Librdkafka (https://github.com/confluentinc/confluent-kafka-go)
```bash
	git clone https://github.com/edenhill/librdkafka.git
	cd librdkafka
	./configure --prefix /usr
	make
	sudo make install
```
## Swager
	1. Docs -> https://github.com/swaggo/gin-swagger
	2. Init swagger -> swag init

## Chạy service:
	1. Tạo config/common.config theo config/common.example.config (đổi lại thông tin theo môi trường)
	2. make run_server