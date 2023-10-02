## RHTAP CLI
Command line for working with your Applications on RHTAP

## Setup
- `git clone https://github.com/redhat-appstudio/rhtap-cli.git` 
- `cd rhtap-cli`
- `go build -o rhtap cmd/rhtap/main.go`


## Usage

```
./rhtap export application --help  
Usage:
  rhtap export application [flags]

Flags:
  -a, --all-applications   When set, all Applications in the current namespace will be cloned.
  -p, --all-projects       When set, all namespaces/projects will be cloned.
  -f, --from string        Namespace from which the Application is being cloned.
  -h, --help               help for application
  -s, --skip string        List of components to be skipped
  -t, --to string          Namespace to which the Application is being cloned.
  -w, --write-to string    Local filesystem path where the YAML would be written out to.
```

### Examples

#### Export all Applications across all Namespaces

1. Login to your RHTAP namespace using `kubelogin` or `kubectl`.

2.  ```
    ./rhtap export application --all-projects
    ```

A directory named `data` will be created by the tool, inside which exported YAML files
would be written out to.


#### Export/Clone one Application from one Namespace to another

1. Login to your RHTAP namespace using `kubelogin` or `kubectl`.

2.  ```
    ./rhtap export application build  --from rhtap-build-tenant --to shbose-tenant --write-to output.yaml
    ```

