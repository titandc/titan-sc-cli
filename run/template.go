package run

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (run *RunMiddleware) TemplateList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	templates, apiReturn, err := run.API.ListTemplates()
	// Render error output
	if err != nil || apiReturn != nil {
		run.handleErrorAndGenericOutput(apiReturn, err)
		return
	}

	if run.JSONOutput {
		printAsJson(templates)
	} else {
		// Separate system templates from user images
		type templateRow struct {
			OS      string
			Version string
			OID     string
		}
		type imageRow struct {
			Name     string
			OID      string
			Base     string
			DiskSize string
		}

		var systemRows []templateRow
		var imageRows []imageRow

		for _, template := range templates {
			if len(template.Versions) == 0 {
				continue
			}
			if template.IsImage {
				for _, version := range template.Versions {
					name := version.Name
					if name == "" {
						name = fmt.Sprintf("%s %s", version.OS, version.Version)
					}
					diskSize := "-"
					if version.ImageInfo != nil && version.ImageInfo.DiskSize > 0 {
						diskSize = fmt.Sprintf("%d GB", version.ImageInfo.DiskSize)
					}
					imageRows = append(imageRows, imageRow{
						Name:     name,
						OID:      version.OID,
						Base:     fmt.Sprintf("%s %s", version.OS, version.Version),
						DiskSize: diskSize,
					})
				}
			} else {
				for _, v := range template.Versions {
					systemRows = append(systemRows, templateRow{
						OS:      template.OS,
						Version: v.Version,
						OID:     v.OID,
					})
				}
			}
		}

		// Print system templates table
		if len(systemRows) > 0 {
			fmt.Println(run.Colorize("System Templates:", "cyan"))
			table := NewTable("OS", "VERSION", "OID")
			table.SetNoColor(!run.Color)
			for _, row := range systemRows {
				table.AddRow(
					Col(row.OS),
					Col(row.Version),
					ColOID(row.OID),
				)
			}
			table.Print()
		}

		// Print user images table
		if len(imageRows) > 0 {
			if len(systemRows) > 0 {
				fmt.Println()
			}
			fmt.Println(run.Colorize("User Images:", "cyan"))
			table := NewTable("NAME", "OID", "BASE", "SIZE")
			table.SetNoColor(!run.Color)
			for _, row := range imageRows {
				table.AddRow(
					Col(row.Name),
					ColOID(row.OID),
					Col(row.Base),
					Col(row.DiskSize),
				)
			}
			table.Print()
		}
	}
}

func (run *RunMiddleware) TemplateShow(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	templateOID, _ := cmd.Flags().GetString("template-oid")

	template, err := run.API.GetTemplateByOID(templateOID)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(template)
	} else {
		fmt.Printf("%s %s\n", run.Colorize("Template:", "cyan"), template.OID)
		if template.Name != "" {
			fmt.Printf("  Name: %s\n", run.Colorize(template.Name, "green"))
		}
		fmt.Printf("  OS: %s\n", template.OS)
		fmt.Printf("  Version: %s\n", template.Version)
		fmt.Printf("  Type: %s\n", template.Type)
		fmt.Printf("  UUID: %s\n", run.Colorize(template.UUID, "dim"))
		fmt.Printf("  Enabled: %t\n", template.Enabled)

		if template.HasLicense != nil {
			fmt.Printf("  Has License: %t\n", *template.HasLicense)
		}

		if template.ImageInfo != nil {
			fmt.Printf("  %s\n", run.Colorize("Image Info:", "cyan"))
			fmt.Printf("    Disk Size: %d GB\n", template.ImageInfo.DiskSize)
			fmt.Printf("    Base Template: %s\n", template.ImageInfo.TemplateOID)
		}
	}
}
