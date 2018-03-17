package command

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/motki/core/log"
	"github.com/motki/core/model"
	"github.com/motki/core/proto/client"

	"strings"

	"github.com/motki/cli"
	"github.com/motki/cli/text"
	"github.com/peterh/liner"
)

// InventoryCommand provides an interactive manager for Inventory.
type InventoryCommand struct {
	character *model.Character
	corp      *model.Corporation
	corpID    int

	env    *cli.Prompter
	logger log.Logger
	client client.Client
}

func NewInventoryCommand(cl client.Client, p *cli.Prompter, logger log.Logger) InventoryCommand {
	var corp *model.Corporation
	var char *model.Character
	var corpID int
	char, err := cl.CharacterForRole(model.RoleLogistics)
	if err == nil {
		corp, err = cl.GetCorporation(char.CorporationID)
		if err != nil {
			logger.Warnf("command: unable to get corporation details: %s", err.Error())
		} else {
			corpID = corp.CorporationID
		}
	}
	if err != nil && err != client.ErrNotAuthenticated {
		logger.Debugf("command: unable to load auth details: %s", err.Error())
	}
	return InventoryCommand{
		char,
		corp,
		corpID,
		p,
		logger,
		cl}
}

func (c InventoryCommand) RequiresAuth() bool {
	return true
}

func (c InventoryCommand) Prefixes() []string {
	return []string{"inventory", "inv"}
}

func (c InventoryCommand) Description() string {
	if c.corp == nil {
		return "Manipulate inventory for corpID 0."
	}
	return fmt.Sprintf("Manipulate inventory for %s.", c.corp.Name)
}

func (c InventoryCommand) Handle(subcmd string, args ...string) {
	switch {
	case len(subcmd) == 0:
		c.PrintHelp()

	case subcmd == "new" || subcmd == "add" || subcmd == "create":
		c.newItem(args...)

	//case subcmd == "show" || subcmd == "view":
	//	c.showInventory(args...)

	case subcmd == "list":
		var alertsOnly bool
		if len(args) > 0 && args[0] == "alerts" {
			alertsOnly = true
		}
		c.listInventory(alertsOnly)

	case subcmd == "alerts":
		c.listInventory(true)

	//case subcmd == "edit":
	//	c.editInventory(args...)

	default:
		fmt.Printf("Unknown subcommand: %s\n", subcmd)
		c.PrintHelp()
	}
}

func (c InventoryCommand) PrintHelp() {
	colWidth := 20
	fmt.Println()
	fmt.Println(fmt.Sprintf(`Command "%s" can be used to manipulate inventory thresholds and alerting.`, text.Boldf("inventory")))
	fmt.Println()
	fmt.Println(text.WrapText(fmt.Sprintf(`When invoking a subcommand, if the optional parameter is omitted, an interactive prompt will begin to collect the necessary details.`), text.StandardTerminalWidthInChars))
	fmt.Println()
	if c.corp != nil {
		fmt.Println(text.Boldf("Character linked!"))
		fmt.Println(fmt.Sprintf("You are logged in as %s for %s.",
			text.Boldf(c.character.Name),
			text.Boldf(c.corp.Name)))
		fmt.Println()
		fmt.Println(text.WrapText(`This command will operate on inventory for your corporation. Corporation-owned assets will be tracked to display available and missing materials.`, text.StandardTerminalWidthInChars))
		fmt.Println()
	}
	fmt.Printf(`Subcommands:
  %s Create a new inventory threshold.
  %s List all inventory items.
  %s List all inventory items below their configured threshold.
  %s Edit an existing inventory threshold.
`,
		text.Boldf(text.PadTextRight("add [typeID]", colWidth)),
		text.Boldf(text.PadTextRight("list", colWidth)),
		text.Boldf(text.PadTextRight("alerts", colWidth)),
		text.Boldf(text.PadTextRight("edit [typeID]", colWidth)))
	fmt.Println()
}

// getInventoryName returns the given Inventory's name.
func (c InventoryCommand) getInventoryName(p *model.InventoryItem) string {
	t, err := c.client.GetItemType(p.TypeID)
	if err != nil {
		c.logger.Debugf("unable to get item name: %s", err.Error())
		return strconv.Itoa(p.TypeID)
	}
	return t.Name
}

// getRegionName returns the given location's name.
func (c InventoryCommand) getLocationName(locationID int) (name string, system string) {
	r, err := c.client.GetLocation(locationID)
	if err != nil {
		c.logger.Debugf("unable to get location: %s", err.Error())
		return "[Error]", ""
	}
	if r.IsCitadel() {
		name = r.Structure.Name
	} else if r.IsStation() {
		name = r.Station.Name
	}
	parts := strings.SplitN(name, " - ", 2)
	if len(parts) == 2 {
		name = parts[1]
	}
	system = r.System.Name
	return name, system
}

func (c InventoryCommand) groupByLocation(items []*model.InventoryItem, err error) (map[int][]*model.InventoryItem, error) {
	if err != nil {
		return nil, err
	}
	res := make(map[int][]*model.InventoryItem)
	for _, it := range items {
		loc, err := c.client.GetLocation(it.LocationID)
		if err != nil {
			return nil, err
		}
		id := loc.ParentID()
		res[id] = append(res[id], it)
	}
	return res, nil
}

// listInventory lists all the Inventory.
func (c InventoryCommand) listInventory(alertsOnly bool) {
	its, err := c.groupByLocation(c.client.GetInventory())
	if err != nil {
		c.logger.Debugf("unable to fetch inventory: %s", err.Error())
		fmt.Println("Error loading inventory from db, try again.")
		return
	}
	var ct int
	for _, group := range its {
		ct += len(group)
	}
	fmt.Println("Listing", ct, "inventory items in", len(its), "locations.")
	fmt.Println()
	if ct == 0 {
		fmt.Println("There are no items. Create a new inventory item with")
		fmt.Println("  inventory add")
		return
	}
	col0 := 10
	col1 := 50
	col2 := 4
	col3 := 12
	format := "%s%s%s%s\n" // TypeID, Name, Alert, Quantity

	for loc, group := range its {
		name, system := c.getLocationName(loc)
		if name != "" {
			fmt.Println(text.CenterText(name, text.StandardTerminalWidthInChars))
		}
		fmt.Println(text.CenterText(system+" System", text.StandardTerminalWidthInChars))
		fmt.Println()
		fmt.Printf(
			format,
			text.PadTextRight("Type ID", col0),
			text.PadTextRight("Name", col1),
			text.PadTextLeft("", col2),
			text.PadTextLeft("Qty", col3))
		for _, it := range group {
			belowThreshold := it.CurrentLevel < it.MinimumLevel
			if alertsOnly && !belowThreshold {
				continue
			}
			qty := strconv.Itoa(it.CurrentLevel)
			var alert string
			if belowThreshold {
				alert = "!"
			}
			fmt.Printf(
				format,
				text.PadTextRight(strconv.Itoa(it.TypeID), col0),
				text.PadTextRight(c.getInventoryName(it), col1),
				text.PadTextLeft(alert, col2),
				text.PadTextLeft(qty, col3),
			)
		}
		fmt.Println()
	}
}

func (c InventoryCommand) newItem(args ...string) {
	item, ok := c.env.PromptItemTypeDetail("Specify Item Type", strings.Join(args, " "))
	if !ok {
		return
	}
	loc, ok := c.promptLocation("Specify Location", "")
	if !ok {
		return
	}
	it, err := c.client.NewInventoryItem(item.ItemType.ID, loc.LocationID)
	if err != nil {
		c.logger.Debugf("unable to load production chain: %s", err.Error())
		fmt.Println("Error loading production chain from db, try again.")
		return
	}
	if it.MinimumLevel, ok = c.env.PromptInt("Enter minimum threshold:", nil); ok {
		if err = c.client.SaveInventoryItem(it); err != nil {
			c.logger.Debugf("unable to save inventory item: %s", err.Error())
			fmt.Println("Error saving inventory item from db, try again.")
			return
		}
	}
	c.printItemDetail(it)
}

func (c InventoryCommand) printItemDetail(it *model.InventoryItem) {
	ty, err := c.client.GetItemTypeDetail(it.TypeID)
	if err != nil {
		c.logger.Debugf("unable to load item type detail: %s", err.Error())
		fmt.Println("Error loading item details from db, try again.")
		return
	}
	loc, err := c.client.GetLocation(it.LocationID)
	if err != nil {
		c.logger.Debugf("unable to load location: %s", err.Error())
		fmt.Println("Error loading location from db, try again.")
		return
	}
	fmt.Printf("Inventory details for %s\n", ty.Name)
	fmt.Printf("  in %s\n\n", loc.String())
	fmt.Printf("%s%d\n", text.PadTextLeft("Current:", 12), it.CurrentLevel)
	fmt.Printf("%s%d\n", text.PadTextLeft("Threshold:", 12), it.MinimumLevel)
	fmt.Printf("%s%s\n", text.PadTextLeft("Updated:", 12), it.FetchedAt)
	fmt.Println()
}

// promptLocation prompts the user for a valid item type input.
//
// If the user enters an integer, it is treated as the item's Type ID.
// Otherwise, the value is used to lookup item types.
//
// This function also accepts an initial input that should be used to
// as the first round of prompt input.
func (c InventoryCommand) promptLocation(prompt string, initialInput string) (*model.Location, bool) {
	var val *model.Location
	var id int
	var err error
	valStr := initialInput
	prompt = fmt.Sprintf("%s: ", prompt)
	for {
		// This loop is ordered in such a way that it does the input validation
		// first. This allows us to specify an initial input and test that first,
		// before actually prompting the user for input.
		if valStr == "" {
			goto prompt
		}
		id, err = strconv.Atoi(valStr)
		if err != nil {
			// Input wasn't a number, presume its a search query.
			locs, err := c.client.QueryLocations(valStr)
			if err != nil || len(locs) == 0 {
				if err != nil {
					c.logger.Debugf("error querying locations: %s", err.Error())
				}
				fmt.Printf("Nothing found for \"%s\".\n", valStr)
				goto prompt
			}
			fmt.Printf("Top %d results for \"%s\":\n", len(locs), valStr)
			for _, loc := range locs {
				fmt.Printf("%s  %s\n", text.PadTextLeft(fmt.Sprintf("%d", loc.LocationID), 16), loc.String())
			}
			goto prompt
		}
		val, err = c.client.GetLocation(id)
		if err != nil {
			fmt.Printf("No location exists with ID %d.\n", id)
			c.logger.Warnf("error fetching location: %s", err.Error())
			goto prompt
		}
		// We have a valid value, break out of the loop.
		break

	prompt:
		valStr, err = c.env.Prompt(prompt)
		if err != nil {
			if err == liner.ErrPromptAborted {
				return nil, false
			}
			if err == io.EOF {
				err = errors.New("unexpected EOF")
				fmt.Println()
			}
			c.logger.Debugf("unable to read input: %s", err.Error())
			goto prompt
		}
		// Loop back around and check valStr for valid input.
		continue
	}
	return val, true
}
