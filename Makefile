run:
	go run -buildvcs=true .
build-backend:
	env GOOS=linux GOARCH=amd64 go build -buildvcs=true -o ./local/backend .

build-command:
	env GOOS=darwin GOARCH=amd64 go build -buildvcs=true -o ./local/command_amd64 command.go
	env GOOS=darwin GOARCH=arm64 go build -buildvcs=true -o ./local/command_arm64 command.go
	env GOOS=windows GOARCH=amd64 go build -buildvcs=true -o ./local/command_amd64.exe command.go

swag:
	swag init -g ./common/swagger/swagger.go -o ./common/swagger --parseDependency
