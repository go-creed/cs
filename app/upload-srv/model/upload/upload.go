package upload

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"

	uploadSrv "cs/app/upload-srv/proto/upload"
	"cs/public/util"

	//"github.com/micro/go-micro/v2/config/source/file"
)

type service struct {
}

func (s *service) FileDetail(db *gorm.DB, data *uploadSrv.FileMate, condition ...string) error {
	err := s.detailFileMate(db, data)
	if err != nil {
		return fmt.Errorf("[Upload][FileDetail] 获取文件详情失败, err:%s", err.Error())
	}
	return nil
}

func (s *service) WriteDB(db *gorm.DB, data *uploadSrv.FileMate) error {
	err := s.insertFileMate(db, data)
	if err != nil {
		return fmt.Errorf("[Upload][WriteDB] 数据库写入失败, err:%s", err)
	}
	return nil
}

func (s *service) Hash(file *os.File) (hashName string, err error) {
	//file.Seek(0, 0) //重置文件游标
	//all, err := ioutil.ReadAll(file)
	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", fmt.Errorf("[Upload][Hash] 数据拷贝失败，err:%s", err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (s *service) MergeFile(fileName string, filesha256 string) (src string, file string, err error) {
	var (
		pwd        = s.staticPath()
		path       = pwd + fileName
		verifyFile bool
	)
	src, file = filepath.Split(path)
	dest := src + "../" + file
	if _, err = util.MergeFile(src, dest); err != nil {
		return
	}
	if verifyFile, err = util.VerifyFile(filesha256, dest); err != nil || !verifyFile {
		return
	}
	_ = os.RemoveAll(src)
	src, file = filepath.Split(dest)
	return src, file, nil
}

func (s *service) CreateFile(path string) (*os.File, string, error) {
	dir := filepath.Dir(path)
	pwd := s.staticPath()
	if err := os.MkdirAll(pwd+dir, os.ModePerm); err != nil {
		return nil, "", fmt.Errorf("[Upload][CreateFile] 创建文件目录, err:%s", err)
	}
	location := pwd + path
	file, err := os.Create(location)
	if err != nil {
		return nil, "", fmt.Errorf("[Upload][SendBytes] 打开文件失败, err:%s", err)
	}
	return file, location, nil
}

func (s *service) Write(file *os.File, bytes []byte) (err error) {
	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("[Upload][SendBytes] 写入文件失败, err:%s", err)
	}
	return nil
}

func GetService() (Service, error) {
	if s == nil {
		return nil, fmt.Errorf("[GetService] Upload's Model Init Failure")
	}
	return s, nil
}

type FileBytes struct {
	Bytes    *uploadSrv.Bytes `json:"-"`
	File     *os.File
	FilePath string `json:"file_path"`
}

type Service interface {
	Write(file *os.File, bytes []byte) error                                     //写图片
	CreateFile(path string) (*os.File, string, error)                            //创建文件
	Hash(file *os.File) (string, error)                                          //Hash
	WriteDB(db *gorm.DB, data *uploadSrv.FileMate) error                         //写入db文件
	FileDetail(db *gorm.DB, data *uploadSrv.FileMate, condition ...string) error //获取文件详情
	MergeFile(fileName string, filesha256 string) (string, string, error)
}
