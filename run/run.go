package run

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
	"titan-sc/api"
)

type RunMiddleware struct {
	HumanReadable bool
	Color         bool
	CLIVersion    string
	CLIos         string
	API           *api.API
}

func NewRunMiddleware(api *api.API) *RunMiddleware {
	return &RunMiddleware{
		API:           api,
		HumanReadable: false,
		Color:         true,
	}
}

func (run *RunMiddleware) ParseGlobalFlags(cmd *cobra.Command) {
	var err error

	run.HumanReadable, err = cmd.Flags().GetBool("human")
	if err != nil {
		run.HumanReadable = false
	}

	run.Color, err = cmd.Flags().GetBool("color")
	if err != nil {
		run.Color = false
	}
}

func (run *RunMiddleware) handleErrorAndGenericOutput(apiReturn *api.APIReturn, err error) {
	// Communication or marshalling error
	if err != nil {
		run.OutputError(err)
		return
	}

	// API parsed error (automatically handle JSON vs string)
	if apiReturn != nil {
		run.printAPIReturn(apiReturn)
		return
	}
}

func (run *RunMiddleware) OutputError(err error) {
	if !run.HumanReadable {
		printAsJson(err.Error())
	} else {
		fmt.Println(err.Error())
	}
}

func (run *RunMiddleware) printAPIReturn(apiReturn *api.APIReturn) {
	if run.HumanReadable {
		printAPIReturnAsString(apiReturn)
	} else {
		printAsJson(apiReturn)
	}
	return
}

func printAsJson(data interface{}) {
	switch data.(type) {
	case []byte:
		// do nothing
		fmt.Println(string(data.([]byte)))
	default:
		dataToPrint, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Printf("{'error': '%s'}", err.Error())
			return
		}
		fmt.Printf(string(dataToPrint))
	}

}

func printAPIReturnAsString(apiReturn *api.APIReturn) {
	if apiReturn.Error != "" {
		fmt.Printf("Error: %s", apiReturn.Error)
	}
	if apiReturn.Success != "" {
		fmt.Printf("Success: %s", apiReturn.Success)
	}
	fmt.Printf(" (code: %s)\n", apiReturn.Code)
}

func millisecondsToTime(timestamp int64) time.Time {
	return time.Unix(0, timestamp*int64(time.Millisecond))
}

func keyboardPromptToLower(promptString string) string {
	// Read user input
	fmt.Print(promptString)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	lowerText := strings.ToLower(text)
	return lowerText
}

func serverStateSetColor(color bool, state string) string {
	if !color {
		return state
	}

	colorState := state
	switch state {
	case "deleted":
		colorState = "\033[1;31m" + state + "\033[0m"
	case "started":
		colorState = "\033[1;32m" + state + "\033[0m"
	case "stopped":
		colorState = "\033[1;33m" + state + "\033[0m"
	}
	return colorState
}
