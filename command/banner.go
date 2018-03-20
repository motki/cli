package command

import (
	"fmt"
	"strings"

	"path/filepath"

	"github.com/motki/cli"
	"github.com/motki/cli/text"
	"github.com/motki/cli/text/banner"
)

// BannerCommand provides item type lookup and display.
type BannerCommand struct {
	env  *cli.Prompter
	font *banner.Font
}

func NewBannerCommand(prompter *cli.Prompter) *BannerCommand {
	return &BannerCommand{prompter, banner.DefaultFont}
}

func (c *BannerCommand) Prefixes() []string {
	return []string{"banner", "b"}
}

func (c *BannerCommand) Description() string {
	return "Print ASCII art banners."
}

func (c *BannerCommand) Handle(subcmd string, args ...string) {
	var defVal string
	switch {
	case subcmd == "font":
		if len(args) < 1 {
			fmt.Println("you must specify a file to load.")
			return
		}
		fmt.Println("Loading", filepath.Base(args[0]))
		f, err := banner.NewFont(args[0])
		if err != nil {
			fmt.Printf("error loading font from file: %s\n", err.Error())
		}
		c.font = f
		fmt.Println("Loaded! Font size:", f.Size())

	case subcmd == "show":
		b := banner.New(c.font, strings.Join(args, " "))
		fmt.Print(b.String())
		return

	case subcmd != "":
		defVal = strings.Join(append([]string{subcmd}, args...), " ")

	default:
	}
	str, ok := c.env.PromptString("Enter text", &defVal)
	if !ok {
		return
	}
	fmt.Print(banner.New(c.font, str).String())
}

func (c *BannerCommand) PrintHelp() {
	colWidth := 20
	fmt.Println()
	fmt.Println(text.WrapText(fmt.Sprintf(`Command "%s" can be used to create ASCII art banners with arbitrary text.`, text.Boldf("banner")), text.StandardTerminalWidthInChars))
	fmt.Println()
	fmt.Printf(`Subcommands:
  %s Print the given text using the current font.
  %s Load the font data located at the given path.

Run banner with no parameters and you will be prompted for text to print.
`,
		text.Boldf(text.PadTextRight("show [text]", colWidth)),
		text.Boldf(text.PadTextRight("font <path>", colWidth)))
	fmt.Println()
}
