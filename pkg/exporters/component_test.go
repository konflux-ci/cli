package exporters

import (
	"context"
	"reflect"
	"testing"

	rhtapAPI "github.com/konflux-ci/cli/api/v1alpha1"
	"github.com/konflux-ci/cli/cmd/rhtap/commands/config"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestTransformComponent(t *testing.T) {
	gitSource := rhtapAPI.ComponentSource{}
	gitSource.GitSource = &rhtapAPI.GitSource{
		URL: "https://github.com/a/b",
	}

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
			name: "embargo with 3 components",
			want: []runtime.Object{
				&rhtapAPI.Component{
					ObjectMeta: v1.ObjectMeta{
						Namespace: "foo",
						Name:      "c1",
						Annotations: map[string]string{
							"skip-initial-checks": "true",
						},
					},
					Spec: rhtapAPI.ComponentSpec{
						Application: "app-name",
					},
				},
				&rhtapAPI.Component{
					// to be imported as an image.
					ObjectMeta: v1.ObjectMeta{
						Name:      "c2",
						Namespace: "foo",
						Annotations: map[string]string{
							"skip-initial-checks": "true",
						},
					},
					Spec: rhtapAPI.ComponentSpec{
						Application: "app-name",
					},
				},
			},
			wantErr: false,
			args: args{
				ctx: context.Background(),
				cloneConfig: config.CloneConfig{
					TargetNamespace:        "foo",
					ApplicatioName:         "app-name",
					AllApplications:        false,
					ComponentSourceURLskip: "https://github.com/a/b",
					AsPrebuiltImages:       true,
				},
				fetchedResourceList: &rhtapAPI.ComponentList{
					Items: []rhtapAPI.Component{
						{
							// to be imported as image
							ObjectMeta: v1.ObjectMeta{
								Name:      "c1",
								Namespace: "source-ns",
							},
							Spec: rhtapAPI.ComponentSpec{
								Application: "app-name",
							},
						},
						{
							// to be skipped
							ObjectMeta: v1.ObjectMeta{
								Name:      "c3",
								Namespace: "source-ns",
							},
							Spec: rhtapAPI.ComponentSpec{
								Application: "app-name",
								Source:      gitSource,
							},
						},
						{
							// to be imported as an image.
							ObjectMeta: v1.ObjectMeta{
								Name:      "c2",
								Namespace: "source-ns",
							},
							Spec: rhtapAPI.ComponentSpec{
								Application: "app-name",
							},
						},
					},
				},
			},
		},
		{
			name: "golden path",
			args: args{
				ctx: context.Background(),
				cloneConfig: config.CloneConfig{
					TargetNamespace:        "foo",
					ApplicatioName:         "app-name",
					AllApplications:        false,
					ComponentSourceURLskip: "https://github.com/a/b",
				},
				fetchedResourceList: &rhtapAPI.ComponentList{
					Items: []rhtapAPI.Component{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "c1",
								Namespace: "source-ns",
							},
							Spec: rhtapAPI.ComponentSpec{
								Application: "app-name",
							},
						},
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "c3",
								Namespace: "source-ns",
							},
							Spec: rhtapAPI.ComponentSpec{
								Application: "app-name",
								Source:      gitSource,
							},
						},
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "c2",
								Namespace: "source-ns",
							},
							Spec: rhtapAPI.ComponentSpec{
								Application: "not-app-name",
							},
						},
					},
				},
			},
			want: []runtime.Object{
				&rhtapAPI.Component{
					ObjectMeta: v1.ObjectMeta{
						Namespace: "foo",
						Name:      "c1",
						Annotations: map[string]string{
							"skip-initial-checks":       "true",
							"image.redhat.com/generate": `{"visibility": "public"}`,
						},
					},
					Spec: rhtapAPI.ComponentSpec{
						Application: "app-name",
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

func Test_shouldSkip(t *testing.T) {
	type args struct {
		listOfURLsToBeSkipped string
		sourceCodeURL         string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should skip - true",
			args: args{
				listOfURLsToBeSkipped: "github.com/foo/bar,github.com/a/b",
				sourceCodeURL:         "github.com/a/b",
			},
			want: true,
		},
		{
			name: "should skip - false",
			args: args{
				listOfURLsToBeSkipped: "",
				sourceCodeURL:         "github.com/a/b",
			},
			want: false,
		},
		{
			name: "should skip - false",
			args: args{
				listOfURLsToBeSkipped: "github.com/foo/bar,github.com/foo/bar2",
				sourceCodeURL:         "github.com/a/b",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldSkip(tt.args.listOfURLsToBeSkipped, tt.args.sourceCodeURL); got != tt.want {
				t.Errorf("shouldSkip() = %v, want %v", got, tt.want)
			}
		})
	}
}
