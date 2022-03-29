package installer

import (
	"context"
	"fmt"
	"github.com/DataWorkbench/glog"
	hc "github.com/mittwald/go-helm-client"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"time"
)

const (
	DefaultHelmRepositoryConfigFmt = "%s/.config/helm/repositories.yaml"
	DefaultHelmRepositoryCacheFmt  = "%s/.cache/helm/repository"

	DefaultDurationSec = 20

	ReleaseNotFoundErr = "release: not found"
)

// **************************************************
// the type Values used for helm when install release
// **************************************************
type Values map[string]interface{}

func (v Values) parse() (string, error) {
	valueBytes, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(valueBytes), nil
}

func (v Values) isEmpty() bool {
	return len(v) == 0
}

// ******************************************************************
// helm client proxy to handle helm release, implement Helm interface
// ******************************************************************
type Proxy struct {
	ctx             context.Context
	namespace       string
	repositoryCache string
	client          hc.Client // helm client
	logger          *glog.Logger
	kclient         *KClient
}

func NewProxy(ctx context.Context, namespace string, logger *glog.Logger, debug bool) (*Proxy, error) {

	debugLog := func(format string, v ...interface{}) {
		// Change this to your own logger. Default is 'log.Printf(format, v...)'.
	}
	if debug {
		debugLog = func(format string, v ...interface{}) {
			logger.Debug().Msg(fmt.Sprintf(format, v)).Fire()
		}
	}
	kubeConfPath := fmt.Sprintf(DefaultKubeConfFmt, os.Getenv("HOME"))
	HelmRepoConf := fmt.Sprintf(DefaultHelmRepositoryConfigFmt, os.Getenv("HOME"))
	HelmRepoCache := fmt.Sprintf(DefaultHelmRepositoryCacheFmt, os.Getenv("HOME"))
	opts := &hc.Options{
		Namespace:        namespace, // Change this to the namespace you wish to install the chart in.
		RepositoryCache:  HelmRepoCache,
		RepositoryConfig: HelmRepoConf,
		Debug:            debug,
		Linting:          true, // Change this to false if you don't want linting.
		DebugLog:         debugLog,
	}

	kubeConf, err := clientcmd.BuildConfigFromFlags("", kubeConfPath)
	if err != nil {
		return nil, err
	}
	restConfopts := &hc.RestConfClientOptions{
		Options:    opts,
		RestConfig: kubeConf,
	}
	var c hc.Client
	if c, err = hc.NewClientFromRestConf(restConfopts); err == nil {
		return &Proxy{
			ctx:             ctx,
			namespace:       namespace,
			repositoryCache: HelmRepoCache,
			client:          c,
			logger:          logger,
		}, nil
	}
	return nil, err
}

func (p *Proxy) install(chart Chart) error {
	var name = chart.getReleaseName()
	var chartName = chart.getChartName()

	values, err := chart.parseValues()
	if err != nil {
		p.logger.Error().String("chart with name", chartName).Error("parse values error", err).Fire()
		return err
	}
	var valuesStr string
	if !values.isEmpty() { // parse values to str
		valuesStr, err = values.parse()
		if err != nil {
			p.logger.Error().Error("marshal values to string error", err).Fire()
			return err
		}
	}

	p.logger.Info().String("create namespace", p.namespace).Fire()
	if p.kclient, err = NewKClient(); err != nil {
		p.logger.Error().Error("new kube client error", err).Fire()
		return err
	}
	if err = p.kclient.CreateNamespace(p.ctx, p.namespace); err != nil {
		p.logger.Error().Error("create namespace error", err).Fire()
		return err
	}

	p.logger.Info().String("helm install release", name).String("with chart", chartName).Fire()
	p.logger.Debug().Any("values", values).Fire()
	chartSpec := &hc.ChartSpec{
		ReleaseName: name,
		ChartName:   fmt.Sprintf("%s/%s", p.repositoryCache, chartName),
		Namespace:   p.namespace,
		ValuesYaml:  valuesStr,
		Recreate:    true,
	}
	_, err = p.client.InstallOrUpgradeChart(p.ctx, chartSpec)
	if err != nil {
		p.logger.Error().Error("helm install error", err).Fire()
		return err
	}

	if chart.waitingReady() {
		err = p.waitingReady(chart)
	}
	return err
}

func (p *Proxy) waitingReady(chart Chart) error {
	name := chart.getReleaseName()
	p.logger.Info().String("waiting release", name).Msg("ready..").Fire()

	labelMap := chart.getLabels()
	ops := v1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelMap).String(),
	}

	ready := false
	var err error
	p.kclient, err = NewKClient()
	if err != nil {
		p.logger.Error().Error("new kube client error", err).Fire()
		return err
	}

	timeoutSec := chart.getTimeoutSecond()
	duration := time.Duration(DefaultDurationSec) * time.Second
	spendSecs := 0
	for spendSecs < timeoutSec {
		time.Sleep(duration)
		spendSecs += DefaultDurationSec

		ready, err = p.isReady(ops)
		if err != nil {
			p.logger.Error().Error("check status ready error", err).Fire()
			return err
		}
		if ready {
			p.logger.Info().String("all pods ready of release", name).
				String("in namespace", p.namespace).Int("spend seconds", spendSecs).Fire()
			return nil
		}
	}
	msg := fmt.Sprintf("install release=%s failed, timeout.", name)
	p.logger.Error().Msg(msg).Fire()
	return errors.New(msg)
}

// Note: need to init p.kubeClient before
func (p *Proxy) isReady(ops v1.ListOptions) (bool, error) {
	// get PodLists
	pods, err := p.kclient.CoreV1().Pods(p.namespace).List(p.ctx, ops)
	if err != nil {
		return false, err
	}
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodSucceeded {
			continue
		}

		for _, condition := range pod.Status.Conditions {
			if condition.Status != corev1.ConditionTrue {
				p.logger.Info().String("pod", pod.GetName()).
					String("is not ready, status of conditionType", string(condition.Type)).
					String("is not true, reason", condition.Reason).
					String("message", condition.Message).Fire()
				return false, nil
			}
		}
	}
	return true, nil
}

func (p *Proxy) exist(releaseName string) (bool, error) {
	_, err := p.client.GetRelease(releaseName)
	if err != nil {
		if errors.Cause(err).Error() == ReleaseNotFoundErr { // release not found
			return false, nil
		}
		return false, err
	}
	return true, err
}
