package queue

import (
	v1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/util/workqueue"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db"
	"l6p.io/kun/api/pkg/core/db/vo"
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
	var podStatue string
	var podReady, podFinished int64
	var podRestartCount int32

	if pod.Status.Phase == v1.PodFailed {
		for _, container := range pod.Status.ContainerStatuses {
			startAt := container.State.Terminated.StartedAt.UnixNano() / 1e6
			if startAt > podReady {
				podReady = startAt
			}

			if container.State.Terminated.ExitCode > 0 {
				podStatue = container.State.Terminated.Reason
				podFinished = container.State.Terminated.FinishedAt.UnixNano() / 1e6
			}
		}
	} else if pod.Status.Phase == v1.PodSucceeded {
		for _, container := range pod.Status.ContainerStatuses {
			startAt := container.State.Terminated.StartedAt.UnixNano() / 1e6
			if startAt > podReady {
				podReady = startAt
			}

			finished := container.State.Terminated.FinishedAt.UnixNano() / 1e6
			if finished > podFinished {
				podStatue = container.State.Terminated.Reason
				podFinished = finished
			}
		}
	} else if pod.Status.Phase == v1.PodRunning {
		podStatue = string(v1.PodRunning)
		for _, container := range pod.Status.ContainerStatuses {
			if container.State.Running != nil {
				startAt := container.State.Running.StartedAt.UnixNano() / 1e6
				if startAt > podReady {
					podReady = startAt
				}
			}

			if container.State.Terminated != nil {
				startAt := container.State.Terminated.StartedAt.UnixNano() / 1e6
				if startAt > podReady {
					podReady = startAt
				}
			}
		}
	} else if pod.Status.Phase == v1.PodPending {
		for _, container := range pod.Status.ContainerStatuses {
			if container.State.Waiting != nil {
				podStatue = container.State.Waiting.Reason
			}
		}
	}

	for _, container := range pod.Status.ContainerStatuses {
		podRestartCount += container.RestartCount
	}

	p.queue.Add(
		&vo.Pod{
			Name:         pod.Name,
			Namespace:    pod.Namespace,
			Phase:        pod.Status.Phase,
			Status:       podStatue,
			Ready:        podReady,
			Finished:     podFinished,
			RestartCount: podRestartCount,
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
