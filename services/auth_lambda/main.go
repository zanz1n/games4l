package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/services/auth_lambda/src"
)

func main() {
	utils.DefaultErrorList.Apply(utils.PtBtMessages)
	os.Setenv("NO_COLOR", "1")
	lambda.Start(src.Handler)
}
