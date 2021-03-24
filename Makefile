cli:
	go build -mod vendor -o bin/server cmd/server/main.go

docker:
	cp $(DATABASE) whosonfirst.db
	docker build -f Dockerfile -t $(CONTAINER) .
	rm whosonfirst.db
