package upload

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	uploadSrv "cs/app/upload-srv/proto/upload"
	log "github.com/micro/go-micro/v2/logger"

	//"github.com/micro/go-micro/v2/config/source/file"
)

const (
	path = "/app/upload-srv/static/file/"
)

var (
	once sync.Once
	s    *service
)

type service struct {
}

func (s *service) WriteDB(db *sql.DB, data *uploadSrv.FileInfo) error {
	err := s.insertFileMate(db, data)
	if err != nil {
		return fmt.Errorf("[Upload][WriteDB] 数据库写入失败, err:%s", err)
	}
	return nil
}

func (s *service) Hash(file *os.File) (hashName string, err error) {
	file.Seek(0, 0) //重置文件游标
	all, err := ioutil.ReadAll(file)
	//_, err = io.Copy(hash, file)
	if err != nil {
		return "", fmt.Errorf("[Upload][Hash] 数据拷贝失败，err:%s", err)
	}
	hash := sha256.New()
	return hex.EncodeToString(hash.Sum(all)), nil
}

func (s *service) CreateFile(fileName string) (*os.File, error) {
	pwd, _ := os.Getwd()
	pwd += path
	if err := os.MkdirAll(pwd, os.ModePerm); err != nil {
		return nil, fmt.Errorf("[Upload][CreateFile] 创建文件目录, err:%s", err)
	}
	pwd += fileName
	file, err := os.Create(pwd)
	if err != nil {
		return nil, fmt.Errorf("[Upload][SendBytes] 打开文件失败, err:%s", err)
	}
	return file, nil
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
	Write(file *os.File, bytes []byte) error            //写图片
	CreateFile(fileName string) (*os.File, error)       //创建文件
	Hash(file *os.File) (string, error)                 //Hash
	WriteDB(db *sql.DB, data *uploadSrv.FileInfo) error //写入db文件
}

// Init Service Model Like Redis, Mysql ....
func Init() {
	once.Do(func() {
		log.Info("[Upload][Model] init service model like redis,mysql...")
		s = &service{}
	})
}
