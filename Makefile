CMDS := node_top
DEPS := $(addprefix kubectl-, $(CMDS))
BUILD_DEPS := $(addprefix build/, $(DEPS))
STATIC_FLAG=-extldflags=-static
KUBECTL_LOCATION := $(shell which kubectl)

build/kubectl-%: ./cmd/%/*
	go build -o $@ \
		-trimpath \
		-ldflags "${STATIC_FLAG} -s -w" \
		./cmd/$*/...

install: $(BUILD_DEPS)
	$(foreach dep,$(DEPS),\
		cp build/$(dep) $(KUBECTL_LOCATION)-$(subst kubectl-,,$(dep)); \
	)

uninstall:
	$(foreach dep,$(DEPS),\
		rm -f $(KUBECTL_LOCATION)-$(subst kubectl-,,$(dep)); \
	)
