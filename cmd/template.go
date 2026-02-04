package cmd

import "github.com/spf13/cobra"

func (cmd *CMD) TemplateCmdAdd() {
	templateList := &cobra.Command{
		Use:   "list",
		Short: "List all available templates.",
		Long:  "List all available templates (operating systems and user images).",
		Run:   cmd.runMiddleware.TemplateList,
	}

	templateShow := &cobra.Command{
		Use:   "show --template-oid TEMPLATE_OID",
		Short: "Show details of a specific template.",
		Long:  "Show details of a specific template by OID.",
		Run:   cmd.runMiddleware.TemplateShow,
	}
	templateShow.Flags().StringP("template-oid", "o", "", "Template OID.")
	_ = templateShow.MarkFlagRequired("template-oid")

	templateCMD := &cobra.Command{
		Use:     "template",
		Short:   "Manage templates.",
		Long:    "Manage templates (operating systems and user images).",
		GroupID: "resources",
	}

	cmd.RootCommand.AddCommand(templateCMD)
	templateCMD.AddCommand(templateList, templateShow)
}
