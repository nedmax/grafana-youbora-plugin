setup:
	yarn install
	yarn build
	go get -u github.com/grafana/grafana-plugin-sdk-go
	go mod tidy

build:
	yarn build
	mage -v

run: build
	docker compose up grafana

package:
	yarn install --pure-lockfile
	yarn build
	mv dist/ nedimar-youbora-datasource
	zip nedimar-youbora-datasource-$(shell cat package.json| jq -r '.version').zip nedimar-youbora-datasource -r
	mv nedimar-youbora-datasource dist