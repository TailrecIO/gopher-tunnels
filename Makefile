build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/server/register server/register.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/server/respond server/respond.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/server/webhook server/webhook.go
	env GOOS=darwin go build -ldflags="-s -w" -o bin/client/config client/client_config.go
	env GOOS=darwin go build -ldflags="-s -w" -o bin/client/gopher client/gopher.go