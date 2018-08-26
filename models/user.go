package models

import (
	"net"
)

type User struct {
	UserId int `xorm:"user_id not null BIGINT pk autoincr"`
	UserName string `xorm:"user_name not null"`
	HwAddr string `xorm:"mac_address not null"`
	IconPath string `xorm:"icon_path"`
}

func (u *User) GetUserId() int { return u.UserId }

func (u *User) GetUserName() string { return u.UserName }

func (u *User) GetHwAddr() ( net.HardwareAddr, error) {
	hwAddr, err := net.ParseMAC(u.HwAddr)
	if err != nil {
		return nil, err
	}

	return hwAddr, nil
}

func (u *User) GetIconPath() string { return u.IconPath }