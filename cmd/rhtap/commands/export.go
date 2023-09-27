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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
	//intergrationServiceApi "github.com/redhat-appstudio/integration-service/api/v1beta1"
)

// Export is a function to export an Application and its associated resources to
// an importable YAML file.
func Export(args []string, cloneConfig *CloneConfig) error {

	// ensure we've been able to pass everything OK

	if len(args) == 0 {
		fmt.Println("Application not specified")
		return fmt.Errorf("application not specified")
	}
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

	err = writeKubernetesResourceToFile(file, inBytes)
	if err != nil {
		return err
	}

	var components = &rhapAPI.ComponentList{}
	err = client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1alpha1/namespaces/%s/components", cloneConfig.SourceNamespace)).
		Do(context.TODO()).Into(components)

	if err != nil {
		return fmt.Errorf("could not fetch components in namespace %s: %v", cloneConfig.SourceNamespace, err.Error())
	}

	for _, c := range components.Items {

		if strings.Compare(c.Spec.Application, cloneConfig.ApplicatioName) == 0 {
			exportableComponent := generateExportableComponent(c, cloneConfig.TargetNamespace, cloneConfig.ComponentSourceURLOverrides)
			// TODO check for nil and report error

			inBytes, err := yaml.Marshal(exportableComponent)
			if err != nil {
				fmt.Println("Error exporting resources ", err.Error())
			}
			err = writeKubernetesResourceToFile(file, inBytes)
			if err != nil {
				return err
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

	if err != nil {
		return fmt.Errorf("could not fetch IntegrationTestScenarios in namespace %s: %v", cloneConfig.SourceNamespace, err.Error())
	}

	// TODO: Remove non-relevant IntegrationTestScenarios
	for _, itc := range integrationTestScenarios.Items {
		if itc.Spec.Application == cloneConfig.ApplicatioName {
			exportableITS := generateExportableIntegrationTestScenario(itc, cloneConfig.TargetNamespace)
			inBytes, err := yaml.Marshal(exportableITS)
			if err != nil {
				fmt.Println("Error exporting resources ", err.Error())
			}
			err = writeKubernetesResourceToFile(file, inBytes)
			if err != nil {
				return err
			}

		}
	}

	// Get the relevant custom secrets exported.
	// As of now, we only have the snyk token to be exported.
	// There's no good way to know if this Application uses the snyk token.
	// We'll export it anyway. Definitely room for improvement.

	snykTokenSecret := &v1.Secret{}
	err = client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1beta1/namespaces/%s/secrets/snyk-secret", cloneConfig.SourceNamespace)).
		Do(context.TODO()).
		Into(snykTokenSecret)

	if err != nil {
		return fmt.Errorf("could not fetch snyk-secret from the namespace %s : %v", cloneConfig.SourceNamespace, err.Error())
	}

	exportableSecret := generateExportableSecret(snykTokenSecret, cloneConfig.TargetNamespace)
	inBytes, err = yaml.Marshal(exportableSecret)
	if err != nil {
		return fmt.Errorf("could not convert Secret to bytes prior to writing out to file : %s", err.Error())
	}

	err = writeKubernetesResourceToFile(file, inBytes)
	if err != nil {
		return fmt.Errorf("could not write Secret to file : %v", err.Error())
	}

	return err
}

func writeKubernetesResourceToFile(file *os.File, inBytes []byte) error {
	_, err := file.WriteString("---\n")
	if err != nil {
		return fmt.Errorf("yaml delimiter could not be added : %v", err)
	}
	_, err = file.Write(inBytes)
	if err != nil {
		return fmt.Errorf("error exporting kubernetes resource to the file : %v", err)
	}

	return nil
}

func generateOverridesMap(overrides string) map[string]string {
	overridesMap := map[string]string{}
	for _, s := range strings.Fields(overrides) {
		overridesMap[strings.Split(s, "=")[0]] = strings.Split(s, "=")[1]
	}
	return overridesMap
}

func generateExportableSecret(
	fetchedSecret *v1.Secret, targetNamespace string) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fetchedSecret.Name,
			Namespace: targetNamespace,
		},
		TypeMeta: fetchedSecret.TypeMeta,
		Data:     fetchedSecret.Data,
		Type:     fetchedSecret.Type,
	}
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
