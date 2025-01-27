package job

import (
	"github.com/goravel/framework/facades"
	"github.com/robfig/cron/v3"
	"time"
)

var jobList map[string]JobsExec

type JobCore struct {
	EntryId        int
	MisfirePolicy  int
	Name           string
	InvokeTarget   string
	Args           string
	CronExpression string
	Content        string
}

func InitJob() {
	config := facades.Config().Get("cron")

	if config == nil {
		facades.Log().Infof("无任务待添加")
		return
	}
	cfg := config.(map[string]interface{})
	jobList = make(map[string]JobsExec)
	for name, exec := range cfg {
		jobList[name] = exec.(JobsExec)
		facades.Log().Infof("任务 %s -> %#v", name, exec.(JobsExec))
	}
}

type ExecJob struct {
	cron *cron.Cron
	JobCore
}

func (e *ExecJob) Run() {
	startTime := time.Now()
	var obj = jobList[e.InvokeTarget]
	err := CallExec(obj.(JobsExec), e.Args, e.Content)
	if err != nil {
		facades.Log().Errorf("任务错误: %s", err.Error())
		Remove(e.cron, e.EntryId)
	} else {
		latencyTime := time.Now().Sub(startTime)
		facades.Log().Infof("任务耗时: %s -> %fms", e.InvokeTarget, latencyTime.Seconds()*1000)
	}
	if e.MisfirePolicy == 1 {
		Remove(e.cron, e.EntryId)
	}
	return
}

func Setup() {
	ret, err := facades.App().Make("cron.core")
	if err != nil {
		facades.Log().Error(err.Error())
		return
	}
	c := ret.(*cron.Cron)
	if err != nil {
		facades.Log().Error("JobCore Remove entry_id error", err)
	}
	for target, _ := range jobList {
		job := JobCore{
			Name:           target,
			InvokeTarget:   target,
			Args:           `{"name":"hello"}`,
			CronExpression: "*/5 * * * * *",
			MisfirePolicy:  0,
			Content:        "",
		}
		exec := ExecJob{
			cron:    c,
			JobCore: job,
		}
		id, e := exec.addJob(c)
		if e != nil {
			facades.Log().Errorf(e.Error())
			continue
		}
		facades.Log().Infof("id: %d, job: %s", id, target)
	}
	// 其中任务
	c.Start()
	facades.Log().Info("任务启动")
	// 关闭任务
}

func (h *ExecJob) addJob(c *cron.Cron) (int, error) {
	id, err := c.AddJob(h.CronExpression, h)
	if err != nil {
		return 0, err
	}
	h.cron = c
	h.EntryId = int(id)
	return h.EntryId, nil
}

// 移除任务(停止任务)
func Remove(c *cron.Cron, entryID int) {
	c.Remove(cron.EntryID(entryID))
}
