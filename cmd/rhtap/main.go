package main

import (
	"fmt"

	"github.com/redhat-appstudio/rhtap-cli/cmd/rhtap/commands"
	"github.com/spf13/cobra"
)

var cloneConfig = new(commands.CloneConfig)

var rcommand = &cobra.Command{
	Use: "rhtap",
}

var versionCommand = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version 2")
	},
}

func init() {
	rcommand.AddCommand(versionCommand)
	rcommand.AddCommand(exportCommand)

	applicationCommand.PersistentFlags().StringVarP(&cloneConfig.SourceNamespace, "from-namespace", "f", "", "Namespace from which the Application is being cloned.")
	applicationCommand.PersistentFlags().StringVarP(&cloneConfig.TargetNamespace, "to-namespace", "t", "", "Namespace to which the Application is being cloned.")
	applicationCommand.PersistentFlags().StringVarP(&cloneConfig.ComponentSourceURLOverrides, "overrides", "o", "", "Overwrite the source code url for specific components")
	applicationCommand.PersistentFlags().StringVarP(&cloneConfig.ComponentSourceURLOverrides, "skip", "s", "", "List of components to be skipped")
	applicationCommand.PersistentFlags().StringVarP(&cloneConfig.ComponentSourceURLOverrides, "write-to", "w", "", "Local filesystem path where the YAML would be written out to.")

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
		commands.List(args, cloneConfig)
	},
}
