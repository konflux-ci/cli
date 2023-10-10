package exporters

import (
	"context"
	"reflect"
	"testing"

	rhtapAPI "github.com/redhat-appstudio/rhtap-cli/api/v1alpha1"
	"github.com/redhat-appstudio/rhtap-cli/cmd/rhtap/commands/config"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestTransformComponent(t *testing.T) {
	type args struct {
		ctx                 context.Context
		fetchedResourceList runtime.Object
		cloneConfig         config.CloneConfig
		localCache          []runtime.Object
	}
	tests := []struct {
		name    string
		args    args
		want    []runtime.Object
		wantErr bool
	}{
		{
			name: "golden path",
			args: args{
				ctx: context.Background(),
				fetchedResourceList: &rhtapAPI.ComponentList{
					Items: []rhtapAPI.Component{
						{
							ObjectMeta: v1.ObjectMeta{
								Annotations: map[string]string{
									"image.redhat.com/generate": `
									"image":      "quay.io/redhat-user-workloads/image-controller-system/city-transit/billing",
									"visibility": "public",
									"secret":     "billing"`,
								},
							},
						},
					},
				},
				cloneConfig: config.CloneConfig{
					TargetNamespace: "foo",
				},
			},
			want: []runtime.Object{
				&rhtapAPI.Component{
					ObjectMeta: v1.ObjectMeta{
						Namespace: "foo",
						Annotations: map[string]string{
							"skip-initial-checks":       "true",
							"image.redhat.com/generate": `{"visibility": "public"}`,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransformComponent(tt.args.ctx, tt.args.fetchedResourceList, tt.args.cloneConfig, tt.args.localCache)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformComponent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransformComponent() = %v, want %v", got, tt.want)
			}
		})
	}
}
