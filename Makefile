
REMOTE_HOST ?= "192.168.0.130"

server:
	go run ./cmd/server/ probe

server-query:
	go run ./cmd/server/ probe show

clean-dist:
	-rm -rf dist

dist: clean-dist
	mkdir dist

dist/rpi-server: dist
	$(eval GIT_SHORT_HASH=$(shell git rev-parse --short HEAD))
	@echo "GIT_SHORT_HASH: $(GIT_SHORT_HASH)"
	$(eval BUILD_TIME=$(shell $(date+"%Y-%m-%d %H:%M:%S")))
	@echo "BUILD_TIME: $(BUILD_TIME)"
	GOOS=linux GOARCH=arm64 \
		go build \
			-o ./dist/rpi-server \
			-ldflags="-extldflags '-static' -X 'github.com/eldius/rpi-system-monitor/internal/config.BuildDate=$(BUILD_TIME)' -X 'github.com/eldius/rpi-system-monitor/internal/config.Version=dev' -X 'github.com/eldius/rpi-system-monitor/internal/config.Commit=$(GIT_SHORT_HASH)'" \
                  ./cmd/server

remote: dist/rpi-server
	-ssh "$(USER)@$(REMOTE_HOST)" "rm /tmp/rpi-server"
	-ssh "$(USER)@$(REMOTE_HOST)" "rm /tmp/execution.log"
	./scripts/send_file.sh $(USER) $(REMOTE_HOST) "dist/rpi-server" "/tmp/rpi-server"
	./scripts/send_file.sh $(USER) $(REMOTE_HOST) "config.yaml" "/tmp/"
	ssh "$(USER)@$(REMOTE_HOST)" "cd /tmp ; /tmp/rpi-server probe"
	ssh "$(USER)@$(REMOTE_HOST)" "cat /tmp/execution.log"

remote-query: dist/rpi-server
	-ssh "$(USER)@$(REMOTE_HOST)" "rm /tmp/rpi-server"
	-ssh "$(USER)@$(REMOTE_HOST)" "rm /tmp/execution.log"
	./scripts/send_file.sh $(USER) $(REMOTE_HOST) "dist/rpi-server" "/tmp/rpi-server"
	./scripts/send_file.sh $(USER) $(REMOTE_HOST) "config.yaml" "/tmp/"
	ssh "$(USER)@$(REMOTE_HOST)" "cd /tmp ; /tmp/rpi-server probe show"
	ssh "$(USER)@$(REMOTE_HOST)" "cat /tmp/execution.log"

snapshot:
	goreleaser release --snapshot --clean
