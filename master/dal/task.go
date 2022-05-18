package dal

import "github.com/AlaricGilbert/argos-core/master/model"

func GetTaskList() (tasks []model.Task, err error) {
	return tasks, db.Table("tasks").Find(&tasks).Error
}

func GetTask(prefix string) (*model.Task, error) {
	var task model.Task
	if err := db.Table("tasks").Where("prefix = ?", prefix).First(&task).Error; err == nil {
		return &task, nil
	} else {
		return nil, err
	}
}

func CreateTask(prefix, protocol string) error {
	return db.Table("tasks").Create(&model.Task{Prefix: prefix, Protocol: protocol}).Error
}

func UpdateTask(prefix, protocol string) error {
	return db.Table("tasks").Where("prefix = ?", prefix).Update("protocol", protocol).Error
}

func RemoveTask(prefix string) error {
	return db.Table("tasks").Where("prefix = ?", prefix).Delete(&model.Task{}).Error
}
