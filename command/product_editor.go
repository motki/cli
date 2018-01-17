package command

import (
	"fmt"
	"strconv"

	"github.com/motki/core/model"

	"github.com/motki/cli/editor"
	"github.com/motki/cli/text"
)

func (c ProductCommand) newProductEditor(p *model.Product) *editor.Editor {
	lineIndex := c.getProductLineIndex(p)
	shownLineNumberHint := false
	var validLineNumber = func(val int) (int, bool) {
		_, ok := lineIndex[val]
		if !ok {
			fmt.Printf("Invalid line number %d.\n", val)
		}
		return val, ok
	}
	var promptLineNumber = func(prompt string, initVal string) (*model.Product, bool) {
		if v, err := strconv.Atoi(initVal); err == nil {
			if line, ok := validLineNumber(v); ok {
				return lineIndex[line], true
			}
		}
		if !shownLineNumberHint && initVal != "0" {
			fmt.Println(text.WrapText(fmt.Sprintf("Hint: line 0 is the main item, %s.\n", c.getProductName(p)), text.StandardTerminalWidthInChars))
			shownLineNumberHint = true
		}
		line, ok := c.env.PromptInt(prompt, nil, validLineNumber)
		if !ok {
			return nil, false
		}
		// Presume the line exists since promptInt filtered it already.
		return lineIndex[line], true
	}
	var firstArg = func(args []string) string {
		if len(args) == 0 {
			return ""
		}
		return args[0]
	}
	var (
		// Save production chain.
		cmdSave = editor.NewCommand(
			"S",
			"Save the current production chain.",
			[]string{},
			func(_ []string) error {
				if err := c.client.SaveProduct(p); err != nil {
					c.logger.Warnf("unable to save production chain: %s", err.Error())
					fmt.Println("Error saving production chain, try again.")
					return err
				}
				fmt.Println("Production chain saved.")
				return nil
			})

		// Save the production chain and exit the editor.
		cmdSaveQuit = editor.NewCommand(
			"SQ",
			"Save the current production chain and exit the editor.",
			[]string{},
			func(_ []string) error {
				if err := c.client.SaveProduct(p); err != nil {
					c.logger.Warnf("unable to save production chain: %s", err.Error())
					fmt.Println("Error saving production chain, try again.")
					return err
				}
				fmt.Println("Production chain saved.")
				return editor.ErrExitEditor
			})

		// View the production chain details.
		cmdView = editor.NewCommand(
			"V",
			"Print the current production chain's details.",
			[]string{},
			func(_ []string) error {
				fmt.Println()
				c.printProductInfo(p)
				return nil
			})

		// View the blueprint inventory overview.
		cmdBlueprintOverview = editor.NewCommand(
			"O",
			"Print the materials inventory details.",
			[]string{},
			func(_ []string) error {
				c.printBlueprintOverview(p)
				return nil
			})

		// View the details of an item in the chain.
		cmdDetail = editor.NewCommand(
			"D",
			"Show detailed information for a specific chain item.",
			[]string{"[#]"},
			func(args []string) error {
				prod, ok := promptLineNumber("Show detail for which line", firstArg(args))
				if !ok {
					return nil
				}
				if prod.ProductID == p.ProductID {
					fmt.Printf("Already showing detail for %s\n", c.getProductName(prod))
					return nil
				}
				fmt.Printf("Showing detail for %s.\n\n", c.getProductName(prod))
				c.printProductInfo(prod)
				fmt.Println("Enter Q or S to return to the previous product.")
				c.newProductEditor(prod).Loop()
				fmt.Printf("Returned to detail for %s\n", c.getProductName(p))
				return nil
			})

		// Edit the average cost per unit of an item in the chain.
		cmdEditCost = editor.NewCommand(
			"C",
			"Set the cost per unit for a specific chain item.",
			[]string{"[#]"},
			func(args []string) error {
				prod, ok := promptLineNumber("Edit cost for which line", firstArg(args))
				if !ok {
					return nil
				}
				prodName := c.getProductName(prod)
				val, ok := c.env.PromptDecimal(fmt.Sprintf("Enter new cost per unit for %s (current: %s)", prodName, prod.Cost()), nil)
				if !ok {
					return nil
				}
				prod.MarketPrice = val
				fmt.Printf("Updated %s per unit cost to %s.\n", prodName, prod.Cost())
				return nil
			})

		// Edit material efficiency of an item in the chain.
		cmdEditME = editor.NewCommand(
			"F",
			"Set the material efficiency for a specific item in the chain.",
			[]string{"[#]"},
			func(args []string) error {
				prod, ok := promptLineNumber("Edit material efficiency for which line", firstArg(args))
				if !ok {
					return nil
				}
				prodName := c.getProductName(prod)
				val, ok := c.env.PromptDecimal(fmt.Sprintf("Enter new material efficiency for %s (current: %s)", prodName, prod.MaterialEfficiency), nil)
				if !ok {
					return nil
				}
				prod.MaterialEfficiency = val
				fmt.Printf("Updated %s material efficiency to %s.\n", prodName, prod.MaterialEfficiency)
				return nil
			})

		// Edit production mode (buy or sell) of an item in the chain.
		cmdEditMode = editor.NewCommand(
			"M",
			"Set the production mode (buy or build) for a specific chain item.",
			[]string{"[#]"},
			func(args []string) error {
				prod, ok := promptLineNumber("Edit production mode for which line", firstArg(args))
				if !ok {
					return nil
				}
				prodName := c.getProductName(prod)
				val, ok := c.env.PromptString(fmt.Sprintf("Enter new mode for %s (current: %s)", prodName, prod.Kind), nil, validateStringIsOneOf([]string{"buy", "build"}))
				if !ok {
					return nil
				}
				if val == "buy" {
					prod.Kind = model.ProductBuy
				} else {
					prod.Kind = model.ProductBuild
				}
				fmt.Printf("Updated %s production mode to %s.\n", prodName, prod.Kind)
				return nil
			})

		// Edit production batch size of an item in the chain.
		cmdEditBatchSize = editor.NewCommand(
			"B",
			"Set the batch size for a specific chain item.",
			[]string{"[#]"},
			func(args []string) error {
				prod, ok := promptLineNumber("Edit batch size for which line", firstArg(args))
				if !ok {
					return nil
				}
				prodName := c.getProductName(prod)
				val, ok := c.env.PromptInt(
					fmt.Sprintf("Enter new batch size for %s (current: %d)", prodName, prod.BatchSize),
					nil,
					validateIntGreaterThan(0))
				if !ok {
					return nil
				}
				prod.BatchSize = val
				fmt.Printf("Updated %s batch size to %d.\n", prodName, prod.BatchSize)
				return nil
			})

		// Edit the sell price of the final product.
		cmdEditSellPrice = editor.NewCommand(
			"P",
			"Set the sell price per unit for the final product.",
			[]string{},
			func(_ []string) error {
				prodName := c.getProductName(p)
				val, ok := c.env.PromptDecimal(fmt.Sprintf("Enter new sell price for %s (current: %s)", prodName, p.MarketPrice), nil)
				if !ok {
					return nil
				}
				p.MarketPrice = val
				fmt.Printf("Updated %s sell price to %s.\n", prodName, p.MarketPrice)
				return nil
			})

		// Edit the sell region of the final product.
		cmdEditSellRegion = editor.NewCommand(
			"R",
			"Set the market region for the production chain.",
			[]string{},
			func(_ []string) error {
				region, ok := c.env.PromptRegion("Specify Region", "")
				if !ok {
					return nil
				}
				var setRegion func(*model.Product)
				setRegion = func(p *model.Product) {
					p.MarketRegionID = region.RegionID
					for _, m := range p.Materials {
						setRegion(m)
					}
				}
				setRegion(p)
				prod, err := c.client.UpdateProductPrices(p)
				if err != nil {
					c.logger.Errorf("unable to fetch market prices for region %d: %s", region.RegionID, err.Error())
					fmt.Println("Error loading production chain prices, try again.")
				}
				*p = *prod
				fmt.Printf("Updated %s target region to %s.\n", c.getProductName(p), c.getRegionName(p.MarketRegionID))
				return nil
			})

		// Update all product costs per unit to current market values.
		cmdUpdateMarketPrices = editor.NewCommand(
			"U",
			"Update market prices with regional data from evemarketer.com.",
			[]string{},
			func(_ []string) error {
				prod, err := c.client.UpdateProductPrices(p)
				if err != nil {
					c.logger.Errorf("unable to fetch market prices for region %d: %s", p.MarketRegionID, err.Error())
					fmt.Println("Error loading production chain prices, try again.")
					return nil
				}
				*p = *prod
				fmt.Println("Production chain prices updated.")
				return nil
			})
	)

	return editor.New(c.env, `The production chain editor is an interactive application for managing arbitrary production chains. Individual components can be tagged as either "buy" or "build". Cost projections, with material efficiency and batch size accounted for, are updated accordingly. The target market region and target final sell price can also be modified for the production chain as a whole.

When invoking a tool and omitting an optional parameter, an interactive prompt will begin to collect the necessary information.

The current product is always line item 0, which can be used when specifying a line number.`,
		cmdSave,
		cmdSaveQuit,
		cmdView,
		cmdBlueprintOverview,
		cmdDetail,
		cmdEditCost,
		cmdEditBatchSize,
		cmdEditME,
		cmdEditMode,
		cmdEditSellPrice,
		cmdEditSellRegion,
		cmdUpdateMarketPrices,
	)
}
