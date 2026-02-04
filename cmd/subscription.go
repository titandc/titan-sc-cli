package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) SubscriptionCmdAdd() {
	subscription := &cobra.Command{
		Use:     "subscription",
		Aliases: []string{"sub"},
		Short:   "Manage billing subscriptions.",
	}

	subscriptionList := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all subscriptions.",
		Long: `List all active and pending subscriptions for your company.

Subscriptions group recurring invoiced items together. When creating a new server,
you can optionally specify a subscription OID to add the server to an existing
subscription instead of creating a new one.`,
		Run: cmd.runMiddleware.SubscriptionList,
	}

	subscriptionShow := &cobra.Command{
		Use:     "show --subscription-oid OID",
		Aliases: []string{"get"},
		Short:   "Show subscription detail.",
		Long:    "Show detailed information about a specific subscription.",
		Run:     cmd.runMiddleware.SubscriptionDetail,
	}

	cmd.RootCommand.AddCommand(subscription)
	subscription.AddCommand(subscriptionList, subscriptionShow)

	subscriptionList.Flags().StringP("company-oid", "c", "", "Filter by company OID (optional).")
	subscriptionList.Flags().BoolP("all", "a", false, "Show all subscriptions including canceled (default: only ongoing).")
	subscriptionShow.Flags().StringP("subscription-oid", "s", "", "Subscription OID to show.")
	_ = subscriptionShow.MarkFlagRequired("subscription-oid")
}
