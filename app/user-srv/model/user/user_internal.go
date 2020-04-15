package user

import (
	"fmt"
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

var (
	UserInfoIsExist = errors.New("user info is already exists")
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

func (s *service) userLogin(db *gorm.DB, info *userPb.UserInfo) error {
	prepare, err := db.DB().Prepare("select `id` from `cs`.`user_info` where `user_name`=? and `password`=? and `status`=?")
	if err != nil {
		return errors.WithStack(err)
	}
	defer prepare.Close()

	row := prepare.QueryRow(info.UserName, info.Password, Enable)
	fmt.Println(row)
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
	var id int64
	for query.Next() {
		err = query.Scan(&id)
		if err != nil {
			return errors.WithStack(err)
		}
		break
	}
	if id != 0 {
		return UserInfoIsExist
	}
	return nil
}

func (s *service) userUpdate() {}
