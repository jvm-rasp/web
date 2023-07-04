package job

import (
	cron "github.com/robfig/cron/v3"
	"server/config"
	"server/repository"
)

type TableDeleteJob struct {
	Enable         bool   // 是否开启
	CronExpression string // 执行频率
	TableMaxSize   int    // 最大数据 10w
}

func NewTableDeleteJob() *TableDeleteJob {
	return &TableDeleteJob{
		Enable:         config.Conf.TableDeleteJob.Enable,
		TableMaxSize:   config.Conf.TableDeleteJob.TableMaxSize,
		CronExpression: config.Conf.TableDeleteJob.CronExpression,
	}
}

func (j *TableDeleteJob) Run() {
	if !j.Enable {
		return
	}
	c := cron.New()
	// 定时执行
	c.AddFunc(j.CronExpression, func() {
		logRepository := repository.NewOperationLogRepository()
		attackRepository := repository.NewRaspAttackRepository()
		raspErrorLogsRepository := repository.NewRaspErrorLogsRepository()
		logRepository.DeleteOperationLogsByJob(j.TableMaxSize)
		attackRepository.DeleteAttackLogsByJob(j.TableMaxSize)
		raspErrorLogsRepository.DeleteRaspLogsByJob(j.TableMaxSize)
	})

	c.Start()
}
