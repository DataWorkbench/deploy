package helm

import (
	"context"
	"fmt"
	"github.com/DataWorkbench/deploy/internal/common"
	"github.com/DataWorkbench/deploy/internal/k8s"
	"github.com/DataWorkbench/glog"
	hc "github.com/mittwald/go-helm-client"
	"github.com/pkg/errors"
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

// ******************************************************************
// helm client proxy to handle helm release, implement Helm interface
// ******************************************************************
type Proxy struct {
	namespace       string
	repositoryCache string
	client          hc.Client // helm client
	logger          *glog.Logger
	kclient         *k8s.KClient
}

func NewProxy(ctx context.Context, namespace string, logger *glog.Logger) (*Proxy, error) {
	debugLog := func(format string, v ...interface{}) {
		// Change this to your own logger. Default is 'log.Printf(format, v...)'.
	}
	if common.Debug {
		debugLog = func(format string, v ...interface{}) {
			logger.Debug().Msg(fmt.Sprintf(format, v)).Fire()
		}
	}
	kubeConfPath := fmt.Sprintf(k8s.DefaultKubeConfFmt, os.Getenv("HOME"))
	HelmRepoConf := fmt.Sprintf(DefaultHelmRepositoryConfigFmt, os.Getenv("HOME"))
	HelmRepoCache := fmt.Sprintf(DefaultHelmRepositoryCacheFmt, os.Getenv("HOME"))
	opts := &hc.Options{
		Namespace:        namespace, // Change this to the namespace you wish to install the chart in.
		RepositoryCache:  HelmRepoCache,
		RepositoryConfig: HelmRepoConf,
		Debug:            common.Debug,
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
			namespace:       namespace,
			repositoryCache: HelmRepoCache,
			client:          c,
			logger:          logger,
		}, nil
	}
	return nil, err
}

func (p Proxy) Install(ctx context.Context, chart Chart) error {
	var name = chart.GetReleaseName()
	var chartName = chart.GetChartName()

	values, err := chart.ParseValues()
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
	if p.kclient, err = k8s.NewKClient(); err != nil {
		p.logger.Error().Error("new kube client error", err).Fire()
		return err
	}
	if err = p.kclient.CreateNamespace(ctx, p.namespace); err != nil {
		p.logger.Error().Error("create namespace error", err).Fire()
		return err
	}

	p.logger.Info().String("helm install release", name).String("with chart", chartName).Fire()
	p.logger.Debug().Any("values", values).Fire()
	chartSpec := &hc.ChartSpec{
		ReleaseName: name,
		ChartName:   fmt.Sprintf("%s/%s", p.repositoryCache, chartName),
		Namespace:   p.namespace,
		DryRun:      common.DryRun,
		ValuesYaml:  valuesStr,
		Recreate:    true,
	}
	_, err = p.client.InstallOrUpgradeChart(ctx, chartSpec)
	if err != nil {
		p.logger.Error().Error("helm install error", err).Fire()
		return err
	}

	if chart.WaitingReady() && !common.DryRun{
		wCtx, cancel := context.WithTimeout(ctx, chart.GetTimeoutSecond())
		defer cancel()
		err = p.WaitingReady(wCtx, chart)
	}
	return err
}

func (p Proxy) WaitingReady(ctx context.Context, chart Chart) error {
	name := chart.GetReleaseName()
	p.logger.Info().String("waiting release", name).Msg("ready..").Fire()

	labelMap := chart.GetLabels()
	ops := v1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelMap).String(),
	}

	ready := false
	var err error
	p.kclient, err = k8s.NewKClient()
	if err != nil {
		p.logger.Error().Error("new kube client error", err).Fire()
		return err
	}

	duration := time.Duration(DefaultDurationSec) * time.Second
	for {
		select {
		case <- time.After(duration):
			ready, err = p.IsReady(ctx, ops)
			if err != nil {
				p.logger.Error().Error("check status ready error", err).Fire()
				return err
			}
			if ready {
				p.logger.Info().String("all pods ready of release", name).
					String("in namespace", p.namespace).Fire()
				return nil
			}
		case <- ctx.Done():
			p.logger.Warn().Error("waiting-action been canceled, error", ctx.Err()).Fire()
			return nil
		}
	}
}

// Note: need to init p.kubeClient before
func (p Proxy) IsReady(ctx context.Context, ops v1.ListOptions) (bool, error) {
	// get PodLists
	pods, err := p.kclient.CoreV1().Pods(p.namespace).List(ctx, ops)
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

func (p Proxy) Exist(releaseName string) (bool, error) {
	_, err := p.client.GetRelease(releaseName)
	if err != nil {
		if errors.Cause(err).Error() == ReleaseNotFoundErr { // release not found
			return false, nil
		}
		return false, err
	}
	return true, err
}
