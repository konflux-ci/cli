package exporters

import (
	"context"
	"reflect"
	"testing"

	konfluxAPI "github.com/konflux-ci/cli/api/v1alpha1"
	"github.com/konflux-ci/cli/cmd/konflux/commands/config"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestTransformEnvironment(t *testing.T) {
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
				cloneConfig: config.CloneConfig{
					SourceNamespace: "source-ns",
					TargetNamespace: "target-ns",
				},
				fetchedResourceList: &konfluxAPI.EnvironmentList{
					Items: []konfluxAPI.Environment{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "production",
								Namespace: "source-ns",
							},
						},
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "development",
								Namespace: "source-ns",
							},
						},
					},
				},
			},
			want: []runtime.Object{
				&konfluxAPI.Environment{
					ObjectMeta: v1.ObjectMeta{
						Name:      "production",
						Namespace: "target-ns",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransformEnvironment(tt.args.ctx, tt.args.fetchedResourceList, tt.args.cloneConfig, tt.args.localCache)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformEnvironment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransformEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}
