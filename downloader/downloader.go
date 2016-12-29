package downloader

import (
	"io"
)

const (
	Version = "0.0.1"
)

type DownLoader interface {
	// 推送任务
	Push(*Task) *Task
	// 结束推送任务
	Done()
	// 开始下载
	DownLoad()
	// 设置日志写位置
	SetLoggerWriter(io.WriteCloser, int)
}

type Task struct {
	Url			string
	To			string
}
