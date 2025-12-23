BINS=$(notdir $(wildcard cmd/*))

build: $(BINS)

$(BINS):
	go build ./cmd/$@

.PHONY: $(BINS)
