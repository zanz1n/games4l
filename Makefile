build-lambdas:
	cd ./services/auth_lambda && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/main .
	cd ./services/telemetry_lambda && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/main .
	cd ./services/question_lambda && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/main .

run-lambdas:
	go run github.com/games4l/backend/tools/lambda_dev

deploy:
	make build-lambdas
	cd infra && terraform apply
