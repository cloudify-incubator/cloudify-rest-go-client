.PHONY: all
all: bin/cfy-go

PACKAGEPATH := github.com/cloudify-incubator/cloudify-rest-go-client

VERSION := `cd src/${PACKAGEPATH} && git rev-parse --short HEAD`

CLOUDPROVIDER ?= vsphere

.PHONY: reformat
reformat:
	rm -rfv pkg/*
	rm -rfv bin/*
	gofmt -w src/${PACKAGEPATH}/cloudify/rest/*.go
	gofmt -w src/${PACKAGEPATH}/cloudify/utils/*.go
	gofmt -w src/${PACKAGEPATH}/cloudify/tests/*.go
	gofmt -w src/${PACKAGEPATH}/cloudify/*.go
	gofmt -w src/${PACKAGEPATH}/cfy-go/*.go
	gofmt -w src/${PACKAGEPATH}/kubernetes/*.go

define colorecho
	@tput setaf 2
	@echo -n $1
	@tput setaf 3
	@echo $2
	@tput sgr0
endef

# cloudify rest
CLOUDIFYREST := \
	src/${PACKAGEPATH}/cloudify/rest/rest.go \
	src/${PACKAGEPATH}/cloudify/rest/types.go

pkg/linux_amd64/${PACKAGEPATH}/cloudify/rest.a: ${CLOUDIFYREST}
	$(call colorecho,"Build: ", $@)
	go build -v -i -o pkg/linux_amd64/${PACKAGEPATH}/cloudify/rest.a ${CLOUDIFYREST}

# cloudify kubernetes support
CLOUDIFYKUBERNETES := \
	src/${PACKAGEPATH}/kubernetes/mount.go \
	src/${PACKAGEPATH}/kubernetes/types.go

pkg/linux_amd64/${PACKAGEPATH}/kubernetes.a: ${CLOUDIFYKUBERNETES}
	$(call colorecho,"Build: ", $@)
	go build -v -i -o pkg/linux_amd64/${PACKAGEPATH}/kubernetes.a ${CLOUDIFYKUBERNETES}

# cloudify utils
CLOUDIFYUTILS := \
	src/${PACKAGEPATH}/cloudify/utils/utils.go

pkg/linux_amd64/${PACKAGEPATH}/cloudify/utils.a: ${CLOUDIFYUTILS}
	$(call colorecho,"Build: ", $@)
	go build -v -i -o pkg/linux_amd64/${PACKAGEPATH}/cloudify/utils.a ${CLOUDIFYUTILS}

# cloudify
CLOUDIFYCOMMON := \
	src/${PACKAGEPATH}/cloudify/scalegroup.go \
	src/${PACKAGEPATH}/cloudify/client.go \
	src/${PACKAGEPATH}/cloudify/nodes.go \
	src/${PACKAGEPATH}/cloudify/plugins.go \
	src/${PACKAGEPATH}/cloudify/instances.go \
	src/${PACKAGEPATH}/cloudify/events.go \
	src/${PACKAGEPATH}/cloudify/blueprints.go \
	src/${PACKAGEPATH}/cloudify/status.go \
	src/${PACKAGEPATH}/cloudify/executions.go \
	src/${PACKAGEPATH}/cloudify/deployments.go

pkg/linux_amd64/${PACKAGEPATH}/cloudify.a: ${CLOUDIFYCOMMON} pkg/linux_amd64/${PACKAGEPATH}/cloudify/rest.a
	$(call colorecho,"Build: ",$@)
	go build -v -i -o pkg/linux_amd64/${PACKAGEPATH}/cloudify.a ${CLOUDIFYCOMMON}

CFYGOLIBS := \
	pkg/linux_amd64/${PACKAGEPATH}/cloudify/utils.a \
	pkg/linux_amd64/${PACKAGEPATH}/kubernetes.a \
	pkg/linux_amd64/${PACKAGEPATH}/cloudify.a

# cfy-go
CFYGO := \
	src/${PACKAGEPATH}/cfy-go/blueprints.go \
	src/${PACKAGEPATH}/cfy-go/deployments.go \
	src/${PACKAGEPATH}/cfy-go/events.go \
	src/${PACKAGEPATH}/cfy-go/executions.go \
	src/${PACKAGEPATH}/cfy-go/info.go \
	src/${PACKAGEPATH}/cfy-go/instances.go \
	src/${PACKAGEPATH}/cfy-go/kubernetes.go \
	src/${PACKAGEPATH}/cfy-go/main.go \
	src/${PACKAGEPATH}/cfy-go/nodes.go \
	src/${PACKAGEPATH}/cfy-go/plugins.go \
	src/${PACKAGEPATH}/cfy-go/scaling.go

bin/cfy-go: ${CFYGO} ${CFYGOLIBS}
	$(call colorecho,"Install: ", $@)
	# delete -s -w if you want to debug
	go build -v -ldflags "-s -w -X main.versionString=${VERSION}" -o bin/cfy-go ${CFYGO}

.PHONY: test
test:
	go test -cover ./src/${PACKAGEPATH}/...
	go get github.com/golang/lint/golint
	golint ./src/${PACKAGEPATH}/...
