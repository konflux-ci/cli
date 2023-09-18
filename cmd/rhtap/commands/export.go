package commands

import (
	"context"
	"fmt"
	"log"

	// TODO: Figure out how to avoid the dependency hell and use APIs
	// from remote repositories.
	rhapAPI "github.com/redhat-appstudio/rhtap-cli/api"
	//intergrationServiceApi "github.com/redhat-appstudio/integration-service/api/v1beta1"
)

// Export is a function to export an Application and its associated resources to
// an importable YAML file.
func Export(args []string, cloneConfig *CloneConfig) error {

	// ensure we've been able to pass everything OK

	fmt.Printf("application name : %s , source namespace : %s , target namespace : %s, overrides : %v,",
		args[0],
		cloneConfig.SourceNamespace,
		cloneConfig.TargetNamespace,
		cloneConfig.ComponentSourceURLOverrides[0],
	)

	// config from kubeconfig
	client, clientError := NewOpenShiftClient()

	if clientError != nil {
		log.Fatal(clientError)
		return clientError
	}
	var applications = &rhapAPI.ApplicationList{}

	// TODO: No need to get all Applications, only get the relevant one
	// and ensure it actually exists
	err := client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1alpha1/namespaces/%s/applications", cloneConfig.SourceNamespace)).
		Do(context.TODO()).Into(applications)

	if err != nil {
		// TODO: Wrap it into something meaningful.
		return err
	}
	fmt.Println(applications)

	var components = &rhapAPI.ComponentList{}
	err = client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1alpha1/namespaces/%s/components", cloneConfig.SourceNamespace)).
		Do(context.TODO()).Into(components)

	if err != nil {
		// TODO: Wrap it into something meaningful.
		return err
	}
	// TODO: Remove non-relevant Components from the list.
	// TODO: Read-up "overrides" and recreate Components which need to have a new source
	// code URL
	// TODO: Re-create Component resources with image references instead of
	// source code URL references.

	var integrationTestScenarios = &rhapAPI.IntegrationTestScenarioList{}
	err = client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1beta1/namespaces/%s/integrationtestscenarios", cloneConfig.SourceNamespace)).
		Do(context.TODO()).
		Into(integrationTestScenarios)

	// TODO: Remove non-relevant IntegrationTestScenarios

	// TODO: Write out a YAML for future use.

	return err
}
