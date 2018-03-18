package command

import (
	"fmt"
	"strconv"

	"github.com/motki/core/log"
	"github.com/motki/core/model"
	"github.com/motki/core/proto/client"

	"strings"

	"sync/atomic"

	"sort"

	"github.com/gdamore/tcell"
	"github.com/motki/cli"
	"github.com/motki/cli/text"
	"github.com/rivo/tview"
)

// InventoryV2Command provides an interactive manager for Inventory.
type InventoryV2Command struct {
	character *model.Character
	corp      *model.Corporation
	corpID    int

	env    *cli.Prompter
	logger log.Logger
	client client.Client
}

func NewInventoryV2Command(cl client.Client, p *cli.Prompter, logger log.Logger) InventoryV2Command {
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
	return InventoryV2Command{
		char,
		corp,
		corpID,
		p,
		logger,
		cl}
}

func (c InventoryV2Command) RequiresAuth() bool {
	return true
}

func (c InventoryV2Command) Prefixes() []string {
	return []string{"inventory2", "inv2"}
}

func (c InventoryV2Command) Description() string {
	if c.corp == nil {
		return "Requires authentication."
	}
	return fmt.Sprintf("EXPERIMENTAL: Interactive inventory management.")
}

type selectItem struct {
	id    string
	value string
}

type searchFn func(input string) (matches []*selectItem)

type superSelect struct {
	tview.FormItem

	display   *tview.InputField
	selection *tview.DropDown

	selected *selectItem

	search searchFn
}

func newSuperSelect(app *tview.Application, title string, val *selectItem, fn searchFn) *superSelect {
	display := tview.NewInputField()
	display.SetLabel(text.PadTextRight(title, 14))
	selection := tview.NewDropDown()
	selection.SetLabel(text.PadTextRight(title, 14))
	sel := &superSelect{
		FormItem:  display,
		display:   display,
		selection: selection,
		search:    fn,
		selected:  val,
	}
	display.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			v := display.GetText()
			res := sel.search(v)
			var options []string
			for _, r := range res {
				options = append(options, r.value)
			}
			sel.selection.SetOptions(options, func(text string, index int) {
				sel.selected = res[index]
				sel.display.SetText(sel.selected.value)
				sel.FormItem = sel.display
				app.SetFocus(sel)
			})
			sel.FormItem = sel.selection
			app.SetFocus(sel)
			return nil
		}
		return event
	})
	selection.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			sel.FormItem = sel.display
			app.SetFocus(sel)
			return nil
		}
		return event
	})
	return sel
}

func (s *superSelect) SetFormAttributes(label string, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
	s.display.SetFormAttributes(label, labelColor, bgColor, fieldTextColor, fieldBgColor)
	s.selection.SetFormAttributes(label, labelColor, bgColor, fieldTextColor, fieldBgColor)
	return s
}

func (c InventoryV2Command) Handle(subcmd string, args ...string) {
	its, err := c.groupByLocation(c.client.GetInventory())
	if err != nil {
		c.logger.Debugf("unable to fetch inventory: %s", err.Error())
		fmt.Println("Error loading inventory from db, try again.")
		return
	}
	app := tview.NewApplication()

	var newItem func(*model.Location) *tview.Flex
	var editItem func(*model.Location, *model.InventoryItem) *tview.Flex

	var refreshLocations func()

	// Locations listing
	locations := tview.NewTable()
	locations.SetBorder(true).SetTitle("Select Location").SetTitleAlign(tview.AlignLeft)
	locations.SetBorderPadding(0, 0, 1, 1)
	locations.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.Stop()
			return nil
		} else if event.Key() == tcell.KeyCtrlN {
			app.SetRoot(newItem(nil), true)
			return nil
		}
		return event
	})

	// Items table
	items := tview.NewTable()
	items.SetBorder(true).SetTitle("Inventory").SetTitleAlign(tview.AlignLeft)
	items.SetBorderPadding(0, 0, 1, 1)
	items.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.SetFocus(locations)
			return nil
		} else if event.Key() == tcell.KeyCtrlN {
			app.SetRoot(newItem(nil), true)
			return nil
		}
		return event
	})

	// Configure layout
	layout := tview.NewFlex()
	layout.SetDirection(tview.FlexRow)
	layout.AddItem(locations, 0, 1, true)
	layout.AddItem(items, 0, 3, false)

	// Current location selected (the row index)
	var currLoc int64
	// Map of row in location list to slice of items in items table.
	itemMap := make(map[int][]*model.InventoryItem)
	// Map of row in location list to the location it represents.
	locMap := make(map[int]*model.Location)
	// Map of location ID to row in location list.
	rowLocMap := make(map[int]int)

	addItem := func(it *model.InventoryItem) {
		its[it.LocationID] = append(its[it.LocationID], it)
		refreshLocations()
		atomic.StoreInt64(&currLoc, int64(rowLocMap[it.LocationID]))
	}

	refreshLocations = func() {
		locations.Clear()
		var i int
		for loc, items := range its {
			rowLocMap[loc] = i
			locMap[i], err = c.client.GetLocation(loc)
			if err != nil {
				//?
				c.logger.Warnf("unable to get location: %s", err)
				continue
			}
			name, system := c.getLocationName(loc)
			func(i int, items []*model.InventoryItem) {
				locations.SetCell(i, 0, tview.NewTableCell(fmt.Sprintf("%s - %s", name, system)).
					SetExpansion(1).SetTextColor(tcell.ColorLightGrey))
				itemMap[i] = items
			}(i, items)
			i++
		}
		locations.ScrollToBeginning()
	}

	selectLocation := func(row, col int) {
		atomic.StoreInt64(&currLoc, int64(row))
		its := itemMap[row]
		sort.Sort(itemSlice(its))
		items.Clear()
		app.SetFocus(items)
		items.SetCell(0, 0, tview.NewTableCell("Type ID").SetTextColor(tcell.ColorWhite).SetSelectable(false))
		items.SetCell(0, 1, tview.NewTableCell("Name").SetTextColor(tcell.ColorWhite).SetSelectable(false))
		items.SetCell(0, 2, tview.NewTableCell("").SetTextColor(tcell.ColorWhite).SetSelectable(false))
		items.SetCell(0, 3,
			tview.NewTableCell("Qty").SetAlign(tview.AlignRight).SetTextColor(tcell.ColorWhite).SetSelectable(false))
		items.SetFixed(1, 0)
		for j, it := range its {
			r := j + 1
			belowThreshold := it.CurrentLevel < it.MinimumLevel
			var alert string
			alertColor := tcell.ColorLightGrey
			if belowThreshold {
				alert = " !"
				alertColor = tcell.ColorRed
			}
			items.SetCell(r, 0, tview.NewTableCell(strconv.Itoa(it.TypeID)).SetTextColor(tcell.ColorLightGrey))
			items.SetCell(r, 1, tview.NewTableCell(c.getInventoryName(it)).SetExpansion(2).SetTextColor(tcell.ColorLightGrey))
			items.SetCell(r, 2, tview.NewTableCell(alert).SetTextColor(alertColor))
			items.SetCell(r, 3,
				tview.NewTableCell(strconv.Itoa(it.CurrentLevel)).SetAlign(tview.AlignRight).SetTextColor(tcell.ColorLightGrey))
		}
		items.ScrollToBeginning()
	}

	// Helpers.
	cancelFn := func() {
		selectLocation(int(atomic.LoadInt64(&currLoc)), 0)
		app.SetRoot(layout, true)
		app.SetFocus(items)
	}
	readOnlyFn := func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			return event
		} else if event.Key() == tcell.KeyEsc {
			cancelFn()
		}
		return nil
	}
	escFn := func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			cancelFn()
			return nil
		}
		return event
	}

	// Create a new disabled, read-only input with the given title and val.
	newDisabledInput := func(title, val string) *tview.InputField {
		input := tview.NewInputField().SetLabel(text.PadTextLeft(title, 14)).SetText(val)
		input.SetBorderPadding(1, 1, 1, 1)
		input.SetInputCapture(readOnlyFn)
		input.SetLabelColor(tcell.ColorLightGrey)
		input.SetFieldBackgroundColor(tcell.ColorDarkGrey)
		return input
	}

	newItem = func(loc *model.Location) *tview.Flex {
		editor := tview.NewFlex()
		editor.SetDirection(tview.FlexRow)
		editor.SetBorder(true)
		editor.SetTitle("New Inventory Item").SetTitleAlign(tview.AlignLeft)

		typeInput := newSuperSelect(app, "Type", nil, func(input string) (matches []*selectItem) {
			if res, err := c.client.QueryItemTypes(input); err == nil {
				for _, r := range res {
					matches = append(matches, &selectItem{strconv.Itoa(r.ID), r.Name})
				}
			}
			return matches
		})

		locInput := newSuperSelect(app, "Location", nil, func(input string) (matches []*selectItem) {
			if res, err := c.client.QueryLocations(input); err == nil {
				for _, r := range res {
					matches = append(matches, &selectItem{strconv.Itoa(r.LocationID), r.String()})
				}
			}
			return matches
		})

		// Editable minimum level input.
		minInput := tview.NewInputField().
			SetLabel(text.PadTextLeft("Minimum Level", 14)).
			SetText("0")
		minInput.SetInputCapture(escFn)
		minInput.SetAcceptanceFunc(func(in string, last rune) bool {
			if _, err := strconv.Atoi(in); err != nil {
				return false
			}
			return true
		})

		// Closure to be called to save changes.
		saveFn := func() {
			typID, _ := strconv.Atoi(typeInput.selected.id)
			locID, _ := strconv.Atoi(locInput.selected.id)
			min, _ := strconv.Atoi(minInput.GetText())
			cancel := true
			// Suspend the main application.
			app.Suspend(func() {
				// Create a new application to display a modal.
				mApp := tview.NewApplication()
				// Configure the modal.
				modal := tview.NewModal()
				modal.SetTitle("Saving...")
				modal.SetText("Saving item...")

				// Save the item in a goroutine.
				go func() {
					it, err := c.client.NewInventoryItem(typID, locID)
					if err != nil {
						c.logger.Warnf("error saving inventory item: %s", err)
						// Error, so don't cancel the edits.
						cancel = false
					} else if it.MinimumLevel != min {
						it.MinimumLevel = min
						err = c.client.SaveInventoryItem(it)
						if err != nil {
							c.logger.Warnf("error saving inventory item: %s", err)
							// Error, so don't cancel the edits.
							cancel = false
						}
					}
					if err == nil {
						addItem(it)
					}
					// Stop the modal application when we are done saving.
					mApp.Stop()
				}()
				mApp.SetRoot(modal, false).Run()
			})
			if cancel {
				// Exit this editor and return to the main layout.
				cancelFn()
			}
		}

		form := tview.NewForm()
		form.AddFormItem(typeInput).AddFormItem(locInput).AddFormItem(minInput)
		form.AddButton("Save New", saveFn)

		editor.AddItem(form, 0, 1, true)

		return editor
	}

	// editItem launches a "subcommand" to edit the given inventory item.
	editItem = func(loc *model.Location, it *model.InventoryItem) *tview.Flex {
		editor := tview.NewFlex()
		editor.SetDirection(tview.FlexRow)
		editor.SetBorder(true)
		editor.SetTitle("Edit Inventory Item").SetTitleAlign(tview.AlignLeft)

		// Disabled Location input.
		locInput := newDisabledInput("Location", loc.String())
		editor.AddItem(locInput, 0, 1, false)

		// Disabled item name input.
		itemInput := newDisabledInput("Type", c.getInventoryName(it))
		editor.AddItem(itemInput, 0, 1, false)

		// Disabled current level input.
		currInput := newDisabledInput("Current Level", strconv.Itoa(it.CurrentLevel))
		editor.AddItem(currInput, 0, 1, false)

		// Editable minimum level input.
		minInput := tview.NewInputField().
			SetLabel(text.PadTextLeft("Minimum Level", 14)).
			SetText(strconv.Itoa(it.MinimumLevel))
		minInput.SetInputCapture(escFn)
		minInput.SetAcceptanceFunc(func(in string, last rune) bool {
			if _, err := strconv.Atoi(in); err != nil {
				return false
			}
			return true
		})

		// Closure to be called to save changes.
		saveFn := func() {
			val := minInput.GetText()
			v, err := strconv.Atoi(val)
			if err != nil {
				return
			}
			it.MinimumLevel = v
			cancel := true
			// Suspend the main application.
			app.Suspend(func() {
				// Create a new application to display a modal.
				mApp := tview.NewApplication()
				// Configure the modal.
				modal := tview.NewModal()
				modal.SetTitle("Saving...")
				modal.SetText("Saving item...")

				// Save the item in a goroutine.
				go func() {
					err := c.client.SaveInventoryItem(it)
					if err != nil {
						c.logger.Warnf("error saving inventory item: %s", err)
						// Error, so don't cancel the edits.
						cancel = false
					}
					// Stop the modal application when we are done saving.
					mApp.Stop()
				}()
				mApp.SetRoot(modal, false).Run()
			})
			if cancel {
				// Exit this editor and return to the main layout.
				cancelFn()
			}
		}

		// Configure form.
		form := tview.NewForm()
		form.AddFormItem(minInput).
			AddButton("Save Changes", saveFn)
		form.SetCancelFunc(cancelFn)

		editor.AddItem(form, 0, 2, true)

		return editor
	}

	// Selection handlers.
	locations.SetSelectable(true, false).SetSelectedFunc(selectLocation)
	items.SetSelectable(true, false).SetSelectedFunc(func(row, col int) {
		pr := int(atomic.LoadInt64(&currLoc))
		loc := locMap[pr]
		it := itemMap[pr][row-1]
		app.SetRoot(editItem(loc, it), true)
	})

	// Populate the locations table.
	refreshLocations()

	// Start the command.
	if err := app.SetRoot(layout, true).SetFocus(layout).Run(); err != nil {
		panic(err)
	}
}

func (c InventoryV2Command) PrintHelp() {
	fmt.Println()
	fmt.Println(text.WrapText(
		fmt.Sprintf(`Command "%s" is an interactive editor used for manipulating inventory thresholds and alerting.`, text.Boldf("inventory2")),
		text.StandardTerminalWidthInChars))
	fmt.Println()
	fmt.Println(text.Boldf("EXPERIMENTAL"))
	fmt.Println(text.WrapText("This command is experimental and definitely still has issues left to be worked out. You have been warned.", text.StandardTerminalWidthInChars))
	fmt.Println()
	if c.corp != nil {
		fmt.Println(text.Boldf("Character linked!"))
		fmt.Println(fmt.Sprintf("You are logged in as %s for %s.",
			text.Boldf(c.character.Name),
			text.Boldf(c.corp.Name)))
		fmt.Println()
		fmt.Println(text.WrapText(`This command will operate on inventory for your corporation. Corporation-owned structures and assets are used when calculating inventory levels. If an inventory level goes below a configurable threshold, it is made to stand out when viewing the list.`, text.StandardTerminalWidthInChars))
	}
	fmt.Println()
}

// getInventoryName returns the given Inventory's name.
func (c InventoryV2Command) getInventoryName(p *model.InventoryItem) string {
	t, err := c.client.GetItemType(p.TypeID)
	if err != nil {
		c.logger.Debugf("unable to get item name: %s", err.Error())
		return strconv.Itoa(p.TypeID)
	}
	return t.Name
}

// getRegionName returns the given location's name.
func (c InventoryV2Command) getLocationName(locationID int) (name string, system string) {
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

func (c InventoryV2Command) groupByLocation(items []*model.InventoryItem, err error) (map[int][]*model.InventoryItem, error) {
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

type itemSlice []*model.InventoryItem

func belowThreshold(it *model.InventoryItem) bool {
	return it.CurrentLevel < it.MinimumLevel
}

func (s itemSlice) Less(i, j int) bool {
	if belowThreshold(s[i]) {
		if belowThreshold(s[j]) {
			return s[i].TypeID < s[j].TypeID
		}
		return true
	}
	if belowThreshold(s[j]) {
		return false
	}
	return s[i].TypeID < s[j].TypeID
}

// Len is the number of elements in the collection.
func (s itemSlice) Len() int {
	return len(s)
}

// Swap swaps the elements with indexes i and j.
func (s itemSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
