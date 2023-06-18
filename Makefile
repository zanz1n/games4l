build-lambdas:
	cd ./services/auth_lambda && make build
	cd ./services/telemetry_lambda && make build
	cd ./services/question_lambda && make build

deploy:
	make build-lambdas
	cd infra && terraform apply
