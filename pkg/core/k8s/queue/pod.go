package queue

import (
	v1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/util/workqueue"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db"
	"l6p.io/kun/api/pkg/core/db/vo"
	"time"
)

type PodQueue struct {
	conf  *core.Config
	queue workqueue.RateLimitingInterface
}

func NewPodQueue(conf *core.Config) *PodQueue {
	return &PodQueue{
		conf: conf,
		queue: workqueue.NewNamedRateLimitingQueue(
			workqueue.DefaultControllerRateLimiter(),
			"PodQueue",
		),
	}
}

func (p *PodQueue) ShutDown() {
	p.ShutDown()
}

func (p *PodQueue) Push(pod *v1.Pod) {
	timestamp := time.Now().UnixNano() / 1e6

	var statue string
	for _, container := range pod.Status.ContainerStatuses {
		if container.State.Terminated != nil {
			statue = container.State.Terminated.Reason
		} else if container.State.Waiting != nil {
			statue = container.State.Waiting.Reason
		} else {
			statue = "Running"
		}
	}

	p.queue.Add(
		&vo.Pod{
			Timestamp: timestamp,
			Namespace: pod.Namespace,
			Name:      pod.Name,
			Phase:     pod.Status.Phase,
			Status:    statue,
		},
	)
}

func (p *PodQueue) Worker() {
	for p.processNextWorkItem() {
	}
}

func (p *PodQueue) processNextWorkItem() bool {
	obj, quit := p.queue.Get()
	if quit {
		return false
	}
	defer p.queue.Done(obj)

	pod, ok := obj.(*vo.Pod)
	if ok {
		utilruntime.HandleError(db.SavePod(p.conf, pod))
	}
	return true
}
