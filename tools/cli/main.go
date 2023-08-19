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
	args, commands, err = cmd.ParseArgs(os.Args)

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

}
