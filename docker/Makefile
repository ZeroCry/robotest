TARGETS := e2e suite
DOCKER_REPO := quay.io/gravitational
VERSION ?= $(shell git describe --long --tags --always|awk -F'[.-]' '{print $$1 "." $$2 "." $$4}')
LATEST := latest

TERRAFORM_VERSION := 0.9.3
CHROMEDRIVER_VERSION := 2.29
BUILD_ARGS := --build-arg TERRAFORM_VERSION=$(TERRAFORM_VERSION) --build-arg CHROMEDRIVER_VERSION=$(CHROMEDRIVER_VERSION)

.PHONY: containers
containers: $(TARGETS)

BINARIES := $(addprefix ../build/robotest-,$(TARGETS))

$(BINARIES):
	cd .. && $(MAKE) -j build

.PHONY: $(TARGETS)
$(TARGETS): $(BINARIES)
	$(eval TEMPDIR = "$(shell mktemp -d)")
	$(eval IMAGE = $(DOCKER_REPO)/robotest-$@:$(VERSION))
	if [ -z "$(TEMPDIR)" ]; then \
	  echo "TEMPDIR is not set"; exit 1; \
	fi;
	mkdir -p $(TEMPDIR)/build
	cp -r ../assets/terraform $(TEMPDIR)
	cp -a ../build/go-junit-report $(TEMPDIR)/build/
	cp -a ../build/robotest-$@ $(TEMPDIR)/build/
	cp -r $@/* $(TEMPDIR)/
	cd $(TEMPDIR) && docker build $(BUILD_ARGS) --rm=true --pull -t $(IMAGE) .
	rm -rf $(TEMPDIR)
	@echo Built $(IMAGE)

# Publish 
 
DOCKER_IMG = $(addprefix $(DOCKER_REPO)/robotest-,$(TARGETS))

.PHONY: publish
publish: $(DOCKER_IMG)

.PHONY: $(DOCKER_IMG)
$(DOCKER_IMG): $(TARGETS)
	docker push $@:$(VERSION)
	docker tag $@:$(VERSION) $@:$(LATEST)
	docker push $@:$(LATEST)

.PHONY: what-version
what-version:
	@echo $(VERSION)