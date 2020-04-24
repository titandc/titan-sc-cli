
package cmd

import (
    "github.com/spf13/cobra"
    . "titan-sc/api"
)

var historyEvent = &cobra.Command{
    Use: "history [--company-uuid | --server-uuid] (--number)",
    Aliases: []string{"hist"},
    Short: "List latest events.",
    Long: "List latest events.",
    Run: API.HistoryCompanyEvent,
}


func historyCmdAdd() {
    rootCmd.AddCommand(historyEvent)
    historyEvent.Flags().IntP("number", "n", 25, "number event.")
    historyEvent.Flags().StringP("company-uuid", "c", "", "Set company UUID.")
    historyEvent.Flags().StringP("server-uuid", "u", "", "Set server UUID.")
}
