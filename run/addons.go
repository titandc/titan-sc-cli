package run

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (run *RunMiddleware) ServerAddon(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")

	addons, err := run.API.ServerAddon(serverOID)
	if err != nil {
		run.OutputError(err)
		return
	}
	if run.JSONOutput {
		printAsJson(addons)
	} else {
		if len(addons.UpgradableItems) == 0 {
			fmt.Println("No addons available.")
			return
		}

		table := NewTable("NAME", "PRICE (HT)", "UNIT", "OID")
		table.SetNoColor(!run.Color)
		for _, addon := range addons.UpgradableItems {
			pricing := fmt.Sprintf("%0.1f%s", float64(addon.PriceUnitHT/1000), addon.Currency)
			unit := fmt.Sprintf("%d %s", addon.ItemUnit.Value, addon.ItemUnit.Unit)

			var priceColorFn func(string) string
			if run.Color {
				priceColorFn = ColorFn("green")
			}

			table.AddRow(
				ColName(addon.Name),
				ColColor(pricing, priceColorFn),
				Col(unit),
				ColOID(addon.OID),
			)
		}
		table.Print()
	}
}
