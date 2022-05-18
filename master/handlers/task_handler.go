package handlers

import (
	"github.com/AlaricGilbert/argos-core/master/dal"
	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	if tasks, err := dal.GetTaskList(); err == nil {
		retData(c, tasks)
	} else {
		retErr(c, err)
	}
}

func WriteTask(c *gin.Context) {
	prefix := c.Query("prefix")
	protocol := c.Query("protocol")

	// var task model.Task

	if prefix == "" {
		retErrMsg(c, "prefix cannot be empty")
		return
	}

	if _, err := dal.GetTask(prefix); err != nil {
		// there is no suck task
		// create new one if protocol is not empty.
		if protocol != "" {
			retUnwarpErr(c, dal.CreateTask(prefix, protocol))
		} else {
			retOK(c)
		}
	} else {
		if protocol == "" {
			// protocol equals empty, remove it.
			retUnwarpErr(c, dal.RemoveTask(prefix))
		} else {
			// updates task.
			retUnwarpErr(c, dal.UpdateTask(prefix, protocol))
		}
	}
}
