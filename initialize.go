package Timer

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

var T = NewTimer()

func initialize() {
	go func() {
		var option []cron.Option
		option = append(option, cron.WithSeconds())

		_, err := T.AddTaskByFunc(cronName, "@daily", func() {
			// task running
			// task  is defined in task.go
			fmt.Println("task running")
		}, taskName, option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		// other tasks are defined here
		// ......
	}()
}
