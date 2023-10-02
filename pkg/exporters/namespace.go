package exporters

import (
	"context"
	"fmt"

	"github.com/redhat-appstudio/rhtap-cli/cmd/rhtap/commands/config"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"

	"sigs.k8s.io/yaml"
)

func TransformNamespace(ctx context.Context, fetchedResourceList runtime.Object, cloneConfig config.CloneConfig, localCache []runtime.Object) ([]runtime.Object, error) {

	var selectedResources []runtime.Object

	namespaceList, ok := fetchedResourceList.(*corev1.NamespaceList)

	if !ok {
		return nil, fmt.Errorf("resource of type ApplicationList was not passsed")
	}

	for _, ns := range namespaceList.Items {
		val, ok := ns.Labels["toolchain.dev.openshift.com/tier"]
		if ok && val == "appstudio" {
			fmt.Println(ns.Name)

			selectedResources = append(selectedResources, &corev1.Namespace{
				TypeMeta: ns.TypeMeta,
				ObjectMeta: v1.ObjectMeta{
					Name: ns.Name,
					// TODO: Which annotations and labels should we have?

				},
				Spec: ns.Spec,
			})
		}
	}

	return selectedResources, nil
}

func GenerateYAMLForNamespace(ctx context.Context, transformedResources []runtime.Object) ([]ResourceYAML, error) {
	var resourcesInYAML []ResourceYAML
	for _, resource := range transformedResources {
		ns := resource.(*corev1.Namespace)
		inBytes, err := yaml.Marshal(ns)
		if err != nil {
			return nil, err
		}
		resourcesInYAML = append(resourcesInYAML, ResourceYAML{inBytes})
	}
	return resourcesInYAML, nil
}

func GetNamespaces(ctx context.Context, namespace string, cloneConfig config.CloneConfig, client *kubernetes.Clientset) (runtime.Object, error) {
	var namespaces = &corev1.NamespaceList{}
	err := client.RESTClient().Get().AbsPath("/api/v1/namespaces").
		Do(context.TODO()).Into(namespaces)
	return namespaces, err
}
