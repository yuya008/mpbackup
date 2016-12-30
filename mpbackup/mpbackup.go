package mpbackup

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/mkideal/cli"
	"github.com/go-xorm/xorm"
	"fmt"
	"log"
	"time"
	"os"
	"github.com/yuya008/mpbackup/downloader"
	"github.com/go-xorm/core"
)

type Cfg struct {
	cli.Helper
	Dbhost string `cli:"dbhost" usage:"数据库host" dft:"0.0.0.0"`
	Dbuser string `cli:"dbuser" usage:"数据库用户名" dft:"root"`
	Dbpwd string `cli:"dbpwd" usage:"数据库密码" dft:""`
	Dbport int `cli:"dbport" usage:"数据库端口" dft:"3306"`
	Dbcharset string `cli:"dbcharset" usage:"数据库客户端字符集" dft:"utf8mb4"`
	Dbname string `cli:"dbname" usage:"数据库名" dft:"mepainew"`
	LogPath string `cli:"logpath" usage:"日志数据文件存储路径" dft:"."`
	BackupPath string `cli:"backuppath" usage:"备份输出目录" dft:"."`
	DownLoaderThreadN int `cli:"downloaderthreadn" usage:"下载器线程数" dft:"1000"`
	DownLoaderTaskQueueLen int `cli:"downLoadertaskqueuelen" usage:"下载器任务队列长度" dft:"10000"`
}

type Mpbackup struct {
	// 日志数据库
	//logdb *leveldb.DB
	// 数据库
	dbEngine *xorm.Engine
	// 下载器
	downloader downloader.DownLoader
	// 备份输出目录
	backupPath string
}

const (
	Varsion = "0.0.1"
	// 10秒一ping
	PingRate = 10
	// 日志文件名
	RunTimeLogFile = "/mpbackup.log"
	// 数据库路径
	DataPath = "/data"
	// 文件和文件夹权限
	FileMode = 0755
	// 下载日志
	DownLoaderLogFile = "/downloader.log"
	// 资源URL前缀
	URLPrefix = "https://images.mepai.me"
)

var (
	runlog *log.Logger
)

func (mp *Mpbackup) pinger() {
	for {
		runlog.Println("PING -> DB")
		err := mp.dbEngine.Ping()
		if err != nil {
			runlog.Panicln(err)
		}
		runlog.Println("PONG <- DB")
		time.Sleep(time.Second * PingRate)
	}
}

func (mp *Mpbackup) init(c *Cfg) {
	var err error
	fileInfo, err := os.Stat(c.LogPath)
	if err != nil {
		log.Panicln(err)
	}
	if !fileInfo.IsDir() {
		log.Panicln(c.LogPath + " not dir")
	}
	file, err := os.OpenFile(c.LogPath + RunTimeLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, FileMode)
	if err != nil {
		log.Panicln(err.Error())
	}
	runlog = log.New(file, "", log.LstdFlags)
	// 初始化数据库
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s", c.Dbuser, c.Dbpwd, c.Dbhost, c.Dbport, c.Dbname, c.Dbcharset,
	)
	runlog.Println("NewEngine " + dataSourceName)
	mp.dbEngine, err = xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		runlog.Panicln(err)
	}
	mp.dbEngine.SetColumnMapper(core.GonicMapper{})
	// 定时ping
	runlog.Println("开启数据库ping任务")
	go mp.pinger()
	runlog.Println("初始化下载器")
	mp.downloader = downloader.NewHttpDownLoader(c.DownLoaderThreadN, c.DownLoaderTaskQueueLen)
	file, err = os.OpenFile(c.LogPath + DownLoaderLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, FileMode)
	if err != nil {
		runlog.Panicln(err.Error())
	}
	mp.downloader.SetLoggerWriter(file, downloader.LogFlag)
	mp.backupPath = c.BackupPath
}

func (mp *Mpbackup) Run(c *Cfg) error {
	mp.init(c)
	go func() {
		runlog.Println("开始推送用户任务")
		mp.pushUserToDownLoader()
		runlog.Println("开始推送作品任务")
		mp.pushWorksAppendixToDownLoader()
		runlog.Println("开始推送活动任务")
		mp.pushActivityToDownLoader()
		mp.downloader.Done()
	}()
	mp.downloader.DownLoad()
	return nil
}
