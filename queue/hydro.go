package queue

import (
	"hydra"
	"time"
)

func Trigger(name string, data []byte, delay time.Duration) uint64 {
	hydra.ConfPath = "/home/luopan/devspace/go-hydra/src/config.json"
	bstalk := hydra.GetBStalkIns()
	event := hydra.NewEvent(name, data, "", 0)
	jobId := bstalk.Trigger(*event, 0, 0)
	return jobId
}
