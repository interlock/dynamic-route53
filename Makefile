build: *.go
	go build -o dynamic-route53 *.go

run: build
	./dynamic-route53