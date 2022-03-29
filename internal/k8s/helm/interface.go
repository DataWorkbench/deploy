package helm

import (
	"github.com/DataWorkbench/deploy/internal/common"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Chart is the proxy of helm chart that with the Values Configuration.
// helm chart interface
type Chart interface {
	// update each field value from global Config if that is ZERO
	UpdateFromConfig(common.Config) error

	// parse the field-values to Values for helm release
	ParseValues() (Values, error)

	// return chart name
	GetChartName() string

	// return relase name
	GetReleaseName() string

	GetLabels() map[string]string

	// whether to wait release ready
	WaitingReady() bool
	GetTimeoutSecond() int

	InitLocalPvDir() error
}

// helm client interface for dataomnis-service
type Helm interface {
	// install Chart to k8s as a Release
	Install(*Chart) error

	// waiting a release ready
	// param:
	//   releaseName
	//   timeoutSec to waiting
	//   durationSec for checking if release is ready
	WaitingReady(string, int64, int64) error

	// check if a release is ready
	IsReady(v1.ListOptions) (bool, error)

	// upgrade a release with Chart
	Upgrade(*Chart) error

	// delete a release by name(string)
	Delete(string) error

	// delete a release by name(string)
	Exist(string) (bool, error)
}
