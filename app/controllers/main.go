package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"github.com/anselal/webcron/app/jobs"
	"github.com/anselal/webcron/app/libs"
	"github.com/anselal/webcron/app/models"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type MainController struct {
	BaseController
}

// Home
func (this *MainController) Index() {
	this.Data["pageTitle"] = "System overview"

	// The task to be performed
	entries := jobs.GetEntries(30)
	jobList := make([]map[string]interface{}, len(entries))
	for k, v := range entries {
		row := make(map[string]interface{})
		job := v.Job.(*jobs.Job)
		row["task_id"] = job.GetId()
		row["task_name"] = job.GetName()
		row["next_time"] = beego.Date(v.Next, "Y-m-d H:i:s")
		jobList[k] = row
	}

	// The most recently executed log
	logs, _ := models.TaskLogGetList(1, 20)
	recentLogs := make([]map[string]interface{}, len(logs))
	for k, v := range logs {
		task, err := models.TaskGetById(v.TaskId)
		taskName := ""
		if err == nil {
			taskName = task.TaskName
		}
		row := make(map[string]interface{})
		row["task_name"] = taskName
		row["id"] = v.Id
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["process_time"] = float64(v.ProcessTime) / 1000
		row["ouput_size"] = libs.SizeFormat(float64(len(v.Output)))
		row["output"] = beego.Substr(v.Output, 0, 100)
		row["status"] = v.Status
		recentLogs[k] = row
	}

	// The most recent failed log
	logs, _ = models.TaskLogGetList(1, 20, "status__lt", 0)
	errLogs := make([]map[string]interface{}, len(logs))
	for k, v := range logs {
		task, err := models.TaskGetById(v.TaskId)
		taskName := ""
		if err == nil {
			taskName = task.TaskName
		}
		row := make(map[string]interface{})
		row["task_name"] = taskName
		row["id"] = v.Id
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["process_time"] = float64(v.ProcessTime) / 1000
		row["ouput_size"] = libs.SizeFormat(float64(len(v.Output)))
		row["error"] = beego.Substr(v.Error, 0, 100)
		row["status"] = v.Status
		errLogs[k] = row
	}

	this.Data["recentLogs"] = recentLogs
	this.Data["errLogs"] = errLogs
	this.Data["jobs"] = jobList
	this.Data["cpuNum"] = runtime.NumCPU()
	this.display()
}

// Personal information
func (this *MainController) Profile() {
	beego.ReadFromRequest(&this.Controller)
	user, _ := models.UserGetById(this.userId)

	if this.isPost() {
		flash := beego.NewFlash()
		user.Email = this.GetString("email")
		user.Update()
		password1 := this.GetString("password1")
		password2 := this.GetString("password2")
		if password1 != "" {
			if len(password1) < 6 {
				flash.Error("Password length must be greater than 6 characters")
				flash.Store(&this.Controller)
				this.redirect(beego.URLFor(".Profile"))
			} else if password2 != password1 {
				flash.Error("The two passwords did not match")
				flash.Store(&this.Controller)
				this.redirect(beego.URLFor(".Profile"))
			} else {
				user.Salt = string(utils.RandomCreateBytes(10))
				user.Password = libs.Md5([]byte(password1 + user.Salt))
				user.Update()
			}
		}
		flash.Success("Successfully modified！")
		flash.Store(&this.Controller)
		this.redirect(beego.URLFor(".Profile"))
	}

	this.Data["pageTitle"] = "Personal information"
	this.Data["user"] = user
	this.display()
}

// Log in
func (this *MainController) Login() {
	if this.userId > 0 {
		this.redirect("/")
	}
	beego.ReadFromRequest(&this.Controller)
	if this.isPost() {
		flash := beego.NewFlash()

		username := strings.TrimSpace(this.GetString("username"))
		password := strings.TrimSpace(this.GetString("password"))
		remember := this.GetString("remember")
		if username != "" && password != "" {
			user, err := models.UserGetByName(username)
			errorMsg := ""
			if err != nil || user.Password != libs.Md5([]byte(password+user.Salt)) {
				errorMsg = "Username or password did not match"
			} else if user.Status == -1 {
				errorMsg = "The account is disabled"
			} else {
				user.LastIp = this.getClientIp()
				user.LastLogin = time.Now().Unix()
				models.UserUpdate(user)

				authkey := libs.Md5([]byte(this.getClientIp() + "|" + user.Password + user.Salt))
				if remember == "yes" {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey, 7*86400)
				} else {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey)
				}

				this.redirect(beego.URLFor("TaskController.List"))
			}
			flash.Error(errorMsg)
			flash.Store(&this.Controller)
			this.redirect(beego.URLFor("MainController.Login"))
		}
	}

	this.TplName = "main/login.html"
}

// Sign out
func (this *MainController) Logout() {
	this.Ctx.SetCookie("auth", "")
	this.redirect(beego.URLFor("MainController.Login"))
}

// Gets the system time
func (this *MainController) GetTime() {
	out := make(map[string]interface{})
	out["time"] = time.Now().UnixNano() / int64(time.Millisecond)
	this.jsonResult(out)
}
