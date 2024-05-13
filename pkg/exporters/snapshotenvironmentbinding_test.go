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

func TestTransformSnapshotEnvironmentBindings(t *testing.T) {
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
			name: "for a single ns",
			args: args{
				ctx: context.Background(),
				fetchedResourceList: &konfluxAPI.SnapshotEnvironmentBindingList{
					Items: []konfluxAPI.SnapshotEnvironmentBinding{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "seb1",
								Namespace: "source-namespace",
							},
							Spec: konfluxAPI.SnapshotEnvironmentBindingSpec{
								Application: "source-application",
							},
						},
					},
				},
				cloneConfig: config.CloneConfig{
					AllApplications: false,
					SourceNamespace: "source-namespace",
					TargetNamespace: "target-namespace",
					ApplicatioName:  "source-application",
				},
			},
			wantErr: false,
			want: []runtime.Object{
				&konfluxAPI.SnapshotEnvironmentBinding{
					ObjectMeta: v1.ObjectMeta{
						Name:      "seb1",
						Namespace: "target-namespace",
					},
					Spec: konfluxAPI.SnapshotEnvironmentBindingSpec{
						Application: "source-application",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransformSnapshotEnvironmentBindings(tt.args.ctx, tt.args.fetchedResourceList, tt.args.cloneConfig, tt.args.localCache)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformSnapshotEnvironmentBindings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransformSnapshotEnvironmentBindings() = %v, want %v", got, tt.want)
			}
		})
	}
}
