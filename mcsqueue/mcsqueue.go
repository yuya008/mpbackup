package mcsqueue

import (
	"sync"
	"fmt"
)

// mcsqueue对外接口
type Mcsqueue interface {
	Put()
	Get()
}

type mcschan struct {
	c chan interface{}
	mutex sync.Mutex
}

type mcsqueue struct {
	chanN 		int32
	perChanSize int
	chans		[]mcschan
}

func New() Mcsqueue {
	return &mcsqueue{
		
	}
}

func (mcs *mcsqueue) Put() {
	fmt.Printf("哈哈")
}

func (mcs *mcsqueue) Get() {
	
}
