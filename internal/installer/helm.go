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
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

const (
	DefaultKubeConfPath = "~/.kube/config"

	DefaultRepositoryConfig = "~/.config/helm/repositories.yaml"
	DefaultRepositoryCache  = "./helm"

	DefaultWaitTimeoutSec = 60 * 20
	DefaultDurationSec = 10
)

// the type Values used for helm when install release
type Values map[string]interface{}

func (v Values) parse() (string, error) {
	valueBytes, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(valueBytes), nil
}

type Proxy struct {
	ctx             context.Context
	namespace       string
	repositoryCache string
	client          hc.Client // helm client
	logger          glog.Logger
}

func NewProxy(ctx context.Context, namespace string, debug bool, logger glog.Logger) (*Proxy, error) {

	debugLog := func(format string, v ...interface{}) {
		// Change this to your own logger. Default is 'log.Printf(format, v...)'.
	}
	if debug {
		debugLog = func(format string, v ...interface{}) {
			logger.Debug().Msg(fmt.Sprintf(format, v)).Fire()
		}
	}
	opts := &hc.Options{
		Namespace:        namespace, // Change this to the namespace you wish to install the chart in.
		RepositoryCache:  DefaultRepositoryCache,
		RepositoryConfig: DefaultRepositoryConfig,
		Debug:            debug,
		Linting:          true, // Change this to false if you don't want linting.
		DebugLog:         debugLog,
	}

	kubeConf, err := clientcmd.BuildConfigFromFlags("", DefaultKubeConfPath)
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
			repositoryCache: DefaultRepositoryCache,
			client:          c,
			logger:          logger,
		}, nil
	}
	return nil, err
}

func NewKubeClient() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", DefaultKubeConfPath)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func (p *Proxy) install(name, chartName string, v Values, waittingReady bool) error {
	valuesStr, err := v.parse()
	if err != nil {
		p.logger.Error().Error("marshal values to string error", err).Fire()
		return err
	}
	chartSpec := &hc.ChartSpec{
		ReleaseName:     name,
		ChartName:       fmt.Sprintf("%s/%s", p.repositoryCache, chartName),
		Namespace:       p.namespace,
		CreateNamespace: true,
		ValuesYaml:      valuesStr,
		Recreate:        true,
	}

	p.logger.Info().String("Helm install release", name).String("with chart", chartSpec.ChartName).Fire()
	_, err = p.client.InstallOrUpgradeChart(p.ctx, chartSpec)
	if err != nil {
		p.logger.Error().Error("helm install error", err).Fire()
		return err
	}
	if waittingReady {
		err = p.waitingReady(name, DefaultWaitTimeoutSec, DefaultDurationSec)
	}
	return err
}

func (p *Proxy) waitingReady(name string, timeoutSec, durationSec uint64) error {
	p.logger.Info().String("Waiting operator", name).String("in namespace", p.namespace).Msg("ready..").Fire()
	nameField := map[string]string{
		"meta.helm.sh/release-name": name,
	}
	ops := v1.ListOptions{
		FieldSelector: fields.SelectorFromSet(nameField).String(),
	}

	ready := false
	var err error
	duration := time.Duration(durationSec) * time.Second
	kubeClient, err := NewKubeClient()
	if err != nil {
		p.logger.Error().Error("new kube client error", err).Fire()
		return err
	}
	for timeoutSec > 0 {
		ready, err = p.isReady(kubeClient, p.namespace, ops)
		if err != nil {
			p.logger.Error().Error("Check status ready error", err).Fire()
			return err
		}
		if ready {
			p.logger.Info().String("All pods ready of operator", name).String(" in namespace", p.namespace).Fire()
			return nil
		}
		time.Sleep(duration)
		timeoutSec -= durationSec
	}
	return errors.New("waiting operator ready, timeout")
}

func (p *Proxy) isReady(kubeClient *kubernetes.Clientset, namespace string, ops v1.ListOptions) (bool, error) {
	// 获取指定 namespace 中的 Pod 列表信息
	pods, err := kubeClient.CoreV1().Pods(namespace).List(p.ctx, ops)
	if err != nil {
		p.logger.Error().Error("List pod error", err).Fire()
		return false, err
	}
	for _, pod := range pods.Items {
		p.logger.Debug().Any("container status", pod.Status).Fire()
		if pod.Status.Phase == corev1.PodPending {
			p.logger.Warn().String("the status of pod", pod.GetName()).Msg("is pending").Fire()
			return false, nil
		}

		for _, cond := range pod.Status.Conditions {
			if cond.Type == corev1.ContainersReady && cond.Status == corev1.ConditionTrue {
				break
			} else if cond.Status != corev1.ConditionTrue {
				p.logger.Debug().String("pod", pod.GetName()).
					String("is not ready, state", string(cond.Type)).
					String(", reason", cond.Reason).Fire()
				return false, nil
			}
		}
	}
	return true, nil
}
