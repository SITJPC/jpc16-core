RUN_FLAGS = -buildvcs=true
BUILD_FLAGS = -tags prod -buildvcs=true
BUILD_COMMAND = go build $(BUILD_FLAGS)

run:
	go run $(RUN_FLAGS) .

run-prod:
	go run $(BUILD_FLAGS) .

build-backend:
	env GOOS=linux GOARCH=amd64 $(BUILD_COMMAND) -o ./local/backend .

build-command:
	env GOOS=darwin GOARCH=amd64 $(BUILD_COMMAND) -o ./local/command_amd64 command.go
	env GOOS=darwin GOARCH=arm64 $(BUILD_COMMAND) -o ./local/command_arm64 command.go
	env GOOS=windows GOARCH=amd64 $(BUILD_COMMAND) -o ./local/command_amd64.exe command.go

swag:
	swag init -g ./common/swagger/swagger.go -o ./common/swagger --parseDependency
