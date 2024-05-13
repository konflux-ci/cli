package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/konflux-ci/cli/cmd/rhtap/commands/config"
)

func Test_validateConfig(t *testing.T) {
	type args struct {
		ctx         context.Context
		cloneConfig *config.CloneConfig
		args        []string
	}
	tests := []struct {
		name            string
		args            args
		wantErr         bool
		wantCloneConfig *config.CloneConfig
	}{
		{
			name: "system backup - all namespaces",
			args: args{
				ctx: context.Background(),
				cloneConfig: &config.CloneConfig{
					AllNamespaces:          true,
					AllApplications:        false,
					ApplicatioName:         "",
					SourceNamespace:        "",
					TargetNamespace:        "",
					ComponentSourceURLskip: "",
					OutputDir:              "",
				},
			},
			wantErr: false,
			wantCloneConfig: &config.CloneConfig{
				AllApplications: true,
				ApplicatioName:  "",
				OutputDir:       "something",
			},
		},
		{
			name: "backup for embargo usecase - application not specified",
			args: args{
				ctx: context.Background(),
				cloneConfig: &config.CloneConfig{
					AllApplications: false,
					ApplicatioName:  "",
					AllNamespaces:   false,
					SourceNamespace: "source-namespace",
					TargetNamespace: "target-namespace",

					ComponentSourceURLskip: "",
					OutputDir:              "",
				},
			},
			wantErr: true,
			wantCloneConfig: &config.CloneConfig{
				AllApplications: false,
				ApplicatioName:  "",
				SourceNamespace: "source-namespace",
				TargetNamespace: "target-namespace",
				OutputDir:       "something",
			},
		},
		{
			name: "backup for embargo usecase - true",
			args: args{
				ctx:  context.Background(),
				args: []string{"foo"},
				cloneConfig: &config.CloneConfig{
					AllApplications: false,
					ApplicatioName:  "app",
					AllNamespaces:   false,
					SourceNamespace: "source-namespace",
					TargetNamespace: "target-namespace",

					ComponentSourceURLskip: "",
					OutputDir:              "something",
				},
			},
			wantErr: false,

			// doesn't matter
			wantCloneConfig: &config.CloneConfig{
				AllApplications: false,
				ApplicatioName:  "foo",
				SourceNamespace: "source-namespace",
				TargetNamespace: "target-namespace",
			},
		},
		{
			name: "target ns not specified",
			args: args{
				ctx:  context.Background(),
				args: []string{"foo"},
				cloneConfig: &config.CloneConfig{
					AllApplications: false,
					ApplicatioName:  "app",
					AllNamespaces:   false,
					SourceNamespace: "source-namespace",
					TargetNamespace: "",

					ComponentSourceURLskip: "",
					OutputDir:              "something",
				},
			},
			wantErr: false,

			// doesn't matter
			wantCloneConfig: &config.CloneConfig{
				AllApplications: false,
				ApplicatioName:  "foo",
				SourceNamespace: "source-namespace",
				TargetNamespace: "source-namespace",
			},
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//if err := validateConfig(tt.args.ctx, tt.args.cloneConfig, tt.args.args); (err != nil) != tt.wantErr {
			//	t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			//}

			err := validateConfig(tt.args.ctx, tt.args.cloneConfig, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.NotEmpty(t, tt.args.cloneConfig.OutputDir)
			assert.Equal(t, tt.wantCloneConfig.ApplicatioName, tt.args.cloneConfig.ApplicatioName)
			assert.Equal(t, tt.wantCloneConfig.AllApplications, tt.args.cloneConfig.AllApplications)
			assert.Equal(t, tt.wantCloneConfig.TargetNamespace, tt.args.cloneConfig.TargetNamespace)
			assert.Equal(t, tt.wantCloneConfig.SourceNamespace, tt.args.cloneConfig.SourceNamespace)
		})
	}
}
