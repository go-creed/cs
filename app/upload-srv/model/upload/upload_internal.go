package upload

import (
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	uploadPb "cs/app/upload-srv/proto/upload"
)

const (
	path = "/static/file/"
)

// get static's path
func (s *service) staticPath() string {
	pwd, _ := os.Getwd()
	pwd += path
	return pwd
}

// detailFileMate
func (s *service) detailFileMate(db *gorm.DB, info *uploadPb.FileMate) error {
	prepare, err := db.DB().Prepare("select * from `cs`.`file_mate` where `filesha256` = ? limit 1")
	if err != nil {
		return errors.WithStack(err)
	}
	defer prepare.Close()
	rows, err := prepare.Query(info.Filesha256)
	if err != nil {
		return errors.WithStack(err)
	}
	for rows.Next() {
		err = rows.Scan(
			&info.Id,
			&info.Filesha256,
			&info.Size,
			&info.Filename,
			&info.Location,
			&info.CreatedAt,
			&info.UpdatedAt,
			&info.DeletedAt)
		if err != nil {
			return errors.WithStack(err)
		}
		break
	}
	return nil
}

// isExistFileMate
func (s *service) isExistFileMate(db *gorm.DB, info *uploadPb.FileMate) (err error) {
	prepare, err := db.DB().Prepare("select id from `cs`.`file_mate` where `filesha256` = ? limit 1")
	if err != nil {
		return errors.WithStack(err)
	}
	defer prepare.Close()
	rows, err := prepare.Query(info.Filesha256)
	if err != nil {
		return errors.WithStack(err)
	}
	for rows.Next() {
		rows.Scan(&info.Id)
		return nil
	}
	return gorm.ErrRecordNotFound
}

// insertFileMate
func (s *service) insertFileMate(db *gorm.DB, info *uploadPb.FileMate) (err error) {
	//search file mate before insert
	if err = s.isExistFileMate(db, info); err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return errors.WithStack(err)
		}
	}
	prepare, err := db.DB().Prepare("insert into `cs`.`file_mate`(`filesha256`,`size`,`filename`,`location`,`created_at`,`updated_at`)VALUES(?,?,?,?,?,?)")
	if err != nil {
		return errors.WithStack(err)
	}
	defer prepare.Close()
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now
	result, err := prepare.Exec(info.Filesha256, info.Size, info.Filename, info.Location, info.CreatedAt, info.UpdatedAt)
	if err != nil {
		return errors.WithStack(err)
	}
	info.Id, err = result.LastInsertId()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
