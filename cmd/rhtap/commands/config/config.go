package config

type CloneConfig struct {
	AllApplications             bool
	AllNamespaces               bool
	ApplicatioName              string
	SourceNamespace             string
	TargetNamespace             string
	ComponentSourceURLOverrides string
	ComponentSourceURLskip      string
	OutputFile                  string
}
