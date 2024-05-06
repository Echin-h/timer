package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var job = mockJob{}

type mockJob struct{}

func (job mockJob) Run() {
	mockFunc()
}

func mockFunc() {
	time.Sleep(time.Second)
	fmt.Println("1s...")
}

func TestNewTimerTask(t *testing.T) {
	tm := NewTimer()
	_tm := tm.(*timer)

	{
		id, err := tm.AddTaskByFunc("func", "@every 1s", mockFunc, "测试mockfunc")
		assert.Nil(t, err)
		fmt.Println(id)
		_, ok := _tm.cronList["func"]
		if !ok {
			t.Error("no find func")
		}
	}

	{
		id, err := tm.AddTaskByJob("job", "@every 1s", job, "测试job mockfunc")
		assert.Nil(t, err)
		fmt.Println(id)
		_, ok := _tm.cronList["job"]
		if !ok {
			t.Error("no find job")
		}
	}

	{
		_, ok := tm.FindCron("func")
		fmt.Println("func finding")
		if !ok {
			t.Error("no find func")
		}
		_, ok = tm.FindCron("job")
		fmt.Println("job finding")
		if !ok {
			t.Error("no find job")
		}
		_, ok = tm.FindCron("none")
		if !ok {
			t.Error("find none")
		}
	}

	{
		tm.RemoveTaskById("func", 1)
		_, ok := _tm.cronList["func"]
		if !ok {
			t.Error("not find func")
		}
		_, ok = _tm.FindTask("func", "测试mockfunc")
		if !ok {
			t.Error("not find 测试mockfunc")
		}
	}

	{
		_, ok := _tm.FindTask("func", "测试mockfunc")
		if !ok {
			t.Error("no find func")
		}
		tm.RemoveTaskByName("func", "测试mockfunc")
		_, ok = _tm.FindTask("func", "测试mockfunc")
		if !ok {
			t.Error("no find func")
		}
	}

	{
		list := tm.FindCronList()
		fmt.Println(list)
	}

	{
		tm.Clear("func")
		_, ok := _tm.cronList["func"]
		if !ok {
			t.Error("no find func")
		}
	}

}
