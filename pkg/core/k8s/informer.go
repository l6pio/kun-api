package k8s

import (
	"github.com/labstack/gommon/log"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db/vo"
	"l6p.io/kun/api/pkg/core/k8s/queue"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var shutdownSignals = []os.Signal{syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT}

type PodController struct {
	conf               *core.Config
	podQueue           *queue.PodQueue
	imageTimelineQueue *queue.ImageTimelineQueue
}

func (c *PodController) Add(obj interface{}) {
	pod, ok := obj.(*v1.Pod)
	if !ok {
		return
	}
	klog.V(4).InfoS("Adding pod", "pod", klog.KObj(pod))
	c.podQueue.Push(pod)
	c.imageTimelineQueue.Push(pod, vo.ImageUp)
}

func (c *PodController) Delete(obj interface{}) {
	pod, ok := obj.(*v1.Pod)
	if !ok {
		return
	}
	klog.V(4).InfoS("Deleting pod", "pod", klog.KObj(pod))
	c.podQueue.Push(pod)
	c.imageTimelineQueue.Push(pod, vo.ImageDown)
}

func StartPodInformer(conf *core.Config) {
	informerFactory := informers.NewSharedInformerFactory(conf.KubeClient, 0)
	podInformer := informerFactory.Core().V1().Pods()

	c := &PodController{
		conf:               conf,
		podQueue:           queue.NewPodQueue(conf),
		imageTimelineQueue: queue.NewImageTimelineQueue(conf),
	}

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.Add,
		UpdateFunc: func(interface{}, interface{}) {},
		DeleteFunc: c.Delete,
	})

	stopCh := setupSignalHandler()
	go informerFactory.Start(stopCh)

	defer runtime.HandleCrash()
	defer c.podQueue.ShutDown()
	defer c.imageTimelineQueue.ShutDown()

	log.Info("start listening for Pod events")
	defer log.Info("stop listening to Pod events")

	if !cache.WaitForCacheSync(stopCh, podInformer.Informer().HasSynced) {
		return
	}

	for i := 0; i < 5; i++ {
		go wait.Until(c.podQueue.Worker, time.Second, stopCh)
		go wait.Until(c.imageTimelineQueue.Worker, time.Second, stopCh)
	}
	<-stopCh
}

func setupSignalHandler() (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	ch := make(chan os.Signal, 3)
	signal.Notify(ch, shutdownSignals...)
	go func() {
		<-ch
		close(stop)
		os.Exit(1)
	}()
	return stop
}
