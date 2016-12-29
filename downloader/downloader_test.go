package downloader

import (
	"testing"
	"fmt"
)

func TestHttpDownLoader(t *testing.T) {
	dl := NewHttpDownLoader(1000, 10)
	
	go func() {
		defer dl.Done()
		for i := 0; i < 1000000; i++ {
			toFile := fmt.Sprintf(
				"downloader_test/%d.jpg", i,
			)
			dl.Push(&Task{
				Url : "https://www.baidu.com/img/bd_logo1.png",
				To : toFile,
			})
		}
	}()
	
	dl.DownLoad()
}
