
REMOTE_HOST ?= "192.168.0.130"

server:
	go run ./cmd/server/

clean-dist:
	-rm -rf dist

dist: clean-dist
	mkdir dist

dist/server-rpi: dist
	GOOS=linux GOARCH=arm64 go build -o ./dist/server-rpi ./cmd/server

remote: dist/server-rpi
	-ssh "$(USER)@$(REMOTE_HOST)" "rm /tmp/server-rpi"
	-ssh "$(USER)@$(REMOTE_HOST)" "rm /tmp/execution.log"
	./scripts/send_file.sh $(USER) $(REMOTE_HOST) "dist/server-rpi" "/tmp/server-rpi"
	./scripts/send_file.sh $(USER) $(REMOTE_HOST) "config.yaml" "/tmp/"
	ssh "$(USER)@$(REMOTE_HOST)" "cd /tmp ; /tmp/server-rpi"
	ssh "$(USER)@$(REMOTE_HOST)" "cat /tmp/execution.log"

snapshot:
	gorelease snapshot