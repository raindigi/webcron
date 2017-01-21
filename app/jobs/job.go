package jobs

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/anselal/webcron/app/mail"
	"github.com/anselal/webcron/app/models"
	"html/template"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"
)

var mailTpl *template.Template

func init() {
	mailTpl, _ = template.New("mail_tpl").Parse(`
	你好 {{.username}}，<br/>

<p>The following are the results of the task execution：</p>

<p>
Task ID：{{.task_id}}<br/>
Task name：{{.task_name}}<br/>       
Start time：{{.start_time}}<br />
Process time：{{.process_time}}秒<br />
Status：{{.status}}
</p>
<p>-------------The following is the task execution output-------------</p>
<p>{{.output}}</p>
<p>
----------------------------------------------------------------------<br />
This message is automatically sent by the system, please do not reply<br />
If you want to cancel the mail notification, please log in to the system to set<br />
</p>
`)

}

type Job struct {
	id         int                                               // Task ID
	logId      int64                                             // Log ID
	name       string                                            // Task name
	task       *models.Task                                      // Task object
	runFunc    func(time.Duration) (string, string, error, bool) // Execute the function
	status     int                                               // Task status, greater than 0 for executing
	Concurrent bool                                              // Whether the same task is allowed to be executed in parallel
}

func NewJobFromTask(task *models.Task) (*Job, error) {
	if task.Id < 1 {
		return nil, fmt.Errorf("ToJob: Lackid")
	}
	job := NewCommandJob(task.Id, task.TaskName, task.Command)
	job.task = task
	job.Concurrent = task.Concurrent == 1
	return job, nil
}

func NewCommandJob(id int, name string, command string) *Job {
	job := &Job{
		id:   id,
		name: name,
	}
	job.runFunc = func(timeout time.Duration) (string, string, error, bool) {
		bufOut := new(bytes.Buffer)
		bufErr := new(bytes.Buffer)
		cmd := exec.Command("/bin/bash", "-c", command)
		cmd.Stdout = bufOut
		cmd.Stderr = bufErr
		cmd.Start()
		err, isTimeout := runCmdWithTimeout(cmd, timeout)

		return bufOut.String(), bufErr.String(), err, isTimeout
	}
	return job
}

func (j *Job) Status() int {
	return j.status
}

func (j *Job) GetName() string {
	return j.name
}

func (j *Job) GetId() int {
	return j.id
}

func (j *Job) GetLogId() int64 {
	return j.logId
}

func (j *Job) Run() {
	if !j.Concurrent && j.status > 0 {
		beego.Warn(fmt.Sprintf("Task [%d] The last execution has not finished and this is ignored.", j.id))
		return
	}

	defer func() {
		if err := recover(); err != nil {
			beego.Error(err, "\n", string(debug.Stack()))
		}
	}()

	if workPool != nil {
		workPool <- true
		defer func() {
			<-workPool
		}()
	}

	beego.Debug(fmt.Sprintf("Start the task: %d", j.id))

	j.status++
	defer func() {
		j.status--
	}()

	t := time.Now()
	timeout := time.Duration(time.Hour * 24)
	if j.task.Timeout > 0 {
		timeout = time.Second * time.Duration(j.task.Timeout)
	}

	cmdOut, cmdErr, err, isTimeout := j.runFunc(timeout)

	ut := time.Now().Sub(t) / time.Millisecond

	// Insert the log
	log := new(models.TaskLog)
	log.TaskId = j.id
	log.Output = cmdOut
	log.Error = cmdErr
	log.ProcessTime = int(ut)
	log.CreateTime = t.Unix()

	if isTimeout {
		log.Status = models.TASK_TIMEOUT
		log.Error = fmt.Sprintf("Task execution exceeded %d 秒\n----------------------\n%s\n", int(timeout/time.Second), cmdErr)
	} else if err != nil {
		log.Status = models.TASK_ERROR
		log.Error = err.Error() + ":" + cmdErr
	}
	j.logId, _ = models.TaskLogAdd(log)

	// Update the last execution time
	j.task.PrevTime = t.Unix()
	j.task.ExecuteTimes++
	j.task.Update("PrevTime", "ExecuteTimes")

	// Send an email notification
	if (j.task.Notify == 1 && err != nil) || j.task.Notify == 2 {
		user, uerr := models.UserGetById(j.task.UserId)
		if uerr != nil {
			return
		}

		var title string

		data := make(map[string]interface{})
		data["task_id"] = j.task.Id
		data["username"] = user.UserName
		data["task_name"] = j.task.TaskName
		data["start_time"] = beego.Date(t, "Y-m-d H:i:s")
		data["process_time"] = float64(ut) / 1000
		data["output"] = cmdOut

		if isTimeout {
			title = fmt.Sprintf("Task execution result notification #%d: %s", j.task.Id, "time out")
			data["status"] = fmt.Sprintf("Time out（%d seconds）", int(timeout/time.Second))
		} else if err != nil {
			title = fmt.Sprintf("Task execution result notification #%d: %s", j.task.Id, "failure")
			data["status"] = "failure（" + err.Error() + "）"
		} else {
			title = fmt.Sprintf("Task execution result notification #%d: %s", j.task.Id, "success")
			data["status"] = "success"
		}

		content := new(bytes.Buffer)
		mailTpl.Execute(content, data)
		ccList := make([]string, 0)
		if j.task.NotifyEmail != "" {
			ccList = strings.Split(j.task.NotifyEmail, "\n")
		}
		if !mail.SendMail(user.Email, user.UserName, title, content.String(), ccList) {
			beego.Error("Sending message timed out：", user.Email)
		}
	}
}
