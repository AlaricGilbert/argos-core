package dal

import "github.com/AlaricGilbert/argos-core/master/model"

func CreateRecord(r *model.Record) error {
	return db.Table("records").Create(r).Error
}

func GetRecordsWithIP(ip string) ([]model.Record, error) {
	var records []model.Record

	return records, db.Table("records").Where("source_ip = ?", ip).Find(&records).Error
}

func GetRecordsWithTxid(txid string) ([]model.Record, error) {
	var records []model.Record

	return records, db.Table("records").Where("txid = ?", txid).Find(&records).Error
}
