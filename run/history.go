package run

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

const (
	// Get N history entry
	NHistory = 25
)

/*
 *******************
 * History functions
 *******************
 */
func (run *RunMiddleware) HistoryCompanyEvent(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	companyUUID, _ := cmd.Flags().GetString("company-uuid")
	number, _ := cmd.Flags().GetInt("number")

	if serverUUID != "" && companyUUID != "" ||
		(companyUUID == "" && serverUUID == "") {
		_ = cmd.Help()
		return
	}

	if number < 1 {
		number = NHistory
	}
	strNumber := fmt.Sprintf("%d", number)

	if companyUUID != "" {
		run.historyByCompany(strNumber, companyUUID)
	} else {
		run.historyByServer(strNumber, serverUUID)
	}
}

func (run *RunMiddleware) historyByCompany(number, companyUUID string) {
	history, apiReturn, err := run.API.HistoryByCompany(number, companyUUID)
	if err != nil {
		run.OutputError(err)
		return
	}
	if apiReturn != nil && apiReturn.Error != "" {
		run.OutputError(errors.New(apiReturn.Error))
		return
	}

	if !run.HumanReadable {
		printAsJson(history)
	} else {
		for _, event := range history {
			date := run.DateFormat(event.Timestamp)
			fmt.Printf("Server UUID:\r\t\t%s\n"+
				"Server Name:\r\t\t%s\n"+
				"Event type:\r\t\t%s (%s)\n"+
				"Event status:\r\t\t%s\n"+
				"Date:\r\t\t%s\n\n",
				event.Server.UUID, event.Server.Name, event.Type,
				event.State, event.Status, date)
		}
	}
}

func (run *RunMiddleware) historyByServer(number, serverUUID string) {
	history, apiReturn, err := run.API.HistoryByServer(number, serverUUID)
	if err != nil {
		run.OutputError(err)
		return
	}
	if apiReturn != nil && apiReturn.Error != "" {
		run.OutputError(errors.New(apiReturn.Error))
		return
	}

	if !run.HumanReadable {
		printAsJson(history)
	} else {
		for _, event := range history {
			date := run.DateFormat(event.Timestamp)
			fmt.Printf("Event type:\r\t\t%s (%s)\n"+
				"Event status:\r\t\t%s\n"+
				"Date:\r\t\t%s\n\n",
				event.Type, event.State, event.Status, date)
		}
	}
}
