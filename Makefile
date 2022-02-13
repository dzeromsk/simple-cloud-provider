CHART=charts/simple-cloud-provider/Chart.yaml
VALUES=charts/simple-cloud-provider/values.yaml

REPOSITORY:=$(shell python3 scripts/repository.py $(VALUES))
VERSION:=$(shell python3 scripts/version.py $(CHART))

build:
	docker build -t $(REPOSITORY):$(VERSION) .

load:
	kind load docker-image $(REPOSITORY):$(VERSION)

clean:
	rm -f simple-cloud-provider-$(VERSION).tgz 
	rm -rf output/

package: output output/index.yaml

output/index.yaml: output/simple-cloud-provider-$(VERSION).tgz
	helm repo index output/

output/simple-cloud-provider-$(VERSION).tgz: simple-cloud-provider-$(VERSION).tgz
	cp -f simple-cloud-provider-$(VERSION).tgz output/

simple-cloud-provider-$(VERSION).tgz:
	helm package charts/simple-cloud-provider

output:
	mkdir -p output/