dockerUp:
	docker compose up -d

dockerDown:
	docker compose down

down:
	go run . migrateDown

vet:
	go vet ./...

fmt:
	go fmt ./...

http: vet fmt
	go run . http