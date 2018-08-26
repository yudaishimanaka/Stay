package models

type User struct {
	UserId int `xorm:"user_id not null BIGINT pk autoincr"`
	UserName string `xorm:"user_name not null"`
	HwAddr string `xorm:"mac_address not null"`
	IconPath string `xorm:"icon_path"`
}
