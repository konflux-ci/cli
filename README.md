## RHTAP CLI
Command line for working with your Applications on RHTAP

## Setup
- `git clone https://github.com/redhat-appstudio/rhtap-cli.git` 
- `cd rhtap-cli`
- `go build -o rhtap cmd/rhtap/main.go`


## Usage

```
Usage:
  rhtap export application [flags]

Flags:
  -f, --from string             Namespace from which the Application is being cloned.
  -h, --help                    help for application
  -o, --overrides string        Overwrite the source code url for specific components
  -s, --skip string             List of components to be skipped
  -t, --to string               Namespace to which the Application is being cloned.
  -w, --write-to string         Local filesyste path where the YAML would be written out to.

```

### Steps

1. Login to your RHTAP namespace using `kubelogin` or `kubectl`.

2. Generate YAML for the `Application` and associated resources ( `Components`,`IntegrationTestScenarios`, etc )


```
./rhtap export application njtransit \
--from corp-shared-tenant \
--to shbose-tenant \
--overrides  "billing-service=github.com/dave/private-repo" \
--write-to /Users/dave/exported-njtransit.yaml

```

3. Switch to the target namespace  

```
kubectl apply -f exported-njtransit.yaml
```
