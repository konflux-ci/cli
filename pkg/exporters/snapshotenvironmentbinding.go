package exporters

import (
	"context"
	"fmt"

	rhtapAPI "github.com/redhat-appstudio/rhtap-cli/api/v1alpha1"
	"github.com/redhat-appstudio/rhtap-cli/cmd/rhtap/commands/config"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

func TransformSnapshotEnvironmentBindings(ctx context.Context, fetchedResourceList runtime.Object, cloneConfig config.CloneConfig, localCache []runtime.Object) ([]runtime.Object, error) {

	var selectedResources []runtime.Object

	itsList, ok := fetchedResourceList.(*rhtapAPI.SnapshotEnvironmentBindingList)

	if !ok {
		return nil, fmt.Errorf("resource of type SnapshotEnvironmentBindingList was not passsed")
	}

	for _, seb := range itsList.Items {

		// TODO: Discard all snapshots which are not relevant.
		if cloneConfig.AllApplications || cloneConfig.ApplicatioName == seb.Spec.Application {
			transformedSnapshot := &rhtapAPI.SnapshotEnvironmentBinding{
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

func GenerateYAMLForSnapshotEnvironmentBindings(ctx context.Context, transformedResources []runtime.Object) ([]ResourceYAML, error) {
	var resourcesInYAML []ResourceYAML
	for _, resource := range transformedResources {
		snapshot := resource.(*rhtapAPI.SnapshotEnvironmentBinding)
		inBytes, err := yaml.Marshal(snapshot)
		if err != nil {
			return nil, err
		}
		resourcesInYAML = append(resourcesInYAML, ResourceYAML{inBytes})
	}
	return resourcesInYAML, nil
}
func GetSnapshotEnvironmentBindings(ctx context.Context, namespace string, cloneConfig config.CloneConfig, client *kubernetes.Clientset) (runtime.Object, error) {
	var its = &rhtapAPI.SnapshotEnvironmentBindingList{}
	err := client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1alpha1/namespaces/%s/snapshotenvironmentbindings", namespace)).
		Do(context.TODO()).Into(its)
	return its, err
}
