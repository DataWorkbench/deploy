package installer

import v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Chart is the proxy of helm chart that with the Values Configuration.
// helm chart interface
type Chart interface {
	// update each field value from global Config if that is ZERO
	updateFromConfig(Config) error

	// parse the field-values to Values for helm release
	parseValues() (Values, error)

	// return chart name
	getChartName() string

	// return relase name
	getReleaseName() string

	// whether to wait release ready
	waitingReady() bool


}


// helm client interface for dataomnis-service
type Helm interface {
	// install Chart to k8s as a Release
	install(*Chart) error

	// waiting a release ready
	// param:
	//   releaseName
	//   timeoutSec to waiting
	//   durationSec for checking if release is ready
	waitingReady(string, int64, int64) error

	// check if a release is ready
	isReady(v1.ListOptions)(bool, error)

	// upgrade a release with Chart
	upgrade(*Chart) error

	// delete a release by name(string)
	delete(string) error
}
