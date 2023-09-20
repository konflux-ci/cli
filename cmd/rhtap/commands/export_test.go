package commands

import (
	"reflect"
	"testing"

	rhapAPI "github.com/redhat-appstudio/rhtap-cli/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_generateOverridesMap(t *testing.T) {
	type args struct {
		overrides string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			"basic test",
			args{
				overrides: "comp1=github.com/foo/bar",
			},
			map[string]string{
				"comp1": "github.com/foo/bar",
			},
		},
		{
			"multiple",
			args{
				overrides: "comp1=github.com/foo/bar   comp2=github.com/foo/bar",
			},
			map[string]string{
				"comp1": "github.com/foo/bar",
				"comp2": "github.com/foo/bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateOverridesMap(tt.args.overrides); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateOverridesMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateExportableComponent(t *testing.T) {
	type args struct {
		fetchedComponent rhapAPI.Component
		targetNamespace  string
		overrides        string
	}
	tests := []struct {
		name string
		args args
		want *rhapAPI.Component
	}{
		{
			"impacted component",
			args{
				fetchedComponent: rhapAPI.Component{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "comp1",
						Namespace: "source-namespace",
					},
					Spec: rhapAPI.ComponentSpec{
						ComponentName: "comp1",
						Source: rhapAPI.ComponentSource{
							ComponentSourceUnion: rhapAPI.ComponentSourceUnion{
								GitSource: &rhapAPI.GitSource{
									URL: "github.com/org/repo",
								},
							},
						},
						ContainerImage: "quay.io/org/repo",
					},
				},
				targetNamespace: "target-namespace",
				overrides:       "comp1=github.com/private/comp1 comp2=github.com/private/comp2",
			},
			&rhapAPI.Component{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "target-namespace",
					Annotations: map[string]string{
						"skip-initial-checks":       "true",
						"image.redhat.com/generate": `{"visibility": "public"}`,
					},
				},
				Spec: rhapAPI.ComponentSpec{
					ComponentName: "comp1",
					Source: rhapAPI.ComponentSource{
						ComponentSourceUnion: rhapAPI.ComponentSourceUnion{
							GitSource: &rhapAPI.GitSource{
								URL: "github.com/private/comp1",
							},
						},
					},
				},
			},
		},
		{
			"non-impacted component",
			args{
				fetchedComponent: rhapAPI.Component{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "name",
						Namespace: "source-namespace",
					},
					Spec: rhapAPI.ComponentSpec{
						ComponentName:  "name",
						ContainerImage: "quay.io/org/repo",
					},
				},
				targetNamespace: "target-namespace",
				overrides:       "comp1=github.com/private/comp1 comp2=github.com/private/comp2",
			},
			&rhapAPI.Component{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "target-namespace",
					Annotations: map[string]string{
						"skip-initial-checks": "true",
					},
				},
				Spec: rhapAPI.ComponentSpec{
					ComponentName:  "name",
					ContainerImage: "quay.io/org/repo",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateExportableComponent(tt.args.fetchedComponent, tt.args.targetNamespace, tt.args.overrides); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateExportableComponent() = %v, want %v", got, tt.want)
			}
		})
	}
}
