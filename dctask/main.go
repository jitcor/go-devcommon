package dctask

import "sync"

type DcTask struct {
	tasks []func(*error)
	wg    *sync.WaitGroup
}

func NewDcTask() *DcTask {
	return &DcTask{
		tasks: make([]func(*error),0),
		wg: new(sync.WaitGroup),
	}
}
func (that *DcTask) Go(task func(*error),err *error) *DcTask {
	if *err!=nil{
		return that
	}
	go func() {
		that.wg.Add(1)
		defer that.wg.Done()
		task(err)
	}()
	that.tasks= append(that.tasks, task)
	return that
}
func (that DcTask) Wait() {
	that.wg.Wait()
}