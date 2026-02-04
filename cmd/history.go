package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) HistoryCmdAdd() {
	historyEvent := &cobra.Command{
		Use:     "history --server-oid SERVER_OID | [--company-oid COMPANY_OID] [--number n --offset n]",
		Aliases: []string{"hist"},
		Short:   "List latest events on a server or a company.",
		Long: `List the n latest events of a server or a whole company.

25 events displayed by default, must not exceed 50.

You can query by:
  - Server OID (-s): Events for a specific server
  - Company OID (-c): Events for an entire company

If --company-oid is not specified, your default company will be used.`,
		Run:     cmd.runMiddleware.EventHistory,
		GroupID: "resources",
	}

	cmd.RootCommand.AddCommand(historyEvent)
	historyEvent.Flags().IntP("number", "n", 25, "Amount of event(s) to retrieve, must not exceed 50.")
	historyEvent.Flags().IntP("offset", "o", 0, "Offset to begin event list.")
	historyEvent.Flags().StringP("company-oid", "c", "", "Company OID (uses your default company if not specified).")
	historyEvent.Flags().StringP("server-oid", "s", "", "Set server OID.")
	historyEvent.MarkFlagsMutuallyExclusive("server-oid", "company-oid")
}
