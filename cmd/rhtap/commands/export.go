package commands

import (
	"context"
	"fmt"
	"log"

	hasApplicationAPI "github.com/redhat-appstudio/rhtap-cli/api"

	//intergrationServiceApi "github.com/redhat-appstudio/integration-service/api/v1beta1"

	v1 "k8s.io/api/core/v1"
)

// List is a function to list number of pods in the cluster
func List(args []string, cloneConfig *CloneConfig) error {

	// ensure we've been able to pass everything OK

	fmt.Printf("application name : %s , source namespace : %s , target namespace : %s, overrides : %v,",
		args[0],
		cloneConfig.SourceNamespace,
		cloneConfig.TargetNamespace,
		cloneConfig.ComponentSourceURLOverrides[0],
	)

	client, clientError := NewOpenShiftClient()

	if clientError != nil {
		log.Fatal(clientError)
		return clientError
	}

	var mycms = &v1.ConfigMapList{}
	err := client.RESTClient().Get().AbsPath("/api/v1/namespaces/shbose-tenant/configmaps").
		Do(context.TODO()).
		Into(mycms)

	fmt.Println(len(mycms.Items))
	fmt.Println(err)

	var applications = &hasApplicationAPI.ApplicationList{}

	err = client.RESTClient().Get().AbsPath("/apis/appstudio.redhat.com/v1alpha1/namespaces/shbose-tenant/applications").
		Do(context.TODO()).
		Into(applications)

	fmt.Println(applications)
	fmt.Println(err)

	var components = &hasApplicationAPI.ComponentList{}
	err = client.RESTClient().Get().AbsPath("/apis/appstudio.redhat.com/v1alpha1/namespaces/shbose-tenant/components").
		Do(context.TODO()).
		Into(components)

	fmt.Println(components)
	fmt.Println(err)

	var integrationTestScenarios = hasApplicationAPI.IntegrationTestScenarioList{}
	err = client.RESTClient().Get().AbsPath("/apis/appstudio.redhat.com/v1beta1/namespaces/shbose-tenant/integrationtestscenarios").
		Do(context.TODO()).
		Into(components)

	fmt.Println(integrationTestScenarios)
	fmt.Println(err)

	return err

	/*
		pods, podlisterror := client.CoreV1().Pods(envVar.Namespace).List(context.TODO(), metav1.ListOptions{})
		if podlisterror == nil {
			fmt.Printf("\nThe number of pods are %d \n", len(pods.Items))
		} else {
			log.Fatal(podlisterror)
		}

		return podlisterror
	*/
}
