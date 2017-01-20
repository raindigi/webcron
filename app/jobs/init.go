package jobs

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/lisijie/webcron/app/models"
	"os/exec"
	"time"
)

func InitJobs() {
	list, _ := models.TaskGetList(1, 1000000, "status", 1)
	for _, task := range list {
		job, err := NewJobFromTask(task)
		if err != nil {
			beego.Error("InitJobs:", err.Error())
			continue
		}
		AddJob(task.CronSpec, job)
	}
}

func runCmdWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	var err error
	select {
	case <-time.After(timeout):
		beego.Warn(fmt.Sprintf("Task execution time of more than %d seconds, the process will be forced to kill: %d", int(timeout/time.Second), cmd.Process.Pid))
		go func() {
			<-done // Read out the above goroutine data, to avoid blocking can not exit
		}()
		if err = cmd.Process.Kill(); err != nil {
			beego.Error(fmt.Sprintf("Process can not kill: %d, error message: %s", cmd.Process.Pid, err))
		}
		return err, true
	case err = <-done:
		return err, false
	}
}
