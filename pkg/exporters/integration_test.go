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

func TestTransformIntegrationTestScenario(t *testing.T) {
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
			name: "one integration test",
			args: args{
				ctx: context.Background(),
				cloneConfig: config.CloneConfig{
					SourceNamespace: "source-namespace",
					TargetNamespace: "target-namespace",
					AllApplications: false,
					ApplicatioName:  "source-app",
				},
				fetchedResourceList: &rhtapAPI.IntegrationTestScenarioList{
					Items: []rhtapAPI.IntegrationTestScenario{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "t1",
								Namespace: "source-namespace",
							},
							Spec: rhtapAPI.IntegrationTestScenarioSpec{
								Application: "not-that-app",
							},
						},
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "t2",
								Namespace: "source-namespace",
							},
							Spec: rhtapAPI.IntegrationTestScenarioSpec{
								Application: "source-app",
							},
						},
					},
				},
			},
			want: []runtime.Object{
				&rhtapAPI.IntegrationTestScenario{
					ObjectMeta: v1.ObjectMeta{
						Name:      "t2",
						Namespace: "target-namespace",
					},
					Spec: rhtapAPI.IntegrationTestScenarioSpec{
						Application: "source-app",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransformIntegrationTestScenario(tt.args.ctx, tt.args.fetchedResourceList, tt.args.cloneConfig, tt.args.localCache)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformIntegrationTestScenario() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) || (len(got) == 0 && len(tt.want) == 0) {
				t.Errorf("TransformIntegrationTestScenario() = %v, want %v, len(got) %d, len(want) %d", got, tt.want, len(got), len(tt.want))
			}
		})
	}
}
