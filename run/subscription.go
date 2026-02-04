package run

import (
	"fmt"
	"time"
	"titan-sc/api"

	"github.com/spf13/cobra"
)

func (run *RunMiddleware) SubscriptionList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	companyOID, err := run.GetDefaultCompanyOID(cmd)
	if err != nil {
		run.OutputError(err)
		return
	}
	allStates, _ := cmd.Flags().GetBool("all")

	subscriptions, err := run.API.GetSubscriptionList(companyOID, !allStates)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(subscriptions)
	} else {
		run.printSubscriptionList(subscriptions)
	}
}

func (run *RunMiddleware) printSubscriptionList(subscriptions []api.Subscription) {
	if len(subscriptions) == 0 {
		fmt.Println("No subscriptions found.")
		return
	}

	table := NewTable("NAME", "DOCUMENT", "STATE", "FREQUENCY", "NEXT BILLING", "AMOUNT HT", "OID")
	table.SetNoColor(!run.Color)

	for _, sub := range subscriptions {
		amountHT := fmt.Sprintf("%.2f€", float64(sub.Amount.HT)/100)
		nextBilling := "-"
		if sub.NextBillingDate > 0 {
			// API returns timestamp - check if it's in milliseconds (too large for seconds)
			ts := sub.NextBillingDate
			if ts > 9999999999 {
				ts = ts / 1000 // Convert from milliseconds to seconds
			}
			nextBilling = time.Unix(ts, 0).Format("2006-01-02")
		}

		var stateColorFn, amountColorFn func(string) string
		if run.Color {
			stateColorFn = StateColorFn(sub.State)
			amountColorFn = ColorFn("green")
		}

		table.AddRow(
			ColName(sub.Name),
			Col(sub.DocumentNumber),
			ColColor(sub.State, stateColorFn),
			Col(sub.Frequency),
			Col(nextBilling),
			ColColor(amountHT, amountColorFn),
			ColOID(sub.OID),
		)
	}

	table.Print()
}

func (run *RunMiddleware) SubscriptionDetail(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	subscriptionOID, _ := cmd.Flags().GetString("subscription-oid")

	subscription, err := run.API.GetSubscription(subscriptionOID)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(subscription)
	} else {
		run.printSubscriptionDetail(subscription)
	}
}

func (run *RunMiddleware) printSubscriptionDetail(sub *api.Subscription) {
	fmt.Printf("%s\n", run.Colorize("Subscription:", "cyan"))
	fmt.Printf("  OID:             %s\n", sub.OID)
	fmt.Printf("  Name:            %s\n", run.Colorize(sub.Name, "cyan"))
	fmt.Printf("  Document Number: %s\n", sub.DocumentNumber)
	fmt.Printf("  State:           %s\n", GetStateColorized(run.Color, sub.State))
	fmt.Printf("  Frequency:       %s\n", sub.Frequency)
	if sub.NextFrequency != "" && sub.NextFrequency != sub.Frequency {
		fmt.Printf("  Next Frequency:  %s\n", sub.NextFrequency)
	}
	if sub.NextBillingDate > 0 {
		ts := sub.NextBillingDate
		if ts > 9999999999 {
			ts = ts / 1000 // Convert from milliseconds to seconds
		}
		fmt.Printf("  Next Billing:    %s\n", time.Unix(ts, 0).Format("2006-01-02"))
	}

	fmt.Printf("\n%s\n", run.Colorize("Company:", "cyan"))
	fmt.Printf("  OID:             %s\n", sub.CompanyOID)
	if sub.Company.Name != "" {
		fmt.Printf("  Name:            %s\n", sub.Company.Name)
	}
	if sub.Company.ClientNumber != "" {
		fmt.Printf("  Client Number:   %s\n", sub.Company.ClientNumber)
	}

	fmt.Printf("\n%s\n", run.Colorize("Billing:", "cyan"))
	fmt.Printf("  Amount HT:       %.2f€\n", float64(sub.Amount.HT)/100)
	fmt.Printf("  Amount TTC:      %.2f€\n", float64(sub.Amount.TTC)/100)
	if sub.PaymentMethodOID != "" {
		fmt.Printf("  Payment Method:  %s\n", sub.PaymentMethodOID)
	}
	if sub.PaymentDisabled {
		fmt.Printf("  Payment:         %s\n", run.Colorize("DISABLED", "yellow"))
	}

	if sub.CreatedAt > 0 {
		fmt.Printf("\n%s\n", run.Colorize("Timestamps:", "cyan"))
		fmt.Printf("  Created:         %s\n", time.Unix(sub.CreatedAt, 0).Format("2006-01-02 15:04:05"))
		if sub.UpdatedAt > 0 {
			fmt.Printf("  Updated:         %s\n", time.Unix(sub.UpdatedAt, 0).Format("2006-01-02 15:04:05"))
		}
	}
}
