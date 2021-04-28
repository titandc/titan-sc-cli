package run

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
	"titan-sc/api"
)

func (run *RunMiddleware) AddonsListAll(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	addons, err := run.API.GetAllAddons()
	if err != nil {
		run.OutputError(err)
		return
	}
	if !run.HumanReadable {
		printAsJson(addons)
	} else {
		var w *tabwriter.Writer
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "UUID\tNAME\tVALUE\tUNIT\tSC1  \tSC2  \tSC3  \t\n")
		for _, addon := range addons {
			pircingSC1 := fmt.Sprintf("%0.1f%s", addon.PricingSC1.Value, addon.PricingSC1.Currency)
			pircingSC2 := fmt.Sprintf("%0.1f%s", addon.PricingSC2.Value, addon.PricingSC2.Currency)
			pircingSC3 := fmt.Sprintf("%0.1f%s", addon.PricingSC3.Value, addon.PricingSC3.Currency)
			_, _ = fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%s\t%s\t\n", addon.UUID,
				addon.Name, addon.Amount.Value, addon.Amount.Unit, pircingSC1, pircingSC2, pircingSC3)
		}
		_ = w.Flush()
	}
}

func (run *RunMiddleware) GetAllAddons() ([]api.APIAddonsItem, error) {
	var addons []api.APIAddonsItem

	addons, err := run.API.GetAllAddons()
	if err != nil {
		run.OutputError(err)
		return nil, err
	}
	return addons, nil
}

func (run *RunMiddleware) AddonGetUUIDByName(addons []api.APIAddonsItem, name string) (string, error) {
	for _, addon := range addons {
		if strings.ToLower(addon.Name) == name {
			return addon.UUID, nil
		}
	}
	return "", fmt.Errorf("UUID for addon name <%s> not found", name)
}
