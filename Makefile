build: *.go
	go build -o dynamic-route53 *.go

run: build
	./dynamic-route53

test:
	go test ./...

docker:
	docker build -t dynamic-route53:latest .
.PHONY: docker