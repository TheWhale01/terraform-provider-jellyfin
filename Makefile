test:
	go test -v ./...

testacc:
	@echo "==> Starting up jellyfin test env..."
	docker compose up --build -d
	@echo "==> Running Terraform acceptance tests"
	@sleep 2s
	TF_ACC=1 go test -v -cover ./... || (docker compose down -v && exit 1)
	@echo "==> Shutting down test env..."
	docker compose down

.PHONY: test testacc
