package service

import (
	"github.com/labstack/gommon/log"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var shutdownSignals = []os.Signal{syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT}

type PodController struct {
	conf  *core.Config
	queue workqueue.RateLimitingInterface
}

func (c *PodController) Add(obj interface{}) {
	pod, ok := obj.(*v1.Pod)
	if !ok {
		return
	}
	klog.V(4).InfoS("Adding pod", "pod", klog.KObj(pod))
	c.pushEvent(pod, core.PodCreate)
}

func (c *PodController) Delete(obj interface{}) {
	pod, ok := obj.(*v1.Pod)
	if !ok {
		return
	}
	klog.V(4).InfoS("Deleting pod", "pod", klog.KObj(pod))
	c.pushEvent(pod, core.PodDelete)
}

func (c *PodController) pushEvent(pod *v1.Pod, status core.PodStatus) {
	namespace := pod.Namespace
	podName := pod.Name

	images := make(map[string]string)
	for _, container := range pod.Status.ContainerStatuses {
		images[container.ImageID] = container.Image
	}

	timestamp := time.Now().UnixNano() / 1e6
	for imageId, image := range images {
		c.queue.Add(
			&core.PodEvent{
				Timestamp: timestamp,
				Namespace: namespace,
				PodName:   podName,
				ImageId:   imageId,
				Image:     image,
				Status:    status,
			},
		)
	}
}

func (c *PodController) worker() {
	for c.processNextWorkItem() {
	}
}

func (c *PodController) processNextWorkItem() bool {
	obj, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(obj)

	podEvent := obj.(*core.PodEvent)

	if podEvent.Status == core.PodCreate {
		utilruntime.HandleError(processPodCreate(c.conf, podEvent))
	} else {
		utilruntime.HandleError(processPodDelete(c.conf, podEvent))
	}
	return true
}

func processPodCreate(conf *core.Config, event *core.PodEvent) error {
	exists, err := db.ImageExists(conf, event.ImageId)
	if err != nil {
		return err
	}

	if !exists {
		report, err := Scan(event.Image)
		if err != nil {
			return err
		}

		if len(report.Matches) == 0 {
			log.Info("no vulnerabilities found")
		}

		log.Info("start saving scan results")
		Insert(conf, event.ImageId, report)
		log.Infof("scan results of '%v' has been saved", event.ImageId)
	}

	if err := db.SavePodEvent(conf, event); err != nil {
		return err
	}
	return db.UpdateImagePods(conf, event.ImageId, core.PodCreate)
}

func processPodDelete(conf *core.Config, event *core.PodEvent) error {
	if err := db.SavePodEvent(conf, event); err != nil {
		return err
	}
	return db.UpdateImagePods(conf, event.ImageId, core.PodDelete)
}

func StartPodInformer(conf *core.Config) {
	informerFactory := informers.NewSharedInformerFactory(conf.KubeClient, 0)
	podInformer := informerFactory.Core().V1().Pods()

	c := &PodController{
		conf: conf,
		queue: workqueue.NewNamedRateLimitingQueue(
			workqueue.DefaultControllerRateLimiter(),
			"kun",
		),
	}

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.Add,
		UpdateFunc: func(interface{}, interface{}) {},
		DeleteFunc: c.Delete,
	})

	stopCh := setupSignalHandler()
	go informerFactory.Start(stopCh)

	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	log.Info("start listening for Pod events")
	defer log.Info("stop listening to Pod events")

	if !cache.WaitForCacheSync(stopCh, podInformer.Informer().HasSynced) {
		return
	}

	for i := 0; i < 5; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
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
