package exporters

import (
	"context"

	"github.com/konflux-ci/cli/cmd/konflux/commands/config"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

type ResourceExport struct {
	Transform    func(ctx context.Context, fetchedResourceList runtime.Object, cloneConfig config.CloneConfig, localCache []runtime.Object) ([]runtime.Object, error)
	GenerateYAML func(ctx context.Context, transformedResources []runtime.Object) ([][]byte, error)
	Get          func(ctx context.Context, namespace string, cloneConfig config.CloneConfig, client *kubernetes.Clientset) (runtime.Object, error)
	Sensitive    bool
}
