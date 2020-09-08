package cmd

import (
	"github.com/spf13/cobra"
	. "titan-sc/api"
)

var historyEvent = &cobra.Command{
	Use:     "history --server-uuid SERVER_UUID | --company-uuid COMPANY_UUID [--amount n]",
	Aliases: []string{"hist"},
	Short:   "List latest events on a server or a company.",
	Long:    "List the n latest events of a server or a whole company (25 events displayed by default, must not exceed 50).",
	Run:     API.HistoryCompanyEvent,
}

func historyCmdAdd() {
	rootCmd.AddCommand(historyEvent)
	historyEvent.Flags().IntP("amount", "n", 25, "Amount of event(s) to retrieve, must not exceed 50.")
	historyEvent.Flags().StringP("company-uuid", "c", "", "Set company UUID.")
	historyEvent.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
}
