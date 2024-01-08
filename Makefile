DEV_FLAGS = -buildvcs=true
PROD_FLAGS = -tags prod -buildvcs=true
BUILD_PROD = go build $(PROD_FLAGS)

run:
	go run $(DEV_FLAGS) .

run-prod:
	go run $(PROD_FLAGS) .

build-backend:
	env GOOS=linux GOARCH=amd64 $(BUILD_PROD) -o ./local/backend .

build-command:
	env GOOS=darwin GOARCH=amd64 $(BUILD_PROD) -o ./local/command_amd64 command.go
	env GOOS=darwin GOARCH=arm64 $(BUILD_PROD) -o ./local/command_arm64 command.go
	env GOOS=windows GOARCH=amd64 $(BUILD_PROD) -o ./local/command_amd64.exe command.go

swag:
	swag init -g ./common/swagger/swagger.go -o ./common/swagger --parseDependency
