package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	// TODO: Figure out how to avoid the dependency hell and use APIs
	// from remote repositories.
	rhapAPI "github.com/redhat-appstudio/rhtap-cli/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
	//intergrationServiceApi "github.com/redhat-appstudio/integration-service/api/v1beta1"
)

// Export is a function to export an Application and its associated resources to
// an importable YAML file.
func Export(args []string, cloneConfig *CloneConfig) error {

	// ensure we've been able to pass everything OK

	// TODO: Handle error if empty.
	cloneConfig.ApplicatioName = args[0]

	// config from kubeconfig
	client, clientError := NewOpenShiftClient()

	if clientError != nil {
		log.Fatal(clientError)
		return clientError
	}
	var singleApplication = rhapAPI.Application{}

	err := client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1alpha1/namespaces/%s/applications/%s", cloneConfig.SourceNamespace, args[0])).
		Do(context.TODO()).Into(&singleApplication)

	if err != nil {
		// TODO: Wrap it into something meaningful.
		return err
	}

	// TODO: use explicit file modes?
	file, err := os.Create(cloneConfig.OutputFile)
	if err != nil {
		fmt.Println("Could not open ", cloneConfig.OutputFile)
		return err
	}

	defer file.Close()

	if err != nil {
		fmt.Println("Could not open ", cloneConfig.OutputFile)
		return err
	}

	exportableApplication := generateExportableApplication(singleApplication, cloneConfig.TargetNamespace)
	inBytes, err := yaml.Marshal(exportableApplication)
	if err != nil {
		fmt.Println("Error exporting resources ", err.Error())
	}
	_, err = file.Write(inBytes)
	if err != nil {
		fmt.Println("Error exporting resources ", err.Error())
	}

	var components = &rhapAPI.ComponentList{}
	err = client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1alpha1/namespaces/%s/components", cloneConfig.SourceNamespace)).
		Do(context.TODO()).Into(components)

	if err != nil {
		// TODO: Wrap it into something meaningful.
		return err
	}

	for _, c := range components.Items {

		if strings.Compare(c.Spec.Application, cloneConfig.ApplicatioName) == 0 {
			exportableComponent := generateExportableComponent(c, cloneConfig.TargetNamespace, cloneConfig.ComponentSourceURLOverrides)
			// TODO check for nil and report error

			inBytes, err := yaml.Marshal(exportableComponent)
			if err != nil {
				fmt.Println("Error exporting resources ", err.Error())
			}
			file.WriteString("---\n")
			_, err = file.Write(inBytes)
			if err != nil {
				fmt.Println("Error exporting resources ", err.Error())
			}

		}
	}

	// TODO: Read-up "overrides" and recreate Components which need to have a new source
	// code URL
	// TODO: Re-create Component resources with image references instead of
	// source code URL references.

	integrationTestScenarios := &rhapAPI.IntegrationTestScenarioList{}
	err = client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1beta1/namespaces/%s/integrationtestscenarios", cloneConfig.SourceNamespace)).
		Do(context.TODO()).
		Into(integrationTestScenarios)

	// TODO: Remove non-relevant IntegrationTestScenarios
	for _, itc := range integrationTestScenarios.Items {
		if itc.Spec.Application == cloneConfig.ApplicatioName {
			exportableITS := generateExportableIntegrationTestScenario(itc, cloneConfig.TargetNamespace)
			inBytes, err := yaml.Marshal(exportableITS)
			if err != nil {
				fmt.Println("Error exporting resources ", err.Error())
			}
			file.WriteString("---\n")
			_, err = file.Write(inBytes)
			if err != nil {
				fmt.Println("Error exporting resources ", err.Error())
			}

		}
	}
	return err
}

func generateOverridesMap(overrides string) map[string]string {
	overridesMap := map[string]string{}
	for _, s := range strings.Fields(overrides) {
		overridesMap[strings.Split(s, "=")[0]] = strings.Split(s, "=")[1]
	}
	return overridesMap
}

func generateExportableIntegrationTestScenario(
	fetchedIntegrationTestScenario rhapAPI.IntegrationTestScenario,
	targetNamespace string) *rhapAPI.IntegrationTestScenario {

	return &rhapAPI.IntegrationTestScenario{
		TypeMeta: fetchedIntegrationTestScenario.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:      fetchedIntegrationTestScenario.Name,
			Namespace: targetNamespace,
		},
		Spec: fetchedIntegrationTestScenario.Spec,
	}

}
func generateExportableApplication(fetchedApplication rhapAPI.Application, targetNamespace string) *rhapAPI.Application {
	return &rhapAPI.Application{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "appstudio.redhat.com/v1alpha1",
			Kind:       "Application",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fetchedApplication.Name,
			Namespace: targetNamespace,
			Annotations: map[string]string{
				"application.thumbnail": "1",
			},
		},
		Spec: rhapAPI.ApplicationSpec{
			DisplayName: fetchedApplication.Name,
		},
	}
}

func generateExportableComponent(fetchedComponent rhapAPI.Component, targetNamespace string, overrides string) *rhapAPI.Component {
	exportableComponent := rhapAPI.Component{}
	overridesMap := generateOverridesMap(overrides)
	val, ok := overridesMap[fetchedComponent.Name]
	if ok {
		fmt.Printf("Found an override for %s : %s \n", fetchedComponent.Name, val)

		exportableComponent = rhapAPI.Component{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "appstudio.redhat.com/v1alpha1",
				Kind:       "Component",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fetchedComponent.Name,
				Namespace: targetNamespace,
				Annotations: map[string]string{
					"skip-initial-checks":       "true",
					"image.redhat.com/generate": `{"visibility": "public"}`,
				},
			},
			Spec: rhapAPI.ComponentSpec{
				Application:   fetchedComponent.Spec.Application,
				ComponentName: fetchedComponent.Spec.ComponentName,
				Source: rhapAPI.ComponentSource{
					ComponentSourceUnion: rhapAPI.ComponentSourceUnion{
						GitSource: &rhapAPI.GitSource{
							URL:           val,
							Context:       fetchedComponent.Spec.Source.GitSource.Revision,
							Revision:      fetchedComponent.Spec.Source.GitSource.Revision,
							DockerfileURL: fetchedComponent.Spec.Source.GitSource.DockerfileURL,
						},
					},
				},
				Replicas:   fetchedComponent.Spec.Replicas,
				Resources:  fetchedComponent.Spec.Resources,
				Env:        fetchedComponent.Spec.Env,
				TargetPort: fetchedComponent.Spec.TargetPort,
			},
		}
	} else {
		// import as image
		exportableComponent = rhapAPI.Component{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "appstudio.redhat.com/v1alpha1",
				Kind:       "Component",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fetchedComponent.Name,
				Namespace: targetNamespace,
				Annotations: map[string]string{
					"skip-initial-checks": "true",
				},
			},
			Spec: rhapAPI.ComponentSpec{
				Application:    fetchedComponent.Spec.Application,
				ComponentName:  fetchedComponent.Spec.ComponentName,
				Replicas:       fetchedComponent.Spec.Replicas,
				Resources:      fetchedComponent.Spec.Resources,
				Env:            fetchedComponent.Spec.Env,
				TargetPort:     fetchedComponent.Spec.TargetPort,
				ContainerImage: fetchedComponent.Spec.ContainerImage,
				Source: rhapAPI.ComponentSource{
					ComponentSourceUnion: rhapAPI.ComponentSourceUnion{
						GitSource: &rhapAPI.GitSource{
							URL:           fetchedComponent.Spec.Source.GitSource.URL,
							Context:       fetchedComponent.Spec.Source.GitSource.Revision,
							Revision:      fetchedComponent.Spec.Source.GitSource.Revision,
							DockerfileURL: fetchedComponent.Spec.Source.GitSource.DockerfileURL,
						},
					},
				},
			},
		}
	}
	return &exportableComponent
}
