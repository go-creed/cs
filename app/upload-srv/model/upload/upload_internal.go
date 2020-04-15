package upload

import (
	"database/sql"
	"path/filepath"
	"time"

	uploadPb "cs/app/upload-srv/proto/upload"
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

//INSERT INTO `cs`.`user`(`id`, `userName`, `password`) VALUES (1, '123', '123');
func (s *service) insertFileMate(db *sql.DB, info *uploadPb.FileInfo) (err error) {
	prepare, err := db.Prepare("insert into `cs`.`file_mate`(`filesha256`,`filename`,`location`,`created_at`,`updated_at`)VALUES(?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer prepare.Close()
	//now := ptypes.TimestampNow()
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now
	result, err := prepare.Exec(info.Filesha256, info.FileName, info.Location, info.CreatedAt, info.UpdatedAt)
	if err != nil {
		return err
	}
	info.Id, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
