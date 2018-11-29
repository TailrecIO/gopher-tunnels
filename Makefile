build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/server/register server/register.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/server/respond server/respond.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/server/webhook server/webhook.go
	go build -ldflags="-s -w" -o bin/tools/config tools/client_config.go
	go build -ldflags="-s -w" -o bin/tools/http_echo tools/http_echo.go
	go build -ldflags="-s -w" -o bin/client/gopher client/gopher.go