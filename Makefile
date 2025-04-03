up: down
	docker compose up --build -d

down:
	docker compose down

clean:
	docker compose down -v

lint:
	golangci-lint run
