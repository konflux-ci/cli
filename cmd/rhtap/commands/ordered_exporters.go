package commands

import "github.com/redhat-appstudio/rhtap-cli/pkg/exporters"

func getOrderedExporters() []exporters.ResourceExport {
	return []exporters.ResourceExport{
		{
			Transform:    exporters.TransformApplication,
			GenerateYAML: exporters.GenerateYAMLForApplication,
			Get:          exporters.GetApplications,
		},
		{
			Transform:    exporters.TransformComponent,
			GenerateYAML: exporters.GenerateYAMLForComponent,
			Get:          exporters.GetComponents,
		},
		{
			Transform:    exporters.TransformIntegrationTestScenario,
			GenerateYAML: exporters.GenerateYAMLForIntegrationTestScenario,
			Get:          exporters.GetIntegrationTests,
		},
		{
			Transform:    exporters.TransformSecret,
			GenerateYAML: exporters.GenerateYAMLForSecret,
			Get:          exporters.GetSecrets,
		},
		{
			Transform:    exporters.TransformEnvironment,
			GenerateYAML: exporters.GenerateYAMLForEnvironments,
			Get:          exporters.GetEnvironments,
		},
		{
			Transform:    exporters.TransformSnapshots,
			GenerateYAML: exporters.GenerateYAMLForSnapshots,
			Get:          exporters.GetSnapshots,
		},
		{
			Transform:    exporters.TransformSnapshotEnvironmentBindings,
			GenerateYAML: exporters.GenerateYAMLForSnapshotEnvironmentBindings,
			Get:          exporters.GetSnapshotEnvironmentBindings,
		},
	}
}
