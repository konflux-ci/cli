package exporters

//"github.com/redhat-appstudio/rhtap-cli/cmd/rhtap/commands"

/*
func TestUsageOfExportersAndFetchers(t *testing.T) {
	apiExporters := []ResourceExport{
		{
			Transform:    TransformApplication,
			GenerateYAML: GenerateYAMLForApplication,
			Get:          GetApplications,
		},
		{
			Transform:    TransformComponent,
			GenerateYAML: GenerateYAMLForComponent,
			Get:          GetComponents,
		},
		{
			Transform:    TransformIntegrationTestScenario,
			GenerateYAML: GenerateYAMLForIntegrationTestScenario,
			Get:          GetIntegrationTests,
		},
		{
			Transform:    TransformSecret,
			GenerateYAML: GenerateYAMLForSecret,
			Get:          GetSecrets,
		},
		{
			Transform:    TransformEnvironment,
			GenerateYAML: GenerateYAMLForEnvironments,
			Get:          GetEnvironments,
		},
		{
			Transform:    TransformSnapshots,
			GenerateYAML: GenerateYAMLForSnapshots,
			Get:          GetSnapshots,
		},
		{
			Transform:    TransformSnapshotEnvironmentBindings,
			GenerateYAML: GenerateYAMLForSnapshotEnvironmentBindings,
			Get:          GetSnapshotEnvironmentBindings,
		},
	}

	var localCache []runtime.Object
	// config from kubeconfig
	client, clientError := commands.NewOpenShiftClient()
	assert.NoError(t, clientError)
	assert.NotNil(t, client)

	cloneConfig := config.CloneConfig{
		AllApplications: false,
		AllNamespaces:   false,
		ApplicatioName:  "build",
		SourceNamespace: "rhtap-build-tenant",
		TargetNamespace: "shbose-tenant",
		OutputFile:      "output.yaml",
	}

	var allResourcesAsYaml []ResourceYAML

	for _, apiExporter := range apiExporters {

		// TODO: Apply for all namespaces

		//resourceList, err := apiExporter.Get(context.Background(), "rhtap-build-tenant", cloneConfig, client)
		resourceList, err := apiExporter.Get(context.Background(), "rhtap-build-tenant", cloneConfig, client)

		assert.NoError(t, err)

		localCache = append(localCache, resourceList)

		individualResources, err := apiExporter.Transform(context.Background(), resourceList, cloneConfig, localCache)
		assert.NoError(t, err)

		resourcesAsYaml, err := apiExporter.GenerateYAML(context.Background(), individualResources)
		assert.NoError(t, err)

		allResourcesAsYaml = append(allResourcesAsYaml, resourcesAsYaml...)
	}

	file, err := os.Create(cloneConfig.OutputFile)
	assert.NoError(t, err)
	defer file.Close()

	for _, r := range allResourcesAsYaml {
		_, err := file.WriteString("---\n")
		assert.NoError(t, err)
		_, err = file.Write(r.ResourceBytes)
		assert.NoError(t, err)
	}
}
*/
