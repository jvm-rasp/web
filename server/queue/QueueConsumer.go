package queue

import (
	"context"
	"server/controller"
)

func ConsumerLog() {
	logHandle := controller.NewLogController()
	for {
		select {
		case log := <-controller.LogChan:
			logHandle.HandleLog(log)
		case <-context.TODO().Done():
			return
		}
	}
}
