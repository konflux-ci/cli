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

func TransformSnapshots(ctx context.Context, fetchedResourceList runtime.Object, cloneConfig config.CloneConfig, localCache []runtime.Object) ([]runtime.Object, error) {

	var selectedResources []runtime.Object

	itsList, ok := fetchedResourceList.(*konfluxAPI.SnapshotList)

	if !ok {
		return nil, fmt.Errorf("resource of type SnapshotList was not passsed")
	}

	for _, snapshot := range itsList.Items {

		// TODO: Discard all snapshots which are not relevant.
		if cloneConfig.AllApplications || cloneConfig.ApplicatioName == snapshot.Spec.Application {
			transformedSnapshot := &konfluxAPI.Snapshot{
				TypeMeta: snapshot.TypeMeta,
				ObjectMeta: v1.ObjectMeta{
					Name:      snapshot.Name,
					Namespace: cloneConfig.TargetNamespace,
				},
				Spec: snapshot.Spec,
			}
			selectedResources = append(selectedResources, transformedSnapshot)
		}
	}
	return selectedResources, nil
}

func GenerateYAMLForSnapshots(ctx context.Context, transformedResources []runtime.Object) ([][]byte, error) {
	var resourcesInYAML [][]byte
	for _, resource := range transformedResources {
		snapshot := resource.(*konfluxAPI.Snapshot)
		inBytes, err := yaml.Marshal(snapshot)
		if err != nil {
			return nil, err
		}
		resourcesInYAML = append(resourcesInYAML, inBytes)
	}
	return resourcesInYAML, nil
}
func GetSnapshots(ctx context.Context, namespace string, cloneConfig config.CloneConfig, client *kubernetes.Clientset) (runtime.Object, error) {
	var its = &konfluxAPI.SnapshotList{}
	err := client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1alpha1/namespaces/%s/snapshots", namespace)).
		Do(context.TODO()).Into(its)
	return its, err
}
