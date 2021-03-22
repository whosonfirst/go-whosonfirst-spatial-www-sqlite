cli:
	go build -mod vendor -o bin/server cmd/server/main.go

docker:
	cp $(DATABASE) whosonfirst.db
	docker build --build-arg DATABASE=whosonfirt.db -f Dockerfile -t spatial-www-sqlite .
	rm whosonfirst.db
