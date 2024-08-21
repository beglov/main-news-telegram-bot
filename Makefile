.PHONY: build ansible deploy lint-ansible lint-fix-ansible

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
