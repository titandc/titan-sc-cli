package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	DefaultURI = "https://sc.titandc.net/api/v1"
	HTTPGet    = http.MethodGet
	HTTPPut    = http.MethodPut
	HTTPPost   = http.MethodPost
	HTTPDelete = http.MethodDelete
	// Get N history entry
	NHistory = 25
)

/*
 *
 *
 ******************
 * Global variable
 ******************
 *
 *
 */
var API = &APITitan{
	Token:         "",
	URI:           "",
	HumanReadable: false,
	Color:         true,
}

/*
 *
 *
 ******************
 * history function
 ******************
 *
 *
 */
func (API *APITitan) HistoryCompanyEvent(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)

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
		API.HistoryByCompany(strNumber, companyUUID)
	} else {
		API.HistoryByServer(strNumber, serverUUID)
	}
}

func (API *APITitan) HistoryByCompany(number, companyUUID string) {
	err := API.SendAndResponse(HTTPGet, "/compute/servers/events?nb="+
		number+"&company_uuid="+companyUUID, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !API.HumanReadable {
		API.PrintJson()
	} else {
		history := make([]APIHistoryEvent, 0)
		if err := json.Unmarshal(API.RespBody, &history); err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, event := range history {
			date := API.DateFormat(event.Timestamp)
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

func (API *APITitan) HistoryByServer(number, serverUUID string) {
	err := API.SendAndResponse(HTTPGet, "/compute/servers/"+serverUUID+
		"/events?nb="+number, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !API.HumanReadable {
		API.PrintJson()
	} else {
		history := make([]APIHistoryEvent, 0)
		if err := json.Unmarshal(API.RespBody, &history); err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, event := range history {
			date := API.DateFormat(event.Timestamp)
			fmt.Printf("Event type:\r\t\t%s (%s)\n"+
				"Event status:\r\t\t%s\n"+
				"Date:\r\t\t%s\n\n",
				event.Type, event.State, event.Status, date)
		}
	}
}

/*
 *
 *
 *****************
 * IP Kvm function
 *****************
 *
 *
 */
func (API *APITitan) KVMIPGetInfos(cmd *cobra.Command, args []string) {
	serverUUID := args[0]
	API.ParseGlobalFlags(cmd)

	err := API.SendAndResponse(HTTPGet, "/compute/servers/"+serverUUID+"/ipkvm", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !API.HumanReadable {
		API.PrintJson()
	} else {
		API.IPKvmPrint(serverUUID)
	}
}

func (API *APITitan) KVMIPStart(cmd *cobra.Command, args []string) {
	API.ParseGlobalFlags(cmd)
	API.IPKvmAction("start", args[0])
}

func (API *APITitan) KVMIPStop(cmd *cobra.Command, args []string) {
	API.ParseGlobalFlags(cmd)
	API.IPKvmAction("stop", args[0])
}

func (API *APITitan) IPKvmAction(action, serverUUID string) {
	act := APIServerAction{
		Action: action,
	}
	API.SendAndPrintDefaultReply(HTTPPut, "/compute/servers/"+serverUUID+"/ipkvm", act)
}

func (API *APITitan) IPKvmPrint(serverUUID string) {
	ipkvm := &APIKvmIP{}
	if err := json.Unmarshal(API.RespBody, ipkvm); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("IP KVM informations:\n  Server UUID: %s\n  Status: %s\n", serverUUID, ipkvm.Status)
	if ipkvm.URI != "" {
		fmt.Println("  URI:", ipkvm.URI)
	}
}

/*
 *
 *
 *********************
 * Other API function
 *********************
 *
 *
 */
func (API *APITitan) WeatherMap(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)

	if err := API.SendAndResponse(HTTPGet, "/weather", nil); err != nil {
		fmt.Println(err.Error())
		return
	}

	if !API.HumanReadable {
		API.PrintJson()
	} else {
		weatherMap := &APIWeatherMap{}
		if err := json.Unmarshal(API.RespBody, weatherMap); err != nil {
			log.Println(err.Error())
			return
		}
		fmt.Printf("Titan Weather Map:\n"+
			"  Compute: %s\n"+
			"  Storage: %s\n"+
			"  Public network: %s\n"+
			"  Private network: %s\n",
			weatherMap.Compute, weatherMap.Storage,
			weatherMap.PublicNetwork, weatherMap.PrivateNetwork)
	}
}

func (API *APITitan) ManagedServices(cmd *cobra.Command, args []string) {
	API.ParseGlobalFlags(cmd)
	companyUUID := args[0]
	managedServicesOpts := APIManagedServices{
		Company: companyUUID,
	}
	API.SendAndPrintDefaultReply(HTTPPost, "/compute/managed_services", managedServicesOpts)
}

func (API *APITitan) IsAdmin() (bool, error) {
	if err := API.SendAndResponse(HTTPGet, "/auth/user/isadmin", nil); err != nil {
		return false, err
	}

	buffer := IsAdminStruct{}
	if err := json.Unmarshal(API.RespBody, &buffer); err != nil {
		return false, err
	}
	return buffer.Admin, nil
}

/*
 *
 *
 ******************
 * Utils function
 ******************
 *
 *
 */
func (API *APITitan) VersionAPI(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)

	if err := API.SendAndResponse(HTTPGet, "/version", nil); err != nil {
		fmt.Println(err.Error())
		return
	}

	if !API.HumanReadable {
		API.PrintJson()
	} else {
		version := &APIVersion{}
		if err := json.Unmarshal(API.RespBody, version); err != nil {
			log.Println(err.Error())
			return
		}
		fmt.Println("Titan API version:", version.Version)
	}
}

func (API *APITitan) DefaultPrintReturn() {
	APIRet := &APIReturn{}
	err := json.Unmarshal(API.RespBody, APIRet)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if APIRet.Success != "" {
		fmt.Println(APIRet.Success)
	} else {
		fmt.Println("Error:", APIRet.Error)
	}
}

func (API *APITitan) PrintJson() {
	dst := &bytes.Buffer{}
	if err := json.Indent(dst, API.RespBody, "", "  "); err != nil {
		fmt.Printf(string(API.RespBody))
		return
	}
	fmt.Printf(dst.String())
}

func (API *APITitan) ParseGlobalFlags(cmd *cobra.Command) {
	var err error

	API.HumanReadable, err = cmd.Flags().GetBool("human")
	if err != nil {
		API.HumanReadable = false
	}

	API.Color, err = cmd.Flags().GetBool("color")
	if err != nil {
		API.Color = false
	}
}

func (API *APITitan) SendAndResponse(method, path string, req []byte) error {
	request, err := http.NewRequest(method, API.URI+path, bytes.NewBuffer(req))
	if err != nil {
		return err
	}
	request.Header.Add("X-API-KEY", API.Token)
	request.Header.Add("Titan-Cli-Os", API.CLIos)
	request.Header.Add("Titan-Cli-Version", API.CLIVersion)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	API.RespBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Check command error
	APIRet := &APIReturn{}
	err = json.Unmarshal(API.RespBody, APIRet)
	if err == nil && APIRet.Error != "" {
		if API.HumanReadable {
			return fmt.Errorf("Error: %s", APIRet.Error)
		} else {
			return fmt.Errorf(string(API.RespBody))
		}
	}
	return nil
}

func (API *APITitan) SendAndPrintDefaultReply(httpMethod, path string, httpData interface{}) {
	var reqData []byte
	var err error

	reqData = nil
	if httpData != nil {
		reqData, err = json.Marshal(httpData)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	err = API.SendAndResponse(httpMethod, path, reqData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !API.HumanReadable {
		API.PrintJson()
	} else {
		API.DefaultPrintReturn()
	}
}
