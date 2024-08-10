APP_NAME = toolbox

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
	go build -ldflags="-s -w" -o $(APP_NAME) main.go

clean:
	@rm -rf $(APP_NAME)

.PHONY: build cleanls