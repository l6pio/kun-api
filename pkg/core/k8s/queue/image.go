package queue

import (
	"github.com/labstack/gommon/log"
	v1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/util/workqueue"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db"
	"l6p.io/kun/api/pkg/core/db/vo"
	"l6p.io/kun/api/pkg/core/service"
)

type ImageQueue struct {
	conf  *core.Config
	queue workqueue.RateLimitingInterface
}

func NewImageQueue(conf *core.Config) *ImageQueue {
	return &ImageQueue{
		conf: conf,
		queue: workqueue.NewNamedRateLimitingQueue(
			workqueue.DefaultControllerRateLimiter(),
			"ImageQueue",
		),
	}
}

func (p *ImageQueue) ShutDown() {
	p.ShutDown()
}

func (p *ImageQueue) Push(pod *v1.Pod, status vo.ImageStatus) {
	images := make(map[string]*vo.ImageTimeline)
	for _, container := range pod.Status.ContainerStatuses {
		if container.State.Running != nil && container.State.Terminated == nil && container.State.Waiting == nil {
			images[container.ImageID] = &vo.ImageTimeline{
				Timestamp: container.State.Running.StartedAt.UnixNano() / 1e6,
				ImageId:   container.ImageID,
				Image:     container.Image,
				Status:    status,
			}
		}
	}

	for _, timeline := range images {
		p.queue.Add(timeline)
	}
}

func (p *ImageQueue) Worker() {
	for p.processNextWorkItem() {
	}
}

func (p *ImageQueue) processNextWorkItem() bool {
	obj, quit := p.queue.Get()
	if quit {
		return false
	}
	defer p.queue.Done(obj)

	imageTimeline, ok := obj.(*vo.ImageTimeline)
	if ok {
		if imageTimeline.Status == vo.ImageUp {
			utilruntime.HandleError(processImageUp(p.conf, imageTimeline))
		} else {
			utilruntime.HandleError(processImageDown(p.conf, imageTimeline))
		}
	}
	return true
}

func processImageUp(conf *core.Config, event *vo.ImageTimeline) error {
	exists, err := db.ImageExists(conf, event.ImageId)
	if err != nil {
		return err
	}

	if !exists && false {
		report, err := service.ScanImageReport(event.Image)
		if err != nil {
			return err
		}

		if len(report.Matches) == 0 {
			log.Info("no vulnerabilities found")
		}

		log.Info("start saving scan results")
		service.InsertImageReport(conf, event.ImageId, report)
		log.Infof("scan results of '%v' has been saved", event.ImageId)
	}

	if err := db.SaveImageTimelineEvent(conf, event); err != nil {
		return err
	}
	return db.UpdateImagePods(conf, event)
}

func processImageDown(conf *core.Config, event *vo.ImageTimeline) error {
	if err := db.SaveImageTimelineEvent(conf, event); err != nil {
		return err
	}
	return db.UpdateImagePods(conf, event)
}
