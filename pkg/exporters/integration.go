package exporters

import (
	"context"
	"fmt"

	rhtapAPI "github.com/konflux-ci/cli/api/v1alpha1"
	"github.com/konflux-ci/cli/cmd/rhtap/commands/config"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

func TransformIntegrationTestScenario(ctx context.Context, fetchedResourceList runtime.Object, cloneConfig config.CloneConfig, localCache []runtime.Object) ([]runtime.Object, error) {

	var selectedResources []runtime.Object

	itsList, ok := fetchedResourceList.(*rhtapAPI.IntegrationTestScenarioList)

	if !ok {
		return nil, fmt.Errorf("resources of type integrationTestScenarioList were not passsed")
	}

	for _, integrationTest := range itsList.Items {
		if cloneConfig.AllApplications || cloneConfig.ApplicatioName == integrationTest.Spec.Application {
			transformedITS := &rhtapAPI.IntegrationTestScenario{
				TypeMeta: integrationTest.TypeMeta,
				ObjectMeta: v1.ObjectMeta{
					Name:      integrationTest.Name,
					Namespace: cloneConfig.TargetNamespace,
				},
				Spec: integrationTest.Spec,
			}
			selectedResources = append(selectedResources, transformedITS)
		}
	}
	return selectedResources, nil
}

func GenerateYAMLForIntegrationTestScenario(ctx context.Context, transformedResources []runtime.Object) ([][]byte, error) {
	var resourcesInYAML [][]byte
	for _, resource := range transformedResources {
		its := resource.(*rhtapAPI.IntegrationTestScenario)
		inBytes, err := yaml.Marshal(its)
		if err != nil {
			return nil, err
		}
		resourcesInYAML = append(resourcesInYAML, inBytes)
	}
	return resourcesInYAML, nil
}
func GetIntegrationTests(ctx context.Context, namespace string, cloneConfig config.CloneConfig, client *kubernetes.Clientset) (runtime.Object, error) {
	var its = &rhtapAPI.IntegrationTestScenarioList{}
	err := client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1beta1/namespaces/%s/integrationtestscenarios", namespace)).
		Do(context.TODO()).Into(its)
	return its, err
}
