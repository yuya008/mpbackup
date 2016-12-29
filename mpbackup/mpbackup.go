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
	//"github.com/syndtr/goleveldb/leveldb"
	
)

type Cfg struct {
	cli.Helper
	Dbhost		string	`cli:"dbhost" usage:"数据库host" dft:"0.0.0.0"`
	Dbuser		string	`cli:"dbuser" usage:"数据库用户名" dft:"root"`
	Dbpwd		string	`cli:"dbpwd" usage:"数据库密码" dft:""`
	Dbport 		int		`cli:"dbport" usage:"数据库端口" dft:"3306"`
	Dbcharset 	string	`cli:"dbcharset" usage:"数据库客户端字符集" dft:"utf8mb4"`
	Dbname		string  `cli:"dbname" usage:"数据库名" dft:"mepainew"`
	LogPath		string	`cli:"logpath" usage:"日志数据文件存储路径" dft:"."`
	BackupPath	string 	`cli:"backuppath" usage:"备份输出目录" dft:"."`
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
	
	//runlog.Println("初始化leveldb")
	//leveldbDataPath := c.LogPath + DataPath
	//fileInfo, err = os.Stat(leveldbDataPath)
	//if err != nil {
	//	if os.IsNotExist(err) {
	//		err = os.MkdirAll(leveldbDataPath, FileMode)
	//		if err != nil {
	//			runlog.Panicln("创建目录失败 " + err.Error())
	//		}
	//	} else {
	//		runlog.Panicln(err.Error())
	//	}
	//}
	//if !fileInfo.IsDir() {
	//	runlog.Panicln(leveldbDataPath + " 不是目录")
	//}
	//mp.logdb, err = leveldb.OpenFile(leveldbDataPath, nil)
	//if err != nil {
	//	runlog.Panicln(err.Error())
	//}
	
	runlog.Println("初始化下载器")
	mp.downloader = downloader.NewHttpDownLoader(1000, 10000)
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
		//mp.pushUserToDownLoader()
		mp.pushWorksAppendixToDownLoader()
		mp.pushActivityToDownLoader()
		mp.downloader.Done()
	}()
	mp.downloader.DownLoad()
	return nil
}
