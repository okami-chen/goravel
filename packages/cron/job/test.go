package job

import (
	"encoding/json"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
)

type CronTestHandle struct {
	//Name string
}

func (t CronTestHandle) Exec(arg any, content any) error {
	j, _ := json.Marshal(map[string]any{"date": carbon.Now().ToDateTimeString()})
	facades.Log().Info("参数: " + string(j))

	return nil
}
