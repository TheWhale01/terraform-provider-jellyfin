NAME=terraform-provider-jellyfin

$(NAME):
	go build -o $(NAME)

clean:
	rm -rf $(NAME)

re: clean $(NAME)

test:
	go test -v ./...

testacc:
	@echo "==> Starting up jellyfin test env..."
	docker compose up --build -d
	@echo "==> Running Terraform acceptance tests"
	@sleep 2s
	TF_ACC=1 go test -v -cover ./internal/provider

.PHONY: test testacc clean re
