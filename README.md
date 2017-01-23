This is a fork. I have already translated the main app using Google Translate. The next step is to run the app and correct the translation.
Feel free to give it a try and contribute.

# webcron
------------

A timed task manager, based on Go language and beego framework development. For unified management of the project timed tasks, provide a visual configuration interface, the implementation of logging, mail notification and other functions, without relying on * crontab Unix services.

## Background of the project

Development of this project is to solve my company's PHP projects in a number of regular tasks, the use of crontab bad management. Timing tasks in my project is also written in PHP, is part of the whole project, I hope to have a system that can configure these timed tasks, and can view the execution of each task, the task is completed or failed to automatically mail alert development Personnel, therefore did this project.

## Features

* Unified management of a variety of timing tasks.
* Second-level timer, the use of crontab time expression.
* Pause tasks at any time.
* Record the results of each task execution.
* Perform result mail notification.

## Interface screenshots

![webcron](https://raw.githubusercontent.com/lisijie/webcron/master/screenshot.png)


## Installation Notes

Go and MySQL are required.

Get the source code

	$ go get github.com/anselal/webcron
	
Open the configuration file conf/app.conf, modify the relevant configuration.
	

Create the database webcron and import install.sql

	$ mysql -u username -p -D webcron < install.sql

Build

	$ go build

Run
	
	$ ./webcron
	or
	$ nohup ./webcron 2>&1 > error.log &
	Set to run in the background

Access： 

http://localhost:8000

username：admin
password：admin888
