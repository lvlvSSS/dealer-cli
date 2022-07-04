package schedule

import "sync"

type Schedule interface {
	Handle(wg *sync.WaitGroup) error
}
