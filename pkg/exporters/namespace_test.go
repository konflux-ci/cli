package exporters

import (
	"context"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"

	"github.com/redhat-appstudio/rhtap-cli/cmd/rhtap/commands/config"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestTransformNamespace(t *testing.T) {
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
			name: "not an appstudio ns",
			args: args{
				ctx: context.Background(),
				fetchedResourceList: &corev1.NamespaceList{
					Items: []corev1.Namespace{
						{
							ObjectMeta: v1.ObjectMeta{
								Name: "n1",
								Labels: map[string]string{
									"toolchain.dev.openshift.com/tier": "appstudio",
								},
							},
						},
					},
				},
			},
			want: []runtime.Object{
				&corev1.Namespace{
					ObjectMeta: v1.ObjectMeta{
						Name: "n1",
						Labels: map[string]string{
							"toolchain.dev.openshift.com/tier": "appstudio",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransformNamespace(tt.args.ctx, tt.args.fetchedResourceList, tt.args.cloneConfig, tt.args.localCache)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransformNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}
