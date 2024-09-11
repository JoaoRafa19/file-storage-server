.PHONY: up
up:
	@if [ "$(MODE)" = "l" ]; then \
		docker compose up; \
	else \
		docker compose up -d; \
	fi

.PHONY: restart
restart:
	docker compose restart

.PHONY: log
log:
	docker logs -f core

.PHONY: down

down:
	docker compose down


.PHONY:test

test:
	@go test -coverprofile cover.out ./...


.PHONY: local
local:
	air

BINARY_NAME=fss

# remove any binaries that are built
clean:
	rm -f ./bin/$(BINARY_NAME)*

build-debug: clean
	CGO_ENABLED=0 go build -gcflags=all="-N -l" -o bin/$(BINARY_NAME)-debug main.go