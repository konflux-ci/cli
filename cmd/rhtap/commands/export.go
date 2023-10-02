package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redhat-appstudio/rhtap-cli/pkg/exporters"

	"github.com/redhat-appstudio/rhtap-cli/cmd/rhtap/commands/config"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

// Export is a function to export an Application and its associated resources to
// an importable YAML file.
func Export(args []string, cloneConfig *config.CloneConfig) error {
	ctx := context.Background()
	validateConfig(ctx, cloneConfig, args)

	//var localCache []runtime.Object
	// config from kubeconfig
	client, clientError := NewOpenShiftClient()
	if clientError != nil {
		fmt.Println(clientError)
		return clientError
	}

	if cloneConfig.AllNamespaces {
		namepspaces, err := exporters.GetNamespaces(ctx, "", *cloneConfig, client)
		if err != nil {
			return err
		}
		filteredListOfNamespaces, err := exporters.TransformNamespace(ctx, namepspaces, *cloneConfig, nil)
		if err != nil {
			return err
		}
		for _, o := range filteredListOfNamespaces {

			ns, ok := o.(*corev1.Namespace)
			if !ok {
				return fmt.Errorf("filtering of namespaces failed")
			}

			localCloneConfig := &config.CloneConfig{
				AllApplications: true,
				SourceNamespace: ns.Name,
			}
			validateConfig(ctx, localCloneConfig, args)

			err = exportResources(ns.Name, getOrderedExporters(), *localCloneConfig, client)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
		return nil
	}

	return exportResources(cloneConfig.SourceNamespace, getOrderedExporters(), *cloneConfig, client)
}

func exportResources(namespace string, apiExporters []exporters.ResourceExport, cloneConfig config.CloneConfig, client *kubernetes.Clientset) error {

	var localCache []runtime.Object
	var allResourcesAsYaml []exporters.ResourceYAML

	for _, apiExporter := range apiExporters {

		// TODO: Apply for all namespaces

		//resourceList, err := apiExporter.Get(context.Background(), "rhtap-build-tenant", cloneConfig, client)
		resourceList, err := apiExporter.Get(context.Background(), namespace, cloneConfig, client)

		if err != nil {

			return err
		}
		localCache = append(localCache, resourceList)

		individualResources, err := apiExporter.Transform(context.Background(), resourceList, cloneConfig, localCache)
		if err != nil {
			return err
		}

		resourcesAsYaml, err := apiExporter.GenerateYAML(context.Background(), individualResources)
		if err != nil {
			return err
		}

		allResourcesAsYaml = append(allResourcesAsYaml, resourcesAsYaml...)
	}

	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		return err
	}
	file, err := os.Create("data/" + cloneConfig.OutputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, r := range allResourcesAsYaml {
		_, err := file.WriteString("---\n")
		if err != nil {
			return err
		}
		_, err = file.Write(r.ResourceBytes)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateConfig(ctx context.Context, cloneConfig *config.CloneConfig, args []string) error {

	if cloneConfig.OutputFile == "" {
		(*cloneConfig).OutputFile = time.Now().String()
	}

	// system-wide backup
	if cloneConfig.AllNamespaces {
		(*cloneConfig).AllApplications = true
	} else {
		// app not specified.
		if len(args) == 0 {
			// also, not all applications are being asked for
			if !cloneConfig.AllApplications {
				return fmt.Errorf("application not specified")
			}
		} else {
			(*cloneConfig).ApplicatioName = args[0]
		}

		// if not all namespaces, then this should have been specified
		if cloneConfig.SourceNamespace == "" {
			return fmt.Errorf("source namespace not specified")
		}

		// Use case: System-wide backup and restore.
		if cloneConfig.TargetNamespace == "" {
			(*cloneConfig).TargetNamespace = cloneConfig.SourceNamespace
		}
	}

	return nil
}
