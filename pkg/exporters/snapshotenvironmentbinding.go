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

func TransformSnapshotEnvironmentBindings(ctx context.Context, fetchedResourceList runtime.Object, cloneConfig config.CloneConfig, localCache []runtime.Object) ([]runtime.Object, error) {

	var selectedResources []runtime.Object

	itsList, ok := fetchedResourceList.(*konfluxAPI.SnapshotEnvironmentBindingList)

	if !ok {
		return nil, fmt.Errorf("resource of type SnapshotEnvironmentBindingList was not passsed")
	}

	for _, seb := range itsList.Items {

		// TODO: Discard all snapshots which are not relevant.
		if cloneConfig.AllApplications || cloneConfig.ApplicatioName == seb.Spec.Application {
			transformedSnapshot := &konfluxAPI.SnapshotEnvironmentBinding{
				TypeMeta: seb.TypeMeta,
				ObjectMeta: v1.ObjectMeta{
					Name:      seb.Name,
					Namespace: cloneConfig.TargetNamespace,
				},
				Spec: seb.Spec,
			}
			selectedResources = append(selectedResources, transformedSnapshot)
		}
	}
	return selectedResources, nil
}

func GenerateYAMLForSnapshotEnvironmentBindings(ctx context.Context, transformedResources []runtime.Object) ([][]byte, error) {
	var resourcesInYAML [][]byte
	for _, resource := range transformedResources {
		snapshot := resource.(*konfluxAPI.SnapshotEnvironmentBinding)
		inBytes, err := yaml.Marshal(snapshot)
		if err != nil {
			return nil, err
		}
		resourcesInYAML = append(resourcesInYAML, inBytes)
	}
	return resourcesInYAML, nil
}
func GetSnapshotEnvironmentBindings(ctx context.Context, namespace string, cloneConfig config.CloneConfig, client *kubernetes.Clientset) (runtime.Object, error) {
	var its = &konfluxAPI.SnapshotEnvironmentBindingList{}
	err := client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1alpha1/namespaces/%s/snapshotenvironmentbindings", namespace)).
		Do(context.TODO()).Into(its)
	return its, err
}
