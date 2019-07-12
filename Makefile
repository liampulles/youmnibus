build:
	go build ./...
install: build
	go install ./...
docker-build: install
	docker build -t youmnibus-burden -f docker/Dockerfile.youmnibus-burden .
	docker build -t youmnibus-capture -f docker/Dockerfile.youmnibus-capture .
	docker build -t youmnibus-query -f docker/Dockerfile.youmnibus-query .
docker-push:
	docker tag youmnibus-burden:latest lpulles/youmnibus-burden:0.1
	docker tag youmnibus-capture:latest lpulles/youmnibus-capture:0.1
	docker tag youmnibus-query:latest lpulles/youmnibus-query:0.1
	docker push lpulles/youmnibus-burden:0.1
	docker push lpulles/youmnibus-capture:0.1
	docker push lpulles/youmnibus-query:0.1
launch-clean:
	docker rm mongo rabbitmq memcache-subscribers memcache-views memcache-videos
launch-mongo-docker:
	docker run --name mongo -v /data/db:/data/db -d -p 27017:27017 mongo
launch-rabbitmq-docker:
	docker run -d --hostname localhost --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3-management
launch-memcache-docker:
	docker run -d --name memcache-subscribers -p 11211:11211 memcached
	docker run -d --name memcache-views -p 11212:11211 memcached
	docker run -d --name memcache-videos -p 11213:11211 memcached
