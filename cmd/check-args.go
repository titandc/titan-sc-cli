package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "regexp"
)

func cmdNeed1UUID(cmd *cobra.Command, args []string) error {
    _ = cmd
    if len(args) < 1 {
        return fmt.Errorf("Invalid number argument: need 1 UUID;\n")
    }

    if !checkUUIDFormat(args[0]) {
        return fmt.Errorf("Invalid UUID format `%s'\n", args[0])
    }
    return nil
}

func cmdNeed1Args(cmd *cobra.Command, args []string) error {
    _ = cmd
    if len(args) < 1 {
        return fmt.Errorf("Invalid number argument: need 1 argument;\n")
    }
    return nil
}

func checkUUIDFormat(uuid string) bool {
    r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
    return r.MatchString(uuid)
}
