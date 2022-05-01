clean:
	mage clean
	rm -f dist

test:
	yarn test
	mage test
	
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

package_node:
	yarn install --pure-lockfile
	yarn build

package_linux: package_node
	mage build:linux
	yarn sign
	mv dist/ nedimar-youbora-datasource
	zip nedimar-youbora-datasource-$(shell cat package.json| jq -r '.version').linux.zip nedimar-youbora-datasource -r
	mv nedimar-youbora-datasource dist
	mage clean

package_windows: package_node
	mage build:windows
	yarn sign
	mv dist/ nedimar-youbora-datasource
	zip nedimar-youbora-datasource-$(shell cat package.json| jq -r '.version').windows.zip nedimar-youbora-datasource -r
	mv nedimar-youbora-datasource dist
	mage clean

package_darwin: package_node
	mage build:darwin
	yarn sign
	mv dist/ nedimar-youbora-datasource
	zip nedimar-youbora-datasource-$(shell cat package.json| jq -r '.version').darwin.zip nedimar-youbora-datasource -r
	mv nedimar-youbora-datasource dist
	mage clean

package: package_linux package_windows package_darwin
