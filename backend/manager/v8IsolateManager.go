package manager

import (
	"fmt"
)

type IsolateStore struct {
	isoTermChans []chan bool
	ch           chan IsolateManagerMessage
}

type IsolateManagerMessage struct {
	isoTerminationChan chan bool
	deleteAll          bool
}

func NewIsolateStore() IsolateStore {
	return IsolateStore{
		isoTermChans: []chan bool{},
		ch:           make(chan IsolateManagerMessage, 30),
	}
}

func (i *IsolateStore) IsolateManager() {
	for message := range i.ch {

		fmt.Println("Isolate", len(i.isoTermChans))

		if message.deleteAll {
			for _, iso := range i.isoTermChans {
				iso <- true
				// iso.Dispose()
			}
			i.isoTermChans = nil
		} else {
			i.isoTermChans = append(i.isoTermChans, message.isoTerminationChan)
		}
	}
}

func (i *IsolateStore) AddIsolate(isoTerminationChan chan bool) {
	i.ch <- IsolateManagerMessage{isoTerminationChan: isoTerminationChan}
}

func (i *IsolateStore) DeleteAllIsolates() {
	i.ch <- IsolateManagerMessage{deleteAll: true}
}
