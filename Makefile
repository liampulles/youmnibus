build:
	go build ./...
install: build
	go install ./...
dockerize:
	docker build -t youmnibus-base -f Dockerfile.youmnibus-base .
	docker build -t youmnibus-capture -f Dockerfile.youmnibus-capture .
	docker build -t youmnibus-query -f Dockerfile.youmnibus-query .
launch-mongo-docker:
	docker run --name mongo -v /data/db:/data/db -d -p 27017:27017 mongo
launch-rabbitmq-docker:
	docker run -d --hostname localhost --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3-management
launch-memcache-docker:
	docker run -d --name memcache-subscribers -p 11211:11211 memcached
	docker run -d --name memcache-views -p 11212:11211 memcached
	docker run -d --name memcache-videos -p 11213:11211 memcached
