package model

type Record struct {
	ID        int64  `gorm:"column:id" db:"id" json:"-" form:"id"`
	Txid      string `gorm:"column:txid" db:"txid" json:"txid" form:"txid"`
	Timestamp int64  `gorm:"column:timestamp" db:"timestamp" json:"timestamp" form:"timestamp"`
	SourceIp  string `gorm:"column:source_ip" db:"source_ip" json:"source_ip" form:"source_ip"`
	Sniffer   string `gorm:"column:sniffer" db:"sniffer" json:"sniffer" form:"sniffer"`
	Protocol  string `gorm:"column:protocol" db:"protocol" json:"protocol" form:"protocol"`
	Method    string `gorm:"column:method" db:"method" json:"method" form:"method"`
}
