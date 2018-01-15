package command

import (
	"fmt"

	"github.com/pkg/errors"

	"strings"

	"github.com/motki/motki-cli/cli"
	"github.com/motki/motki-cli/cli/text"
)

var errCommandNotFound = errors.New("command not found")
var errExitEditor = errors.New("exit editor")

// An editor handles simple subcommand abstraction.
//
// Note that this structure is not intended to be modified after creation. This
// removes the need for any kind of synchronization.
type editor struct {
	prompter *cli.Prompter

	// Editor help text.
	help string

	// Commands registered with the editor.
	commands []string
	// Lookup map for registered commands.
	registry map[string]editorCommand
	// prompt contains the auto-generated prompt for the editor.
	prompt string

	// Widths are measured in number of characters.
	widthCol0 int // Gutter column.
	widthCol1 int // Command name + args column.
	widthCol2 int // Description column.
}

func newEditor(prompt *cli.Prompter, helpText string, commands ...editorCommand) *editor {
	const widthGutter = 2
	const widthCol0 = 7
	e := &editor{
		prompter: prompt,

		help: helpText,

		commands: []string{},
		registry: make(map[string]editorCommand),

		prompt:    "",
		widthCol0: 2,
		widthCol1: widthCol0,
		widthCol2: text.StandardTerminalWidthInChars - widthCol0 - widthGutter}
	// Resulting order looks like: [Quit command, your commands..., Help command]
	commands = append(
		append([]editorCommand{}, newEditorQuitCommand()),
		append(commands, newEditorHelpCommand(e))...)
	for _, c := range commands {
		s := c.name()
		e.registry[s] = c
		e.commands = append(e.commands, s)
	}
	e.prompt = "Specify operation [" + strings.Join(e.commands, ",") + "]"
	return e
}

func (e *editor) loop() {
	for {
		cmd, args, ok := e.prompter.PromptStringWithArgs(
			e.prompt,
			nil,
			transformStringToCaps,
			validateStringIsOneOf(e.commands))
		cmd = strings.ToUpper(cmd)
		if !ok || cmd == "Q" {
			return
		}
		err := e.handle(cmd, args...)
		if err != nil {
			if err == errExitEditor {
				return
			}
			fmt.Println("error:", err.Error())
		}
	}
}

func (e *editor) handle(cmd string, args ...string) error {
	c, ok := e.registry[cmd]
	if !ok {
		return errCommandNotFound
	}
	return c.handle(args)
}

type editorCommand interface {
	handle(args []string) error
	name() string
	printHelpText(widthCol0, widthCol1, widthCol2 int)
}

type editorCommandFunc struct {
	shortName   string
	description string
	args        []string
	handler     func(args []string) error
}

func (dcc *editorCommandFunc) handle(args []string) error {
	return dcc.handler(args)
}

func (dcc *editorCommandFunc) name() string {
	return dcc.shortName
}

func (dcc *editorCommandFunc) printHelpText(widthCol0, widthCol1, widthCol2 int) {
	arg := ""
	if len(dcc.args) > 0 {
		// TODO: support more than one argument
		arg = dcc.args[0]
	}
	fmt.Printf(
		"%s%s %s\n",
		strings.Repeat(" ", widthCol0),
		text.Boldf(text.PadTextRight(fmt.Sprintf("%s %s", dcc.shortName, arg), widthCol1)),
		text.PadTextRight(dcc.description, widthCol2))
}

func newEditorCommandFunc(shortName, description string, args []string, handler func(args []string) error) editorCommand {
	return &editorCommandFunc{shortName, description, args, handler}
}

func newEditorQuitCommand() editorCommand {
	return newEditorCommandFunc(
		"Q",
		"Quit the editor without saving.",
		[]string{},
		func(_ []string) error {
			return errExitEditor
		})
}

func newEditorHelpCommand(ed *editor) editorCommand {
	return newEditorCommandFunc("?", "Print help text.", []string{}, func(args []string) error {
		fmt.Println()
		fmt.Println(text.WrapText(ed.help, text.StandardTerminalWidthInChars))
		fmt.Println()
		fmt.Println("Available operations:")
		for _, cn := range ed.commands {
			c, ok := ed.registry[cn]
			if !ok {
				return errCommandNotFound
			}
			c.printHelpText(ed.widthCol0, ed.widthCol1, ed.widthCol2)
		}
		fmt.Println()
		return nil
	})
}
