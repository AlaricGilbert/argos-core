package model

type Task struct {
	ID       int64  `gorm:"column:id" db:"id" json:"-" form:"id"`
	Prefix   string `gorm:"column:prefix" db:"prefix" json:"prefix" form:"prefix"`
	Protocol string `gorm:"column:protocol" db:"protocol" json:"protocol" form:"protocol"`
}
