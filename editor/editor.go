// Package editor is an interactive command-line editor that supports
// sub-commands with arguments.
package editor // import "github.com/motki/cli/editor"

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/motki/cli"
	"github.com/motki/cli/text"
)

var ErrCommandNotFound = errors.New("command not found")
var ErrExitEditor = errors.New("exit editor")

// An Editor handles simple subcommand abstraction.
//
// Note that this structure is not intended to be modified after creation. This
// removes the need for any kind of synchronization.
type Editor struct {
	*cli.Prompter

	// Editor help text.
	help string

	// Commands registered with the Editor.
	commands []string
	// Lookup map for registered commands.
	registry map[string]Command
	// prompt contains the auto-generated prompt for the Editor.
	prompt string

	// Widths are measured in number of characters.
	widthCol0 int // Gutter column.
	widthCol1 int // Command name + args column.
	widthCol2 int // Description column.
}

func New(prompt *cli.Prompter, helpText string, commands ...Command) *Editor {
	const widthGutter = 2
	const widthCol0 = 7
	e := &Editor{
		Prompter: prompt,

		help: helpText,

		commands: []string{},
		registry: make(map[string]Command),

		prompt:    "",
		widthCol0: 2,
		widthCol1: widthCol0,
		widthCol2: text.StandardTerminalWidthInChars - widthCol0 - widthGutter}
	// Resulting order looks like: [Quit Command, your commands..., Help Command]
	commands = append(
		append([]Command{}, newQuitCommand()),
		append(commands, newHelpCommand(e))...)
	for _, c := range commands {
		s := c.name()
		e.registry[s] = c
		e.commands = append(e.commands, s)
	}
	e.prompt = "Specify operation [" + strings.Join(e.commands, ",") + "]"
	return e
}

func (e *Editor) Loop() {
	for {
		cmd, args, ok := e.Prompter.PromptStringWithArgs(
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
			if err == ErrExitEditor {
				return
			}
			fmt.Println("error:", err.Error())
		}
	}
}

func (e *Editor) handle(cmd string, args ...string) error {
	c, ok := e.registry[cmd]
	if !ok {
		return ErrCommandNotFound
	}
	return c.handle(args)
}

type Command interface {
	handle(args []string) error
	name() string
	printHelpText(widthCol0, widthCol1, widthCol2 int)
}

type commandImpl struct {
	shortName   string
	description string
	args        []string
	handler     func(args []string) error
}

func NewCommand(shortName, description string, args []string, handler func(args []string) error) Command {
	return &commandImpl{shortName, description, args, handler}
}

func (c *commandImpl) handle(args []string) error {
	return c.handler(args)
}

func (c *commandImpl) name() string {
	return c.shortName
}

func (c *commandImpl) printHelpText(widthCol0, widthCol1, widthCol2 int) {
	arg := ""
	if len(c.args) > 0 {
		// TODO: support more than one argument
		arg = c.args[0]
	}
	fmt.Printf(
		"%s%s %s\n",
		strings.Repeat(" ", widthCol0),
		text.Boldf(text.PadTextRight(fmt.Sprintf("%s %s", c.shortName, arg), widthCol1)),
		text.PadTextRight(c.description, widthCol2))
}

func newQuitCommand() Command {
	return NewCommand(
		"Q",
		"Quit the Editor without saving.",
		[]string{},
		func(_ []string) error {
			return ErrExitEditor
		})
}

func newHelpCommand(ed *Editor) Command {
	return NewCommand("?", "Print help text.", []string{}, func(args []string) error {
		fmt.Println()
		fmt.Println(text.WrapText(ed.help, text.StandardTerminalWidthInChars))
		fmt.Println()
		fmt.Println("Available operations:")
		for _, cn := range ed.commands {
			c, ok := ed.registry[cn]
			if !ok {
				return ErrCommandNotFound
			}
			c.printHelpText(ed.widthCol0, ed.widthCol1, ed.widthCol2)
		}
		fmt.Println()
		return nil
	})
}

// validateStringIsOneOf returns a string validator.
//
// This validator ensures the value received is in the given list of valid values.
func validateStringIsOneOf(valid []string) func(string) (string, bool) {
	return func(val string) (string, bool) {
		for _, v := range valid {
			if val == v {
				return val, true
			}
		}
		return val, false
	}
}

var transformStringToCaps = func(val string) (string, bool) {
	return strings.ToUpper(val), true
}
