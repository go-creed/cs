package user

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	userPb "cs/app/user-srv/proto/user"
)

type status int

const (
	Enable status = iota + 1
	Disable
	Locking
)

func (s *service) userInsert(db *gorm.DB, info *userPb.UserInfo) error {
	prepare, err := db.DB().Prepare("INSERT INTO `cs`.`user_info`(`user_name`,`password`,`phone`,`status`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?)")
	if err != nil {
		return errors.WithStack(err)
	}
	defer prepare.Close()
	now := time.Now().Unix()
	info.UpdatedAt = now
	info.CreatedAt = now
	info.Status = int64(Enable)
	exec, err := prepare.Exec(info.UserName, info.Password, info.Phone, info.Status, info.CreatedAt, info.UpdatedAt)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = exec.RowsAffected()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *service) userIsExist(db *gorm.DB, info *userPb.UserInfo) error {
	prepare, err := db.DB().Prepare("select `id` from `cs`.`user_info` where `phone` = ? limit 1")
	if err != nil {
		return errors.WithStack(err)
	}
	defer prepare.Close()

	// query data
	query, err := prepare.Query(info.Phone)
	if err != nil {
		return errors.WithStack(err)
	}
	// view length
	columns, err := query.Columns()
	if err != nil {
		return errors.WithStack(err)
	}
	if len(columns) != 0 {
		return errors.WithStack(errors.New("user info is already exists"))
	}
	return nil
}

func (s *service) userUpdate() {}
