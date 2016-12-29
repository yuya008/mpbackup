package downloader

import (
	"io"
	"log"
	"os"
	"sync"
	"net/http"
	"io/ioutil"
	"path"
	"errors"
	"time"
)

type httpDownLoader struct {
	threadN int
	taskC chan *Task
	logger *log.Logger
}

const (
	LogPrefix = "downloader"
	LogFlag = log.Ldate | log.Lmicroseconds
)

func NewHttpDownLoader(threadN, taskCLen int) DownLoader {
	return &httpDownLoader{
		threadN : threadN,
		taskC : make(chan *Task, taskCLen),
		logger : log.New(os.Stderr, LogPrefix, LogFlag),
	}
}

func (hdl *httpDownLoader) Push(task *Task) *Task {
	if task == nil || task.Url == "" || task.To == ""{
		return nil
	}
	for {
		select {
		case hdl.taskC <- task:
			hdl.logger.Printf("推送任务 %v 到任务队列\n", task)
			return task
		default:
			hdl.logger.Println("任务队列已满重试")
			time.Sleep(time.Millisecond * 100)
			continue
		}
	}
}

func (hdl *httpDownLoader) Done() {
	hdl.logger.Println("完成任务推送")
	close(hdl.taskC)
}

func (hdl *httpDownLoader) SetLoggerWriter(w io.WriteCloser, flag int) {
	hdl.logger = log.New(w, LogPrefix, flag)
}

func (hdl *httpDownLoader) DownLoad() {
	var wait sync.WaitGroup
	wait.Add(hdl.threadN)
	
	for i := 0; i < hdl.threadN; i++ {
		hdl.logger.Printf("派发线程 %d \n", i)
		go func(thrno int) {
			defer wait.Done()
			for {
				task, ok := <- hdl.taskC
				if !ok {
					hdl.logger.Printf("线程 %d -> 完成任务退出\n", thrno)
					return
				}
				hdl.logger.Printf("线程 %d -> 抢到任务 %v\n", thrno, task)
				err := hdl.download(task)
				if err != nil {
					hdl.logger.Printf("线程 %d -> %s\n", thrno, err.Error())
				}
			}
		}(i)
	}
	
	wait.Wait()
	hdl.logger.Println("downloader退出")
}

func (hdl *httpDownLoader) download(task *Task) error {
	dir := path.Dir(task.To)
	
	if fileinfo, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			hdl.logger.Println(dir, "目录不存在,创建")
			if err = os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	} else if !fileinfo.IsDir() {
		return errors.New(dir + "不是一个目录")
	}
	
	hdl.logger.Println("开始下载", task.Url)
	response, err := http.Get(task.Url)
	
	if err != nil {
		return err
	}
	
	data, err := ioutil.ReadAll(response.Body)
	
	if err != nil {
		return err
	}
	
	if err = ioutil.WriteFile(task.To, data, 0755); err != nil {
		return err
	}
	hdl.logger.Println("下载成功", task.Url, task.To)
	return nil
}
