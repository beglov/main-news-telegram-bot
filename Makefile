.PHONY: test cover format lint build ansible deploy lint-ansible lint-fix-ansible

test:
	go test -cover ./...

cover:
	go test -coverprofile=/tmp/coverage.txt ./...
	go tool cover -html=/tmp/coverage.txt

# Format code with gofumpt (https://github.com/mvdan/gofumpt).
format:
	gofumpt -l -w .

# Lint code with golangci-lint (https://github.com/golangci/golangci-lint).
lint:
	docker run -t --rm -v $(PWD):/app -w /app golangci/golangci-lint:v2.0.2 golangci-lint run

build:
	go build -o bot cmd/bot/bot.go
	go build -o send_news cmd/send_news/send_news.go

ansible:
	ansible-playbook -i ansible/inventories/staging.yml ansible/provision.yml

deploy: build ansible

lint-ansible:
	docker run -it --rm -v ${PWD}:/mnt haxorof/ansible-lint ansible

lint-fix-ansible:
	docker run -it --rm -v ${PWD}:/mnt haxorof/ansible-lint --fix ansible
