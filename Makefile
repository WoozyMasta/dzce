GO          ?= go
LINTER      ?= golangci-lint
ALIGNER     ?= betteralign
VULNCHECK   ?= govulncheck
BENCHSTAT   ?= benchstat
LINTKIT     ?= lintkit
SCHEMADOC   ?= schemadoc
BENCH_COUNT ?= 6
BENCH_REF   ?= bench_baseline.txt

.PHONY: check ci

check: verify tidy fmt vulncheck vet lint-fix align-fix test test-race generate
ci: download tools-ci verify tidy-check fmt-check generate-check vulncheck vet lint align test

.PHONY: test test-race

test:
	$(GO) test ./...

test-race:
	$(GO) test -race ./...

.PHONY: bench bench-fast bench-reset

bench:
	@tmp=$$(mktemp); \
	$(GO) test ./... -run=^$$ -bench 'Benchmark' -benchmem -count=$(BENCH_COUNT) | tee "$$tmp"; \
	if [ -f "$(BENCH_REF)" ]; then \
		$(BENCHSTAT) "$(BENCH_REF)" "$$tmp"; \
	else \
		cp "$$tmp" "$(BENCH_REF)" && echo "Baseline saved to $(BENCH_REF)"; \
	fi; \
	rm -f "$$tmp"

bench-fast:
	$(GO) test ./... -run=^$$ -bench 'Benchmark' -benchmem

bench-reset:
	rm -f "$(BENCH_REF)"

.PHONY: download verify vet tidy tidy-check fmt fmt-check vulncheck lint lint-fix align align-fix

download:
	$(GO) mod download

verify:
	$(GO) mod verify

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

tidy-check:
	@$(GO) mod tidy
	@git diff --stat --exit-code -- go.mod go.sum || ( \
		echo "go mod tidy: repository is not tidy"; \
		exit 1; \
	)

fmt:
	gofmt -w .

fmt-check:
	@files="$$(gofmt -l .)"; \
	if [ -n "$$files" ]; then \
		echo "$$files"; \
		echo "gofmt: files need formatting"; \
		exit 1; \
	fi

vulncheck:
	$(VULNCHECK) ./...

lint:
	$(LINTER) run ./...

lint-fix:
	$(LINTER) run --fix ./...

align:
	$(ALIGNER) ./...

align-fix:
	-$(ALIGNER) -apply ./...
	$(ALIGNER) ./...

.PHONY: tools tools-ci tool-golangci-lint tool-betteralign tool-benchstat

tools: tool-golangci-lint tool-betteralign tool-govulncheck tool-benchstat tool-lintkit tool-schemadoc
tools-ci: tool-golangci-lint tool-govulncheck tool-betteralign tool-lintkit tool-schemadoc

tool-golangci-lint:
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

tool-betteralign:
	$(GO) install github.com/dkorunic/betteralign/cmd/betteralign@latest

tool-govulncheck:
	$(GO) install golang.org/x/vuln/cmd/govulncheck@latest

tool-benchstat:
	$(GO) install golang.org/x/perf/cmd/benchstat@latest

tool-lintkit:
	$(GO) install github.com/woozymasta/lintkit/cmd/lintkit@latest

tool-schemadoc:
	$(GO) install github.com/woozymasta/schemadoc/cmd/schemadoc@latest

.PHONY: generate generate-check schema-gen schema-check lint-gen lint-check

generate: schema-gen lint-gen
generate-check: schema-check lint-check

schema-gen:
	$(SCHEMADOC) build schemadoc.build.yaml

schema-check:
	$(SCHEMADOC) build schemadoc.check.yaml

lint-gen:
	$(LINTKIT) snapshot -f yaml rules.yaml
	$(LINTKIT) doc -t table -w 76 rules.yaml docs/RULES.md

lint-check:
	$(LINTKIT) snapshot -cf yaml rules.yaml
	$(LINTKIT) doc -ct table -w 76 rules.yaml docs/RULES.md

.PHONY: release-notes

release-notes:
	@awk '\
	/^<!--/,/^-->/ { next } \
	/^## \[[0-9]+\.[0-9]+\.[0-9]+\]/ { if (found) exit; found=1; next } \
	found { \
		if (/^## \[/) { exit } \
		if (/^$$/) { flush(); print; next } \
		if (/^\* / || /^- /) { flush(); buf=$$0; next } \
		if (/^###/ || /^\[/) { flush(); print; next } \
		sub(/^[ \t]+/, ""); sub(/[ \t]+$$/, ""); \
		if (buf != "") { buf = buf " " $$0 } else { buf = $$0 } \
		next \
	} \
	function flush() { if (buf != "") { print buf; buf = "" } } \
	END { flush() } \
	' CHANGELOG.md
