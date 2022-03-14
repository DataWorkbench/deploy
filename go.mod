module github.com/DataWorkbench/deploy

go 1.15

require (
	github.com/DataWorkbench/common v0.0.0-20220308091550-bc27f5176476
	github.com/DataWorkbench/glog v0.0.0-20220302035436-25a1ae256704
	github.com/go-playground/validator/v10 v10.4.1
	github.com/mittwald/go-helm-client v0.8.0
	github.com/pkg/errors v0.9.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0
)
