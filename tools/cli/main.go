package main

import (
	"fmt"
	"os"

	"github.com/games4l/backend/tools/cli/cmd"
	"github.com/joho/godotenv"
)

var (
	args     map[string]string
	commands []string
)

func fatalerr(msg string) {
	fmt.Println("ERROR: " + msg)
	os.Exit(1)
}

func init() {
	var err error
	args, commands, err = cmd.ParseArgs(os.Args[1:])

	if err != nil {
		fatalerr(err.Error())
	}

	envFileName, ok := args["env-file"]
	if !ok {
		envFileName = ".env"
	}

	if file, err := os.Open(envFileName); err == nil {
		defer file.Close()

		m, err := godotenv.Parse(file)
		if ok && err != nil {
			fatalerr("The env file is not valid")
		}

		if m != nil {
			for k, v := range m {
				os.Setenv(k, v)
			}
		}
	} else if ok {
		fatalerr("The env file could not be opened")
	}
}

func main() {
	baseUrl, ok := args["base-url"]
	if !ok {
		baseUrl = os.Getenv("APP_URL")
		if baseUrl == "" {
			fatalerr("The base-url arg or the APP_URL env must be provided")
		}
	}

	webhookSig, ok := args["webhook-sig"]
	if !ok {
		webhookSig = os.Getenv("WEBHOOK_SIG")
		if webhookSig == "" {
			fatalerr("The webhook-sig arg or the WEBHOOK_SIG env must be provided")
		}
	}

	if err := cmd.HandleCommand(args, commands, baseUrl, webhookSig); err != nil {
		fatalerr(err.Error())
	}
}
