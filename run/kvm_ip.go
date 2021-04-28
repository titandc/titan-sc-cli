package run

import (
	"fmt"
	"github.com/spf13/cobra"
	"titan-sc/api"
)

func (run *RunMiddleware) KVMIPGetInfos(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	kvmip, err := run.API.KVMIPGet(serverUUID)
	if err != nil {
		run.OutputError(err)
		return
	}

	if !run.HumanReadable {
		printAsJson(kvmip)
	} else {
		run.IPKvmPrint(serverUUID, *kvmip)
	}
}

func (run *RunMiddleware) KVMIPStart(cmd *cobra.Command, args []string) {
	_ = args
	run.IPKvmAction("start", cmd)
}

func (run *RunMiddleware) KVMIPStop(cmd *cobra.Command, args []string) {
	_ = args
	run.IPKvmAction("stop", cmd)
}

func (run *RunMiddleware) IPKvmAction(action string, cmd *cobra.Command) {
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	act := api.APIServerAction{
		Action: action,
	}
	_, apiReturn, err := run.API.SendRequestToAPI(api.HTTPPut, "/compute/servers/"+serverUUID+"/ipkvm", act)
	if err != nil || apiReturn != nil {
		run.handleErrorAndGenericOutput(apiReturn, err)
		return
	}
}

func (run *RunMiddleware) IPKvmPrint(serverUUID string, ipkvm api.APIKvmIP) {

	fmt.Printf("IP KVM informations:\n  Server UUID: %s\n  Status: %s\n", serverUUID, ipkvm.Status)
	if ipkvm.URI != "" {
		fmt.Println("  URI:", ipkvm.URI)
	}
}
