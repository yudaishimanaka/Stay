package models

import (
	"net"
)

type User struct {
	UserId int `xorm:"user_id not null BIGINT pk autoincr" form:"user_id"`
	UserName string `xorm:"user_name not null" form:"user_name"`
	HwAddr string `xorm:"mac_address not null" form:"hw_addr"`
	IconPath string `xorm:"icon_path" form:"icon_path"`
}
func (u *User) GetParseHwAddr() ( net.HardwareAddr, error) {
	hwAddr, err := net.ParseMAC(u.HwAddr)
	if err != nil {
		return nil, err
	}

	return hwAddr, nil
}