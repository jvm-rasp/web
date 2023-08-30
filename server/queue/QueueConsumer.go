package queue

import (
	"context"
	"server/controller"
	"server/socket"
)

func ConsumerLog() {
	for {
		select {
		// todo 加上定时器
		case log := <-socket.LogChan:
			ctl := controller.NewLogController()
			ctl.ReportLogFromSocket(log)
		case <-context.TODO().Done():
			return
		}
	}
}
