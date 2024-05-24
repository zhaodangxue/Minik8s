package hpacontroller

import "minik8s/apiobjects"

type Worker interface {
	Run()
	Done()
	SyncCh() chan<- struct{}
}
type worker struct {
	syncCh chan struct{}
	target *apiobjects.HorizontalPodAutoscaler
}
