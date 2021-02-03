package cmd

import (
	"github.com/spf13/cobra"
	. "titan-sc/api"
)

var server = &cobra.Command{
	Use:     "server",
	Aliases: []string{"srv"},
	Short:   "Manage servers.",
	Long:    "Manage servers.",
}

var serverList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Show detail of all servers in your companies.",
	Long:    "Show detail of all servers in your companies.",
	Run:     API.ServerList,
}

var serverDetail = &cobra.Command{
	Use:     "show --server-uuid SERVER_UUID",
	Aliases: []string{"get"},
	Short:   "Show server detail.",
	Long:    "Show detailed information about a server.",
	Run:     API.ServerDetail,
}

var serverStart = &cobra.Command{
	Use:   "start --server-uuid SERVER_UUID",
	Short: "Send an action request to start a server.",
	Long: "Send an action request to start a server." +
		"\nList of available actions:" +
		"\n  start\n  stop\n  hardstop\n  reboot\n",
	Run:  API.ServerStart,
}

var serverStop = &cobra.Command{
	Use:   "stop --server-uuid SERVER_UUID",
	Short: "Send an action request to stop a server.",
	Long: "Send an action request to stop a server." +
		"\nList of available actions:" +
		"\n  start\n  stop\n  hardstop\n  reboot\n",
	Run:  API.ServerStop,
}

var serverRestart = &cobra.Command{
	Use:     "restart --server-uuid SERVER_UUID",
	Aliases: []string{"reboot"},
	Short:   "Send an action request to restart a server.",
	Long: "Send an action request to restart a server." +
		"\nList of available actions:" +
		"\n  start\n  stop\n  hardstop\n  reboot\n",
	Run:  API.ServerRestart,
}

var serverHardstop = &cobra.Command{
	Use:   "hardstop --server-uuid SERVER_UUID",
	Short: "Send an action request to hardstop a server.",
	Long: "Send an action request to hardstop a server." +
		"\nList of available actions:" +
		"\n  start\n  stop\n  hardstop\n  reboot\n",
	Run:  API.ServerHardstop,
}

var serverChangeName = &cobra.Command{
	Use:   "rename --server-uuid SERVER_UUID --name NEW_NAME",
	Short: "Send a request to change server's name.",
	Long:  "Send a request to change server's name.",
	Run:   API.ServerChangeName,
}

var serverChangeReverse = &cobra.Command{
	Use:   "reverse --server-uuid SERVER_UUID --reverse NEW_REVERSE",
	Short: "Send a request to change server's reverse.",
	Long:  "Send a request to change server's reverse.",
	Run:   API.ServerChangeReverse,
}

var serverLoadISO = &cobra.Command{
	Use:     "load-iso --uri HTTPS_URI --server-uuid SERVER_UUID",
	Aliases: []string{"li"},
	Short:   "Send a request to load an ISO from HTTPS.",
	Long:    "Send a request to load a bootable ISO from HTTPS.",
	Run:     API.ServerLoadISO,
}

var serverUnloadISO = &cobra.Command{
	Use:     "unload-iso --server-uuid SERVER_UUID",
	Aliases: []string{"ui"},
	Short:   "Send a request to unload previously loaded ISO(s).",
	Long:    "Send a request to unload all previously loaded ISO(s).",
	Run:     API.ServerUnloadISO,
}

var ServerAddonsList = &cobra.Command{
	Use:   "addons",
	Short: "List all server addons.",
	Long:  "List all server addons.",
	Run:   API.AddonsListAll,
}

var serverGetTemplateList = &cobra.Command{
	Use:   "templates",
	Short: "List all server template.",
	Long:  "List all server template.",
	Run:   API.ServerGetTemplateList,
}

var serverCreate = &cobra.Command{
	Use:   "create --os OS_NAME --os-version OS_VERSION --plan SC1/SC2/SC3",
	Short: "Send a request for create new server's.",
	Long:  "Send a request for create new server's.\nGet os and os version see: server templates.",
	Run:   API.ServerCreate,
}

var serverDelete = &cobra.Command{
	Use:     "delete --server-uuid SERVER_UUID",
	Aliases: []string{"del"},
	Short:   "Send a request for delete a server's.",
	Long:    "Send a request for delete a server's.",
	Run:     API.ServerDelete,
}

var serverReset = &cobra.Command{
	Use:   "reset --server-uuid SERVER_UUID --os OS_NAME --os-version OS_VERSION",
	Short: "Send a request for reset a server's.",
	Long:  "Send a request for reset a server's.",
	Run:   API.ServerReset,
}

func serverCmdAdd() {
	rootCmd.AddCommand(server)
	server.AddCommand(serverList, serverDetail, serverStart,
		serverStop, serverRestart, serverHardstop, serverLoadISO,
		serverUnloadISO, serverChangeName, serverChangeReverse, ServerAddonsList,
		serverGetTemplateList, serverCreate, serverDelete, serverReset)
	serverList.Flags().StringP("company-uuid", "c", "", "Set company UUID.")

	serverDetail.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	_ = serverDetail.MarkFlagRequired("server-uuid")

	serverStart.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	_ = serverStart.MarkFlagRequired("server-uuid")

	serverStop.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	_ = serverStop.MarkFlagRequired("server-uuid")

	serverRestart.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	_ = serverRestart.MarkFlagRequired("server-uuid")

	serverHardstop.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	_ = serverHardstop.MarkFlagRequired("server-uuid")

	serverLoadISO.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	serverLoadISO.Flags().StringP("uri", "u", "", "Set remote ISO URI (HTTPS only).")
	_ = serverLoadISO.MarkFlagRequired("server-uuid")
	_ = serverLoadISO.MarkFlagRequired("uri")

	serverUnloadISO.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	_ = serverUnloadISO.MarkFlagRequired("server-uuid")

	serverChangeName.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	serverChangeName.Flags().StringP("name", "n", "", "Set new server's name.")
	_ = serverChangeName.MarkFlagRequired("server-uuid")
	_ = serverChangeName.MarkFlagRequired("name")

	serverChangeReverse.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	serverChangeReverse.Flags().StringP("reverse", "r", "", "Set new server's reverse.")
	_ = serverChangeReverse.MarkFlagRequired("server-uuid")
	_ = serverChangeReverse.MarkFlagRequired("reverse")

	// Server create
	serverCreate.Flags().StringP("plan", "p", "", "Choose your server plan.")
	serverCreate.Flags().StringP("os", "", "", "Set you OS name.")
	serverCreate.Flags().StringP("os-version", "", "", "Set your os version.")
	serverCreate.Flags().StringP("login", "", "", "Set login.")
	serverCreate.Flags().StringP("password", "", "", "Set password for login.")
	serverCreate.Flags().StringP("network-uuid", "", "", "Set network UUID for managed network.")
	serverCreate.Flags().Int64P("quantity", "", 1, "Set quantity.")
	serverCreate.Flags().IntP("cpu-addon", "c", 0, "Number CPU addons.")
	serverCreate.Flags().IntP("ram-addon", "r", 0, "Number  RAM size (GB) addons.")
	serverCreate.Flags().IntP("disk-addon", "d", 0, "Number Disk size (GB) addons.")
	serverCreate.Flags().StringP("ssh-keys-name", "", "", "Set ssh keys: keyname1,keyname2,...,keynameX.")
	_ = serverCreate.MarkFlagRequired("plan")
	_ = serverCreate.MarkFlagRequired("os")
	_ = serverCreate.MarkFlagRequired("os-version")

	// server reset
	serverReset.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	serverReset.Flags().StringP("os", "", "", "Set you OS name.")
	serverReset.Flags().StringP("os-version", "", "", "Set your os version.")
	serverReset.Flags().StringP("password", "", "", "Set password for login.")
	serverReset.Flags().StringP("ssh-keys-name", "", "", "Set ssh keys: keyname1,keyname2,...,keynameX.")
	_ = serverReset.MarkFlagRequired("os")
	_ = serverReset.MarkFlagRequired("os-version")
	_ = serverReset.MarkFlagRequired("server-uuid")

	serverDelete.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	_ = serverDelete.MarkFlagRequired("server-uuid")
}
