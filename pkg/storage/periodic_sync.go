package storage

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type PS interface {
	Stop()
}

type PeriodicSync struct {
	File     *os.File
	Interval time.Duration
	Content  interface{}

	stop     chan struct{}
	stopOnce sync.Once
}

func New(file *os.File, interval time.Duration, content interface{}) PS {
	sync := &PeriodicSync{
		File:     file,
		Interval: interval,
		Content:  content,
		stop:     make(chan struct{}),
	}

	go sync.syncLoop()

	return sync
}

func (sync *PeriodicSync) Stop() {
	sync.stopOnce.Do(func() { close(sync.stop) })
}

func (sync *PeriodicSync) save() error {
	sync.File.Truncate(0)
	sync.File.Seek(0, 0)
	return json.NewEncoder(sync.File).Encode(sync.Content)
}

func (sync *PeriodicSync) syncLoop() {
	interval := sync.Interval
	expireTime := time.NewTimer(interval)
	for {
		select {
		case <-sync.stop:
			return
		case <-expireTime.C:
			err := sync.save()
			if err != nil {
				log.Println("Error saving", err)
			}
			expireTime.Reset(interval)
		}
	}
}
