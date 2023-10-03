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
  -a, --all-applications   When set, all Applications in the current namespace will be cloned.
  -p, --all-projects       When set, all namespaces/projects will be cloned.
  -f, --from string        Namespace from which the Application is being cloned.
  -h, --help               help for application
  -k, --key string         Local filesystem path to an existing encryption key
  -s, --skip string        List of components to be skipped
  -t, --to string          Namespace to which the Application is being cloned.
  -w, --write-to string    Local filesystem directory path where the YAML would be written out to.
```

### Examples


Unless otherwise specified by the `--write-to` flag, the exported YAML files would be written out as follows:

* `data/20231003120120/source-tenant.yaml`  
* `data/20231003120120/encrypted-source-tenant.yaml`. 

#### Export all Applications across all Namespaces

1. Login to your RHTAP namespace using `kubelogin` or `kubectl`.

2.  ```
    ./rhtap export application --all-projects
    ```

#### Export/Clone one Application from one Namespace to another

1. Login to your RHTAP namespace using `kubelogin` or `kubectl`.

2.  ```
    ./rhtap export application build  --from rhtap-build-tenant --to shbose-tenant
    ```


#### Encrypt sensitive data extracted from one or more namespaces


Login to your RHTAP namespace using `kubelogin` or `kubectl`.

2.  ```
    ./rhtap export application build  -f rhtap-build-tenant -k /Users/sbose/keys/sbose.pub 
    ```

