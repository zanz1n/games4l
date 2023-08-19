package cmd_test

import (
	"encoding/json"
	"testing"

	"github.com/games4l/backend/tools/cli/cmd"
)

var (
	mustNotParseCmd = [][]string{
		{"--"},
	}

	mustParseCmd = []string{
		"command",
		"subcommand",
		"--bool-statement1",
		"--bool-statement2",
		"--statement1",
		"value1",
		"--bool-statement3",
		"--statement2",
		"value2",
	}
	intendedCmdCommands = []string{"command", "subcommand"}
	intendedCmdArgs     = [][2]string{
		{"bool-statement1", "true"},
		{"bool-statement2", "true"},
		{"bool-statement3", "true"},
		{"statement1", "value1"},
		{"statement2", "value2"},
	}
)

func TestParseArgsCorrect(t *testing.T) {
	args, commands, err := cmd.ParseArgs(mustParseCmd)

	if err != nil {
		t.Fatal("Failed to parse prepared arguments")
	}

	buf, err := json.MarshalIndent(
		map[string]any{
			"args":     args,
			"commands": commands,
		}, "", "  ")

	if err != nil {
		t.Fatal("Json display serialization failed")
	}

	t.Log("Json payload: " + string(buf))

	if len(commands) != len(intendedCmdCommands) {
		t.Fatal("The generated commands length is not consistent with the intended one")
	} else if len(args) != len(intendedCmdArgs) {
		t.Fatal("The generated args length is not consistent with the intended one")
	}

	for i := range intendedCmdCommands {
		if commands[i] != intendedCmdCommands[i] {
			t.Fatalf("Generated command '%s' that differs from '%s'",
				commands[i], intendedCmdCommands[i],
			)
		}
	}

	t.Log("The intended commands matches the result ones")

	for _, arg := range intendedCmdArgs {
		key, value := arg[0], arg[1]

		rvalue, ok := args[key]

		if !ok || value != rvalue {
			t.Fatalf(
				"Generated arg '%s' with '%s' value that differs from '%s'",
				key, rvalue, value,
			)
		}
	}

	t.Log("The intended args matches the result ones")
}

func TestParseArgsInvalid(t *testing.T) {
	var err error

	for _, cmdline := range mustNotParseCmd {
		_, _, err = cmd.ParseArgs(cmdline)

		if err == nil {
			t.Fatalf("The invalid input '%v' didn't generate any errors", cmdline)
		} else {
			t.Logf(
				"The invalid input '%v' generate an error (as expected): %s",
				cmdline, err.Error(),
			)
		}
	}

}
