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
	"time"
)

type ImageTimelineQueue struct {
	conf  *core.Config
	queue workqueue.RateLimitingInterface
}

func NewImageTimelineQueue(conf *core.Config) *ImageTimelineQueue {
	return &ImageTimelineQueue{
		conf: conf,
		queue: workqueue.NewNamedRateLimitingQueue(
			workqueue.DefaultControllerRateLimiter(),
			"ImageTimelineQueue",
		),
	}
}

func (p *ImageTimelineQueue) ShutDown() {
	p.ShutDown()
}

func (p *ImageTimelineQueue) Push(pod *v1.Pod, status vo.ImageStatus) {
	timestamp := time.Now().UnixNano() / 1e6

	images := make(map[string]string)
	for _, container := range pod.Status.ContainerStatuses {
		images[container.ImageID] = container.Image
	}

	for imageId, image := range images {
		p.queue.Add(
			&vo.ImageTimeline{
				Timestamp: timestamp,
				ImageId:   imageId,
				Image:     image,
				Status:    status,
			},
		)
	}
}

func (p *ImageTimelineQueue) Worker() {
	for p.processNextWorkItem() {
	}
}

func (p *ImageTimelineQueue) processNextWorkItem() bool {
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

	if !exists {
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
