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