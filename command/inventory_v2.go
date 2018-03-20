package command

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/motki/cli"
	"github.com/motki/cli/text"
	"github.com/motki/cli/text/banner"

	"github.com/motki/core/log"
	"github.com/motki/core/model"
	"github.com/motki/core/proto/client"
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
	display.SetLabel(text.PadTextRight(title, 14)).
		SetPlaceholder("Enter search terms")
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

// newLoadingScreen creates a new tview.Flex with the given logo text displayed prominently.
// This section is straight up lifted from tview's demo. Thanks rivo!
// https://github.com/rivo/tview/blob/761e3d72da6c156aff91324b479f5cbf2272f53f/demos/presentation/cover.go
func newLoadingScreen(logo string) *tview.Flex {
	// What's the size of the logo?
	lines := strings.Split(logo, "\n")
	logoWidth := 0
	logoHeight := len(lines)
	for _, line := range lines {
		if len(line) > logoWidth {
			logoWidth = len(line)
		}
	}
	logoBox := tview.NewTextView().
		SetTextColor(tcell.ColorBlue)
	fmt.Fprint(logoBox, logo)

	// Create a frame for the subtitle and navigation infos.
	frame := tview.NewFrame(tview.NewBox()).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Loading inventory...", true, tview.AlignCenter, tcell.ColorWhite).
		AddText("", true, tview.AlignCenter, tcell.ColorWhite).
		AddText("Please stand by.", true, tview.AlignCenter, tcell.ColorBlue)

	// Create a Flex layout that centers the logo and subtitle.
	loadScreen := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 7, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(logoBox, logoWidth, 1, true).
			AddItem(tview.NewBox(), 0, 1, false), logoHeight, 1, true).
		AddItem(frame, 0, 10, false)

	return loadScreen
}

type inventoryApp struct {
	client client.Client
	logger log.Logger

	app         *tview.Application
	layout      *tview.Flex
	locListing  *tview.Table
	itemListing *tview.Table

	corpID int
	inv    *inventory

	rowToLocID map[int]int

	currentLoc     *location
	currentItem    *model.InventoryItem
	currentItemRow int
}

func newInventoryApp(cl client.Client, logger log.Logger, corpID int) (*inventoryApp, error) {
	a := &inventoryApp{
		client:      cl,
		logger:      logger,
		app:         tview.NewApplication(),
		layout:      tview.NewFlex(),
		locListing:  tview.NewTable(),
		itemListing: tview.NewTable(),
		corpID:      corpID,
		rowToLocID:  make(map[int]int),
	}

	// Locations listing
	a.locListing.SetBorder(true).SetTitle("Select Location").SetTitleAlign(tview.AlignLeft)
	a.locListing.SetBorderPadding(0, 0, 1, 1)
	a.locListing.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			a.app.Stop()
			return nil
		} else if event.Key() == tcell.KeyCtrlN {
			a.app.SetRoot(a.newItemScreen(nil), true)
			return nil
		}
		return event
	})

	// Items table
	a.itemListing.SetBorder(true).SetTitle("Inventory").SetTitleAlign(tview.AlignLeft)
	a.itemListing.SetBorderPadding(0, 0, 1, 1)
	a.itemListing.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			a.app.SetFocus(a.locListing)
			return nil
		} else if event.Key() == tcell.KeyCtrlN {
			a.app.SetRoot(a.newItemScreen(nil), true)
			return nil
		}
		return event
	})

	// Configure layout
	a.layout.SetDirection(tview.FlexRow)
	a.layout.AddItem(a.locListing, 0, 1, true)
	a.layout.AddItem(a.itemListing, 0, 3, false)

	// Selection handlers.
	a.locListing.SetSelectable(true, false).SetSelectedFunc(func(row, col int) {
		loc, ok := a.inv.locIDIndex[a.rowToLocID[row]]
		if !ok {
			a.logger.Warnf("expected locIDIndex to contain ")
		}
		a.currentLoc = loc
		a.currentItem = nil
		a.currentItemRow = 0
		a.refresh()
	})
	a.itemListing.SetSelectable(true, false).SetSelectedFunc(func(row, col int) {
		it := a.currentLoc.Inventory[row-1]
		a.app.SetRoot(a.editItemScreen(a.currentLoc, it), true)
	})

	// Load the data concurrently.
	go func() {
		var err error
		a.inv, err = newInventory(a.client)
		if err != nil {
			a.logger.Debugf("unable to fetch inventory: %s", err.Error())
			fmt.Println("Error loading inventory from db, try again.")
			a.app.Stop()
			return
		}

		a.reset()
	}()

	return a, nil
}

func (a *inventoryApp) Run() error {
	var logo string
	if corp, err := a.client.GetCorporation(a.corpID); err != nil {
		a.logger.Warnf("error getting corporation details: %s", err.Error())
		logo = "[MOTKI]"
	} else {
		logo = banner.Sprintf("[%s]", corp.Ticker)
	}
	return a.app.SetRoot(newLoadingScreen(logo), true).Run()
}

func (a *inventoryApp) escFn(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc {
		a.reset()
		return nil
	}
	return event
}

// Create a new disabled, read-only input with the given title and val.
func (a *inventoryApp) newDisabledInput(title, val string) *tview.InputField {
	input := tview.NewInputField().SetLabel(text.PadTextLeft(title, 14)).SetText(val)
	input.SetBorderPadding(1, 1, 1, 1)
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			return event
		}
		return a.escFn(event)
	})
	input.SetLabelColor(tcell.ColorLightGrey)
	input.SetFieldBackgroundColor(tcell.ColorDarkGrey)
	return input
}

func (a *inventoryApp) newItemScreen(loc *location) *tview.Flex {
	editor := tview.NewFlex()
	editor.SetDirection(tview.FlexRow)
	editor.SetBorder(true)
	editor.SetTitle("New Inventory Item").SetTitleAlign(tview.AlignLeft)

	typeInput := newSuperSelect(a.app, "Type", nil, func(input string) (matches []*selectItem) {
		if res, err := a.client.QueryItemTypes(input); err == nil {
			for _, r := range res {
				matches = append(matches, &selectItem{strconv.Itoa(r.ID), r.Name})
			}
		}
		return matches
	})

	locInput := newSuperSelect(a.app, "Location", nil, func(input string) (matches []*selectItem) {
		if res, err := a.client.QueryLocations(input); err == nil {
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
	minInput.SetInputCapture(a.escFn)
	minInput.SetAcceptanceFunc(func(in string, last rune) bool {
		if _, err := strconv.Atoi(in); err != nil {
			return false
		}
		return true
	})

	form := tview.NewForm()
	form.AddFormItem(typeInput).AddFormItem(locInput).AddFormItem(minInput)
	form.AddButton("Save New", func() {
		typID, _ := strconv.Atoi(typeInput.selected.id)
		locID, _ := strconv.Atoi(locInput.selected.id)
		min, _ := strconv.Atoi(minInput.GetText())
		a.doSave(nil, func(_ *model.InventoryItem) *model.InventoryItem {
			it, err := a.client.NewInventoryItem(typID, locID)
			if err != nil {
				a.logger.Warnf("error saving inventory item: %s", err)
				return nil
			}
			if it.MinimumLevel != min {
				it.MinimumLevel = min
			}
			return it
		})
	})

	editor.AddItem(form, 0, 1, true)

	return editor
}

func (a *inventoryApp) editItemScreen(loc *location, it *model.InventoryItem) *tview.Flex {
	a.currentItem = it

	editor := tview.NewFlex()
	editor.SetDirection(tview.FlexRow)
	editor.SetBorder(true)
	editor.SetTitle("Edit Inventory Item").SetTitleAlign(tview.AlignLeft)

	// Disabled Location input.
	locInput := a.newDisabledInput("Location", loc.String())
	editor.AddItem(locInput, 0, 1, false)

	// Disabled item name input.
	itemInput := a.newDisabledInput("Type", a.getInventoryName(it))
	editor.AddItem(itemInput, 0, 1, false)

	// Disabled current level input.
	currInput := a.newDisabledInput("Current Level", strconv.Itoa(it.CurrentLevel))
	editor.AddItem(currInput, 0, 1, false)

	// Editable minimum level input.
	minInput := tview.NewInputField().
		SetLabel(text.PadTextLeft("Minimum Level", 14)).
		SetText(strconv.Itoa(it.MinimumLevel))
	minInput.SetInputCapture(a.escFn)
	minInput.SetAcceptanceFunc(func(in string, last rune) bool {
		if _, err := strconv.Atoi(in); err != nil {
			return false
		}
		return true
	})

	// Configure form.
	form := tview.NewForm()
	form.AddFormItem(minInput).
		AddButton("Save Changes", func() {
			val := minInput.GetText()
			v, err := strconv.Atoi(val)
			if err != nil {
				return
			}
			it.MinimumLevel = v
			a.doSave(it, nil)
		})
	form.SetCancelFunc(a.reset)

	editor.AddItem(form, 0, 2, true)

	return editor
}

func (a *inventoryApp) reset() {
	a.refresh()
	a.app.SetRoot(a.layout, true)
	if a.currentLoc != nil {
		a.app.SetFocus(a.itemListing)
	} else {
		a.app.SetFocus(a.locListing)
	}
	a.app.Draw()
}

func (a *inventoryApp) refresh() {
	// Refresh the locations listing.
	a.locListing.Clear()
	a.inv.Sort()
	a.locListing.ScrollToBeginning()
	for i, loc := range a.inv.locations {
		if loc == a.currentLoc {
			a.locListing.Select(i, 0)
		}
		name, system := loc.Name()
		a.locListing.SetCell(i, 0,
			tview.NewTableCell(fmt.Sprintf("%s - %s", system, name)).
				SetExpansion(1).SetTextColor(tcell.ColorLightGrey))
		a.rowToLocID[i] = loc.LocationID

	}
	a.itemListing.Clear()
	a.itemListing.ScrollToBeginning()

	// Next, refresh the items listing.
	loc := a.currentLoc
	if loc == nil {
		return
	}
	loc.Sort()
	a.currentItemRow = 0
	a.app.SetFocus(a.itemListing)
	a.itemListing.SetCell(0, 0,
		tview.NewTableCell("Type ID").
			SetTextColor(tcell.ColorWhite).SetSelectable(false))
	a.itemListing.SetCell(0, 1,
		tview.NewTableCell("Name").
			SetTextColor(tcell.ColorWhite).SetSelectable(false))
	a.itemListing.SetCell(0, 2,
		tview.NewTableCell("").SetTextColor(tcell.ColorWhite).SetSelectable(false))
	a.itemListing.SetCell(0, 3,
		tview.NewTableCell("Qty").
			SetAlign(tview.AlignRight).SetTextColor(tcell.ColorWhite).SetSelectable(false))
	a.itemListing.SetFixed(1, 0)
	for j, it := range loc.Inventory {
		r := j + 1
		if it == a.currentItem {
			a.itemListing.Select(r, 0)
		}
		var alert string
		alertColor := tcell.ColorLightGrey
		if belowThreshold(it) {
			alert = " !"
			alertColor = tcell.ColorRed
		}
		a.itemListing.SetCell(r, 0, tview.NewTableCell(strconv.Itoa(it.TypeID)).
			SetTextColor(tcell.ColorLightGrey))
		a.itemListing.SetCell(r, 1, tview.NewTableCell(a.getInventoryName(it)).
			SetExpansion(2).SetTextColor(tcell.ColorLightGrey))
		a.itemListing.SetCell(r, 2, tview.NewTableCell(alert).SetTextColor(alertColor))
		a.itemListing.SetCell(r, 3,
			tview.NewTableCell(strconv.Itoa(it.CurrentLevel)).SetAlign(tview.AlignRight).
				SetTextColor(tcell.ColorLightGrey))
	}
}

func (a *inventoryApp) doSave(it *model.InventoryItem, fn func(*model.InventoryItem) *model.InventoryItem) {
	// Configure the modal.
	modal := tview.NewModal()
	modal.SetTitle("Saving...")
	modal.SetText("Saving item.")
	a.app.SetRoot(modal, false).Draw()

	done := make(chan struct{})
	go func() {
		select {
		case <-done:
			return

		case <-time.After(3 * time.Second):
			modal := tview.NewModal()
			modal.SetTitle("Still saving...")
			modal.SetText("It's taking longer than expected. Please stand by.")
			a.app.SetRoot(modal, false).Draw()
		}
	}()

	// Save the item in a goroutine.
	go func() {
		defer close(done)
		if fn != nil {
			it = fn(it)
		}
		err := a.client.SaveInventoryItem(it)
		if err != nil {
			a.logger.Warnf("error saving inventory item: %s", err)
			// Error, so don't cancel the edits.
			return
		}
		// Reset back to our listings.
		if err = a.inv.Track(it); err != nil {
			a.logger.Warnf("error tracking item in local inventory: %s", err.Error())
			a.app.Stop()
			return
		}
		a.currentLoc = a.inv.locIDIndex[it.LocationID]
		a.currentItem = it
		a.reset()
	}()
}

func (c InventoryV2Command) Handle(subcmd string, args ...string) {
	a, err := newInventoryApp(c.client, c.logger, c.corpID)
	if err != nil {
		c.logger.Warnf("error initializing inventory app: %s", err.Error())
	}
	err = a.Run()
	if err != nil {
		c.logger.Warnf("error running inventory app: %s", err.Error())
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
func (a *inventoryApp) getInventoryName(p *model.InventoryItem) string {
	t, err := a.client.GetItemType(p.TypeID)
	if err != nil {
		a.logger.Debugf("unable to get item name: %s", err.Error())
		return strconv.Itoa(p.TypeID)
	}
	return t.Name
}

// Name returns the given location's name.
func (l *location) Name() (name string, system string) {
	if l.IsCitadel() {
		name = l.Structure.Name
	} else if l.IsStation() {
		name = l.Station.Name
	}
	parts := strings.SplitN(name, " - ", 2)
	if len(parts) == 2 {
		name = parts[1]
	}
	system = l.System.Name
	return name, system
}

type location struct {
	*model.Location

	Inventory []*model.InventoryItem
}

func (l *location) Sort() {
	sort.Sort(itemSlice(l.Inventory))
}

type inventory struct {
	cl client.Client

	locations []*location

	// LocationID to location.
	locIDIndex map[int]*location

	// ItemID to location.
	itemIDLocIndex map[itemID]*location
}

func newInventory(cl client.Client) (*inventory, error) {
	items, err := cl.GetInventory()
	if err != nil {
		return nil, err
	}
	inv := &inventory{
		cl:             cl,
		locations:      make([]*location, 0),
		locIDIndex:     make(map[int]*location),
		itemIDLocIndex: make(map[itemID]*location)}
	for _, it := range items {
		if err := inv.Track(it); err != nil {
			return nil, err
		}
	}
	return inv, nil
}

func (i *inventory) Sort() {
	sort.Sort(locSlice(i.locations))
}

type itemID struct {
	TypeID     int
	LocationID int
}

func getItemID(it *model.InventoryItem) itemID {
	return itemID{it.TypeID, it.LocationID}
}

func (i *inventory) Track(it *model.InventoryItem) error {
	itID := getItemID(it)
	if _, ok := i.itemIDLocIndex[itID]; ok {
		// Already tracked
		return nil
	}
	loc, ok := i.locIDIndex[it.LocationID]
	if !ok {
		dl, err := i.cl.GetLocation(it.LocationID)
		if err != nil {
			return err
		}
		loc = &location{Location: dl}
		i.locations = append(i.locations, loc)
		i.locIDIndex[it.LocationID] = loc
	}
	i.itemIDLocIndex[itID] = loc
	loc.Inventory = append(loc.Inventory, it)
	return nil
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

type locSlice []*location

func (s locSlice) Less(i, j int) bool {
	return s[i].System.Name < s[j].System.Name
}

// Len is the number of elements in the collection.
func (s locSlice) Len() int {
	return len(s)
}

// Swap swaps the elements with indexes i and j.
func (s locSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
