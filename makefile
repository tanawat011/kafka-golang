def:
	go run main.go

dev:
	fresh -c fresh.conf

mod:
	GO111MODULE=on go mod tidy

start:
	docker-compose up -d --build
stop:
	docker-compose down