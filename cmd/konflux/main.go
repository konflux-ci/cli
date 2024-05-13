package main

import (
	"fmt"

	"github.com/konflux-ci/cli/cmd/konflux/commands"
	"github.com/konflux-ci/cli/cmd/konflux/commands/config"
	"github.com/spf13/cobra"
)

var cloneConfig = new(config.CloneConfig)

var rcommand = &cobra.Command{
	Use: "konflux",
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

	applicationCommand.PersistentFlags().StringVarP(&cloneConfig.ComponentSourceURLskip, "skip", "s", "", "List of components to be skipped")
	applicationCommand.PersistentFlags().BoolVarP(&cloneConfig.AsPrebuiltImages, "as-prebuilt-images", "i", false, "Export components such that they could be imported as pre-built images")

	applicationCommand.PersistentFlags().StringVarP(&cloneConfig.OutputDir, "write-to", "w", "", "Local filesystem directory path where the YAML would be written out to.")
	applicationCommand.PersistentFlags().StringVarP(&cloneConfig.Key, "key", "k", "", "Local filesystem path to an existing encryption key")

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
		err := commands.Export(args, cloneConfig)
		if err != nil {
			fmt.Println(err)
		}
	},
}
