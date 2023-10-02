package main

import (
	"fmt"

	"github.com/redhat-appstudio/rhtap-cli/cmd/rhtap/commands"
	"github.com/redhat-appstudio/rhtap-cli/cmd/rhtap/commands/config"
	"github.com/spf13/cobra"
)

var cloneConfig = new(config.CloneConfig)

var rcommand = &cobra.Command{
	Use: "rhtap",
}

var versionCommand = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Generate from commit / build time vars / lgflags
		fmt.Println("Version 2")
	},
}

func init() {
	rcommand.AddCommand(versionCommand)
	rcommand.AddCommand(exportCommand)

	applicationCommand.PersistentFlags().BoolVarP(&cloneConfig.AllNamespaces, "all-projects", "p", false, "When set, all namespaces/projects will be cloned.")
	applicationCommand.PersistentFlags().BoolVarP(&cloneConfig.AllApplications, "all-applications", "a", false, "When set, all Applications in the current namespace will be cloned.")
	applicationCommand.PersistentFlags().StringVarP(&cloneConfig.SourceNamespace, "from", "f", "", "Namespace from which the Application is being cloned.")
	applicationCommand.PersistentFlags().StringVarP(&cloneConfig.TargetNamespace, "to", "t", "", "Namespace to which the Application is being cloned.")

	// "overrides" is to be used for embargo workflows only but
	// they aren't very useful since the user needs to authenticate with the
	// private repo interactively

	//applicationCommand.PersistentFlags().StringVarP(&cloneConfig.ComponentSourceURLOverrides, "overrides", "o", "", "Overwrite the source code url for specific components")

	applicationCommand.PersistentFlags().StringVarP(&cloneConfig.ComponentSourceURLskip, "skip", "s", "", "List of components to be skipped")
	applicationCommand.PersistentFlags().StringVarP(&cloneConfig.OutputFile, "write-to", "w", "", "Local filesystem path where the YAML would be written out to.")

	exportCommand.AddCommand(applicationCommand)
}

func main() {
	rcommand.Execute()
}

var exportCommand = &cobra.Command{
	Use:   "export",
	Short: "Command to export <resources>",
}

var applicationCommand = &cobra.Command{
	Use: "application",
	Run: func(cm *cobra.Command, args []string) {
		fmt.Println("Export Application and associated resources into a YAML file")
		commands.Export(args, cloneConfig)
	},
}
