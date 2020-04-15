package upload

import (
	"database/sql"
	"path/filepath"
	"time"

	uploadPb "cs/app/upload-srv/proto/upload"
	"github.com/pkg/errors"
)

func (s *service) imagesExt(fileName string) (string, bool) {
	ext := filepath.Ext(fileName)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".ico", ".svg":
		return ext, true
	default:
		return "", false
	}
}

// searchFileMate
func (s *service) isExistFileMate(db *sql.DB, info *uploadPb.FileInfo) (isExist bool, err error) {
	prepare, err := db.Prepare("select id from `cs`.`file_mate` where `filesha256` = ? limit 1")
	if err != nil {
		return false, errors.WithStack(err)
	}
	defer prepare.Close()
	rows, err := prepare.Query(info.Filesha256)
	if err != nil {
		return false, errors.WithStack(err)
	}
	for rows.Next() {
		rows.Scan(&info.Id)
		return true, nil
	}
	return false, nil
}

// insertFileMate
func (s *service) insertFileMate(db *sql.DB, info *uploadPb.FileInfo) (err error) {
	//search file mate before insert
	var isExist bool
	if isExist, err = s.isExistFileMate(db, info); err != nil {
		return errors.WithStack(err)
	}
	if isExist {
		return nil
	}
	prepare, err := db.Prepare("insert into `cs`.`file_mate`(`filesha256`,`size`,`filename`,`location`,`created_at`,`updated_at`)VALUES(?,?,?,?,?,?)")
	if err != nil {
		return errors.WithStack(err)
	}
	defer prepare.Close()
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now
	result, err := prepare.Exec(info.Filesha256, info.Size, info.FileName, info.Location, info.CreatedAt, info.UpdatedAt)
	if err != nil {
		return errors.WithStack(err)
	}
	info.Id, err = result.LastInsertId()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
