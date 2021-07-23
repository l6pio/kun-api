package core

import (
	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	DatabaseName string
	MongoAddr    string
	MongoUser    string
	MongoPass    string
	KubeClient   *kubernetes.Clientset
}

type PodStatus int

const (
	PodCreate PodStatus = 1
	PodDelete PodStatus = 0
)

type PodEvent struct {
	Timestamp int64
	ImageId   string
	Image     string
	Status    PodStatus
}

func NewConfig(mongoAddr, mongoUser, mongoPass, master, kubeconfig string) (*Config, error) {
	k8sConfig, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	//k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}

	return &Config{
		DatabaseName: "kun",
		MongoAddr:    mongoAddr,
		MongoUser:    mongoUser,
		MongoPass:    mongoPass,
		KubeClient:   kubeClient,
	}, nil
}

func WithConfig(conf *Config) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(ctx echo.Context) error {
				ctx.Set("config", conf)
				return next(ctx)
			}
		},
	}
}
