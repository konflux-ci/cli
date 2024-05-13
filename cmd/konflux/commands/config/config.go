package config

type CloneConfig struct {
	AllApplications        bool
	AllNamespaces          bool
	ApplicatioName         string
	SourceNamespace        string
	TargetNamespace        string
	ComponentSourceURLskip string
	OutputDir              string
	Key                    string
	AsPrebuiltImages       bool
}
