
REMOTE_HOST ?= "192.168.0.130"
USER := "eldius"

server:
	go run ./cmd/server/ probe

server-query:
	go run ./cmd/server/ probe show

clean-dist:
	-rm -rf dist

dist: clean-dist
	mkdir dist

push:
	#go tool push -host 192.168.0.42 -v periph.io/x/cmd/bmxx80
	go tool push -host $(REMOTE_HOST) -v ./cmd/agent

dist/rpi-monitor-agent: dist
	$(eval GIT_SHORT_HASH=$(shell git rev-parse --short HEAD))
	@echo "GIT_SHORT_HASH: $(GIT_SHORT_HASH)"
	$(eval BUILD_TIME=$(shell $(date+"%Y-%m-%d %H:%M:%S")))
	@echo "BUILD_TIME: $(BUILD_TIME)"
	GOOS=linux GOARCH=arm64 \
		go build \
			-o ./dist/rpi-monitor-agent \
			-ldflags="-extldflags '-static' -X 'github.com/eldius/rpi-system-monitor/internal/config.BuildDate=$(BUILD_TIME)' -X 'github.com/eldius/rpi-system-monitor/internal/config.Version=dev' -X 'github.com/eldius/rpi-system-monitor/internal/config.Commit=$(GIT_SHORT_HASH)'" \
                  ./cmd/agent

remote: push
	-ssh "$(USER)@$(REMOTE_HOST)" "rm ~/execution.log"
	./scripts/send_file.sh $(USER) $(REMOTE_HOST) "config.yaml" "~"
	ssh "$(USER)@$(REMOTE_HOST)" "~/agent probe"
	ssh "$(USER)@$(REMOTE_HOST)" "cat ~/execution.log"

remote-query: push
	-ssh "$(USER)@$(REMOTE_HOST)" "rm ~/execution.log"
	./scripts/send_file.sh $(USER) $(REMOTE_HOST) "config.yaml" "~"
	ssh "$(USER)@$(REMOTE_HOST)" "~/agent probe show"
	ssh "$(USER)@$(REMOTE_HOST)" "cat ~/execution.log"

probe-show:
	go run ./cmd/agent/ probe show

probe:
	go run ./cmd/agent/ probe

monitor:
	go run ./cmd/agent/ monitor

snapshot:
	goreleaser release --snapshot --clean

test:
	go test -cover ./...

lint:
	golangci-lint run

vulncheck:
	go tool govulncheck ./...

validate: lint vulncheck test
	@echo ""
	@echo ""
	@echo "#######################"
	@echo "# Validating finished #"
	@echo "#######################"
	@echo ""
