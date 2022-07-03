package schedule

import "sync"

type Remote interface {
	Handle(wg *sync.WaitGroup) error
}
