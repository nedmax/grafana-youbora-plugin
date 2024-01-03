GRAFANA_PLUGIN_ID := $(shell cat src/plugin.json | jq -r .id)
GRAFANA_PLUGIN_VERSION := $(shell cat package.json| jq -r '.version')
GRAFANA_PLUGIN_TYPE := $(shell cat src/plugin.json | jq -r .type)

MD5CMD ?= "md5"

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
	@yarn build
	@mage -v

build-with-docker:
	@docker compose run --remove-orphans go-builder
	@docker compose run --remove-orphans js-builder 

run:
	@docker compose up grafana

package_node:
	@yarn install --pure-lockfile
	@yarn build

check_os:
	@if [ -z "${OS}" ]; then echo "OS not set. Please run \`make {linux,windows,darwin} package\`"; exit 1; fi ||:

linux:
	@:$(eval OS=linux)
	@:$(eval GRAFANA_PLUGIN_ARTIFACT="${GRAFANA_PLUGIN_ID}-${GRAFANA_PLUGIN_VERSION}-${OS}.zip")
	@:$(eval GRAFANA_PLUGIN_ARTIFACT_CHECKSUM="${GRAFANA_PLUGIN_ARTIFACT}.md5")
	
package: check_os package_node
	@mage build:${OS}
	@yarn sign
	@mv dist/ ${GRAFANA_PLUGIN_ID}
	@zip ${GRAFANA_PLUGIN_ARTIFACT} ${GRAFANA_PLUGIN_ID} -r
	@$(MD5CMD) ${GRAFANA_PLUGIN_ARTIFACT} > ${GRAFANA_PLUGIN_ARTIFACT_CHECKSUM}
	@plugincheck ${GRAFANA_PLUGIN_ARTIFACT}
	@rm -r ${GRAFANA_PLUGIN_ID}
