# Variables
GO = go
MAIN = main.go
TEMPLATE = templates/mono.conf
COMMON = values/common.yaml
OUTDIR = artifacts
VALUE_FILES = values/server.yaml values/macbook.yaml values/ios.yaml

.PHONY: all clean reveal

all: reveal $(OUTDIR)/server.conf $(OUTDIR)/macbook.conf $(OUTDIR)/ios.conf

reveal:
	git secret reveal $(COMMON) $(VALUE_FILES)

$(OUTDIR)/server.conf: $(MAIN) $(TEMPLATE) $(COMMON) values/server.yaml | $(OUTDIR)
	$(GO) run $(MAIN) -template $(TEMPLATE) $(COMMON) values/server.yaml > $@

$(OUTDIR)/macbook.conf: $(MAIN) $(TEMPLATE) $(COMMON) values/macbook.yaml | $(OUTDIR)
	$(GO) run $(MAIN) -template $(TEMPLATE) $(COMMON) values/macbook.yaml > $@

$(OUTDIR)/ios.conf: $(MAIN) $(TEMPLATE) $(COMMON) values/ios.yaml | $(OUTDIR)
	$(GO) run $(MAIN) -template $(TEMPLATE) $(COMMON) values/ios.yaml > $@

$(OUTDIR):
	mkdir -p $(OUTDIR)

clean:
	rm -rf $(OUTDIR)
	git secret hide
