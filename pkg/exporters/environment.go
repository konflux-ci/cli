package exporters

import (
	"context"
	"fmt"

	konfluxAPI "github.com/konflux-ci/cli/api/v1alpha1"
	"github.com/konflux-ci/cli/cmd/konflux/commands/config"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

func TransformEnvironment(ctx context.Context, fetchedResourceList runtime.Object, cloneConfig config.CloneConfig, localCache []runtime.Object) ([]runtime.Object, error) {

	var selectedResources []runtime.Object

	envList, ok := fetchedResourceList.(*konfluxAPI.EnvironmentList)

	if !ok {
		return nil, fmt.Errorf("resources of type integrationTestScenarioList were not passsed")
	}

	for _, environment := range envList.Items {
		// skip copying the "development environment"
		if environment.Name != "development" {
			selectedResources = append(selectedResources, &konfluxAPI.Environment{
				TypeMeta: environment.TypeMeta,
				ObjectMeta: v1.ObjectMeta{
					Name:      environment.Name,
					Namespace: cloneConfig.TargetNamespace,
				},
				Spec: environment.Spec,
			})
		}
	}

	return selectedResources, nil
}

func GenerateYAMLForEnvironments(ctx context.Context, transformedResources []runtime.Object) ([][]byte, error) {
	var resourcesInYAML [][]byte
	for _, resource := range transformedResources {
		its := resource.(*konfluxAPI.Environment)
		inBytes, err := yaml.Marshal(its)
		if err != nil {
			return nil, err
		}
		resourcesInYAML = append(resourcesInYAML, inBytes)
	}
	return resourcesInYAML, nil
}
func GetEnvironments(ctx context.Context, namespace string, cloneConfig config.CloneConfig, client *kubernetes.Clientset) (runtime.Object, error) {
	var its = &konfluxAPI.EnvironmentList{}
	err := client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1alpha1/namespaces/%s/environments", namespace)).
		Do(context.TODO()).Into(its)
	return its, err
}
