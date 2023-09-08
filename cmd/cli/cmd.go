package main

import "errors"

func HandleCommand(
	args map[string]string,
	cmds []string,
	baseUrl string,
	sigKey string,
) error {
	if len(cmds) != 2 {
		return errors.New("invalid commands length")
	}

	var err error

	switch cmds[0] {
	case "user":
		switch cmds[1] {
		case "create":
			if err = HandleUserCreate(args, baseUrl, sigKey); err != nil {
				return err
			}
		default:
			return errors.New("invalid subcommand " + cmds[1])
		}
	case "question":
		switch cmds[1] {
		case "convert":
			if err = HandleQuestionConvert(args); err != nil {
				return err
			}
		}
	default:
		return errors.New("invalid command " + cmds[1])
	}

	return nil
}
