work-tidy:
	cd ./libs/auth && go mod tidy -e
	cd ./libs/logger && go mod tidy -e
	cd ./libs/telemetry && go mod tidy -e
	cd ./libs/utils && go mod tidy -e
	cd ./libs/user && go mod tidy -e
	cd ./libs/question && go mod tidy -e

	cd ./services/telemetry_lambda && go mod tidy -e
	cd ./services/auth_lambda && go mod tidy -e
	cd ./services/question_lambda && go mod tidy -e

	cd ./tools/lambda-dev && go mod tidy -e
	cd ./tools/cli && go mod tidy -e

build-lambdas:
	cd ./services/auth_lambda && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/main .
	cd ./services/telemetry_lambda && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/main .
	cd ./services/question_lambda && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/main .

run-lambdas:
	go run github.com/games4l/backend/tools/lambda_dev

cli:
	go run github.com/games4l/backend/tools/cli

test:
	go test ./tools/cli/... -v --race

deploy:
	make build-lambdas
	cd infra && terraform apply
