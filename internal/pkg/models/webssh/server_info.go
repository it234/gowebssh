package webssh

import (
	"time"

	"github.com/it234/gowebssh/internal/pkg/models/basemodel"

	"github.com/jinzhu/gorm"
)

// 服务器信息
type ServerInfo struct {
	basemodel.Model
	Memo          string `gorm:"column:memo;size:128;" json:"memo" form:"memo"`                                        //备注
	ServerAddress string `gorm:"column:server_address;size:128;not null;" json:"server_address" form:"server_address"` // 服务器地址
	UserName      string `gorm:"column:user_name;size:128;not null;" json:"user_name" form:"user_name"`                // 登录用户名
	Password      string `gorm:"column:password;type:char(128);not null;" json:"password" form:"password"`             // 登录密码
	AliasName     string `gorm:"column:alias_name;size:64;" json:"alias_name" form:"alias_name"`                       // 服务器别名
}

// 表名
func (ServerInfo) TableName() string {
	return TableName("server_info")
}

// 添加前
func (m *ServerInfo) BeforeCreate(scope *gorm.Scope) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// 更新前
func (m *ServerInfo) BeforeUpdate(scope *gorm.Scope) error {
	m.UpdatedAt = time.Now()
	return nil
}
