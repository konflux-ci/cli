package exporters

import (
	"context"
	"fmt"
	"strings"

	rhtapAPI "github.com/redhat-appstudio/rhtap-cli/api/v1alpha1"
	"github.com/redhat-appstudio/rhtap-cli/cmd/rhtap/commands/config"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

/*
// Transform transforms a Component for exporting into a backup/retarget or cloning for use in a different namespace
func (c *ComponentExport) Transform(ctx context.Context, obj runtime.Object, cloneConfig config.CloneConfig, asIs bool) (runtime.Object, error) {

	fetchedComponent, OK := obj.(*rhtapAPI.Component)
	if !OK {
		return nil, fmt.Errorf("did not find a Component resource")
	}
	transformedComponent := &rhtapAPI.Component{
		TypeMeta: fetchedComponent.TypeMeta,

		ObjectMeta: v1.ObjectMeta{
			Name:      fetchedComponent.Name,
			Namespace: cloneConfig.TargetNamespace,
		},
		Spec: fetchedComponent.Spec,
	}
	return transformedComponent, nil
}
*/

func TransformComponent(ctx context.Context, fetchedResourceList runtime.Object, cloneConfig config.CloneConfig, localCache []runtime.Object) ([]runtime.Object, error) {
	var selectedResources []runtime.Object

	componentList, ok := fetchedResourceList.(*rhtapAPI.ComponentList)

	if !ok {
		return nil, fmt.Errorf("resources of type componentList were not passsed")
	}

	for _, component := range componentList.Items {
		var transformedComponent *rhtapAPI.Component

		if cloneConfig.AllApplications || cloneConfig.ApplicatioName == component.Spec.Application {

			if component.Spec.Source.GitSource != nil && shouldSkip(cloneConfig.ComponentSourceURLskip, component.Spec.Source.GitSource.URL) {
				continue
			}

			transformedComponent = &rhtapAPI.Component{
				TypeMeta: component.TypeMeta,

				ObjectMeta: v1.ObjectMeta{
					Name:      component.Name,
					Namespace: cloneConfig.TargetNamespace,
					Annotations: map[string]string{
						"skip-initial-checks": "true",
					},
				},
				Spec: component.Spec,
			}
			if !cloneConfig.AsPrebuiltImages {
				// for embargo flows, this annotation would be skipped.
				// for backup and restore, this annotation would reset the robot account token.
				// for re-target namespaces, this annotation would create a new image repo.

				// TODO: match the visibility of the original Component/repo.

				transformedComponent.ObjectMeta.Annotations["image.redhat.com/generate"] = `{"visibility": "public"}`
			}
			selectedResources = append(selectedResources, transformedComponent)
		}

	}

	return selectedResources, nil
}

func GenerateYAMLForComponent(ctx context.Context, transformedResources []runtime.Object) ([][]byte, error) {
	var resourcesInYAML [][]byte
	for _, resource := range transformedResources {
		component := resource.(*rhtapAPI.Component)
		inBytes, err := yaml.Marshal(component)
		if err != nil {
			return nil, err
		}
		resourcesInYAML = append(resourcesInYAML, inBytes)
	}
	return resourcesInYAML, nil
}

/*
type ComponentFetch struct{}

var _ ResourceFetcher = &ComponentFetch{}

func (c *ComponentFetch) Get(ctx context.Context, namespace string, cloneConfig config.CloneConfig, client *kubernetes.Clientset) (runtime.Object, error) {
	var components = &rhtapAPI.ComponentList{}
	err := client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1alpha1/namespaces/%s/components", namespace)).
		Do(context.TODO()).Into(components)
	return components, err
}
*/

func GetComponents(ctx context.Context, namespace string, cloneConfig config.CloneConfig, client *kubernetes.Clientset) (runtime.Object, error) {
	var components = &rhtapAPI.ComponentList{}
	err := client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1alpha1/namespaces/%s/components", namespace)).
		Do(context.TODO()).Into(components)
	return components, err
}

func shouldSkip(listOfURLsToBeSkipped string, sourceCodeURL string) bool {
	URLList := strings.Split(listOfURLsToBeSkipped, ",")
	for _, url := range URLList {
		if url == sourceCodeURL {
			return true
		}
	}
	return false
}
