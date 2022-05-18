package dal

import (
	"errors"

	"github.com/AlaricGilbert/argos-core/master/model"
	"gorm.io/gorm"
)

type ConclusionQuery struct {
	Offset   int64 // when offset set, GetConclusions will return ids > Offset
	Limits   int   // largest query result number
	TimeFrom int64
	TimeTo   int64
	Method   string
}

func CreateOrUpdateConclustion(r *model.Record) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var record model.Record
		if err := tx.Table("conclusions").Where("txid = ? AND method = ?", r.Txid, r.Method).First(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// create
				return tx.Table("conclusions").Create(r).Error
			}
			return err
		} else {
			if r.Timestamp < record.Timestamp {
				return tx.Table("conclusions").Model(&record).Updates(r).Error
			}
			return nil
		}
	})
}

func GetConclusions(q ConclusionQuery) ([]model.Record, error) {
	var records []model.Record
	query := db.Debug().Table("conclusions").Where("method = ?", q.Method)

	if q.Offset != 0 {
		query = query.Where("id > ?", q.Offset)
	}

	if q.TimeFrom != 0 && q.TimeTo != 0 {
		query = query.Where("timestamp between ? and ?", q.TimeFrom, q.TimeTo)
	}

	if q.Limits < 20 || q.Limits > 100 {
		q.Limits = 20
	}

	return records, query.Group("source_ip").Order("timestamp desc").Limit(q.Limits).Find(&records).Error
}
