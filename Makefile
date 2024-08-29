# Variables
GO = go
MAIN = main.go
TEMPLATE = templates/mono.conf
COMMON = values/common.yaml
OUTDIR = artifacts

# Targets
.PHONY: all clean

all: $(OUTDIR)/server.conf $(OUTDIR)/macbook.conf $(OUTDIR)/ios.conf $(OUTDIR)/office.conf

$(OUTDIR)/server.conf: $(MAIN) $(TEMPLATE) $(COMMON) values/server.yaml | $(OUTDIR)
	$(GO) run $(MAIN) -template $(TEMPLATE) $(COMMON) values/server.yaml > $@

$(OUTDIR)/macbook.conf: $(MAIN) $(TEMPLATE) $(COMMON) values/macbook.yaml | $(OUTDIR)
	$(GO) run $(MAIN) -template $(TEMPLATE) $(COMMON) values/macbook.yaml > $@

$(OUTDIR)/ios.conf: $(MAIN) $(TEMPLATE) $(COMMON) values/ios.yaml | $(OUTDIR)
	$(GO) run $(MAIN) -template $(TEMPLATE) $(COMMON) values/ios.yaml > $@

$(OUTDIR)/office.conf: $(MAIN) $(TEMPLATE) $(COMMON) values/office.yaml | $(OUTDIR)
	$(GO) run $(MAIN) -template $(TEMPLATE) $(COMMON) values/office.yaml > $@

$(OUTDIR):
	mkdir -p $(OUTDIR)

clean:
	rm -rf $(OUTDIR)
