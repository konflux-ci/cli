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

func TransformApplication(ctx context.Context, fetchedResourceList runtime.Object, cloneConfig config.CloneConfig, localCache []runtime.Object) ([]runtime.Object, error) {

	var selectedResources []runtime.Object

	applicationsList, ok := fetchedResourceList.(*rhtapAPI.ApplicationList)

	if !ok {
		return nil, fmt.Errorf("resource of type ApplicationList was not passsed")
	}

	for _, application := range applicationsList.Items {
		if cloneConfig.AllApplications || application.Name == cloneConfig.ApplicatioName {
			transformedApplication := &rhtapAPI.Application{
				TypeMeta: application.TypeMeta,
				ObjectMeta: v1.ObjectMeta{
					Name:      application.Name,
					Namespace: cloneConfig.TargetNamespace,
					Annotations: map[string]string{
						"application.thumbnail": "1",
					},
				},
				Spec: application.Spec,
			}
			selectedResources = append(selectedResources, transformedApplication)
		}
	}
	return selectedResources, nil
}

func GenerateYAMLForApplication(ctx context.Context, transformedResources []runtime.Object) ([]ResourceYAML, error) {
	var resourcesInYAML []ResourceYAML
	for _, resource := range transformedResources {
		application := resource.(*rhtapAPI.Application)
		inBytes, err := yaml.Marshal(application)
		if err != nil {
			return nil, err
		}
		resourcesInYAML = append(resourcesInYAML, ResourceYAML{inBytes})
	}
	return resourcesInYAML, nil
}

func GetApplications(ctx context.Context, namespace string, cloneConfig config.CloneConfig, client *kubernetes.Clientset) (runtime.Object, error) {
	var applications = &rhtapAPI.ApplicationList{}
	err := client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1alpha1/namespaces/%s/applications", namespace)).
		Do(context.TODO()).Into(applications)
	return applications, err
}
