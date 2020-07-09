package user

import (
	"time"

	log "github.com/micro/go-micro/v2/logger"
	"micro_demo/basic/db"
	proto "micro_demo/proto/user"
)

// User ...
type User struct {
	ID        int64      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserID    int64      `json:"user_id" gorm:"column:user_id"`
	UserName  string     `json:"user_name"  gorm:"column:user_id"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

// QueryUserByName 查询名字
func (s *service) QueryUserByName(userName string) (ret *proto.User, err error) {
	queryString := `SELECT user_id, user_name, pwd FROM user WHERE user_name = ?`

	// 获取数据库
	o := db.GetDB()

	ret = &proto.User{}

	// 查询
	err = o.QueryRow(queryString, userName).Scan(&ret.Id, &ret.Name, &ret.Pwd)
	if err != nil {
		log.Infof("[QueryUserByName] 查询数据失败，err：%s", err)
		return
	}
	return
}
