setup:
	yard install

build:
	yarn dev

console:
	docker compose run console

run: build
	docker compose up grafana
