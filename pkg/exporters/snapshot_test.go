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

func TestTransformSnapshots(t *testing.T) {
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
			name: "snapshots by application",
			args: args{
				ctx: context.Background(),
				fetchedResourceList: &konfluxAPI.SnapshotList{
					Items: []konfluxAPI.Snapshot{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "s1",
								Namespace: "source-namepsace",
							},
							Spec: konfluxAPI.SnapshotSpec{
								Application: "source-application",
							},
						},
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "s2",
								Namespace: "source-namepsace",
							},
							Spec: konfluxAPI.SnapshotSpec{
								Application: "not-source-application",
							},
						},
					},
				},
				cloneConfig: config.CloneConfig{
					ApplicatioName:  "source-application",
					AllApplications: false,
					SourceNamespace: "source-namespace",
					TargetNamespace: "target-namespace",
				},
			},
			wantErr: false,
			want: []runtime.Object{
				&konfluxAPI.Snapshot{
					ObjectMeta: v1.ObjectMeta{
						Name:      "s1",
						Namespace: "target-namepsace",
					},
					Spec: konfluxAPI.SnapshotSpec{
						Application: "source-application",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransformSnapshots(tt.args.ctx, tt.args.fetchedResourceList, tt.args.cloneConfig, tt.args.localCache)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformSnapshots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// TODO :fix failure
			if !reflect.DeepEqual(got[0], tt.want[0]) {
				t.Errorf("TransformSnapshots() = %v, want %v", got[0], tt.want[0])
			}
		})
	}
}
