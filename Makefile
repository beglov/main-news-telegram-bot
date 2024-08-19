.PHONY: build ansible deploy

build:
	go build -o bot cmd/bot/bot.go
	go build -o send_news cmd/send_news/send_news.go

ansible:
	ansible-playbook -i ansible/inventories/staging.yml ansible/provision.yml

deploy: build ansible
