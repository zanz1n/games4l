build-lambdas:
	GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/lambda_auth/main github.com/games4l/cmd/lambda_auth
	GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/lambda_question/main github.com/games4l/cmd/lambda_question
	GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/lambda_telemetry/main github.com/games4l/cmd/lambda_telemetry

run-lambdas:
	go run github.com/games4l/cmd/lambda_server

test:
	go test ./... -v --race

deploy:
	rm -rf ./apps/memories/dist
	pnpm --filter memories build
	make build-lambdas
	cd infra && terraform apply
