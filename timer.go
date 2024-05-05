package Timer

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type Timer interface {
	// FindCronList 寻找所有Cron
	FindCronList() map[string]*taskManager
	// AddTaskByFuncWithSecond 添加Task 方法形式以秒的形式加入
	AddTaskByFuncWithSecond(cronName string, spec string, fun func(), taskName string, option ...cron.Option) (cron.EntryID, error) // 添加Task Func以秒的形式加入
	// AddTaskByJobWithSeconds 添加Task 接口形式以秒的形式加入
	AddTaskByJobWithSeconds(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error)
	// AddTaskByFunc 通过函数的方法添加任务
	AddTaskByFunc(cronName string, spec string, task func(), taskName string, option ...cron.Option) (cron.EntryID, error)
	// AddTaskByJob 通过接口的方法添加任务 要实现一个带有 Run方法的接口触发
	AddTaskByJob(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error)
	// FindCron 获取对应taskName的cron 可能会为空
	FindCron(cronName string) (*taskManager, bool)
	// StartCron 指定cron开始执行
	StartCron(cronName string)
	// StopCron 指定cron停止执行
	StopCron(cronName string)
	// FindTask 查找指定cron下的指定task
	FindTask(cronName string, taskName string) (*task, bool)
	// RemoveTaskById 根据id删除指定cron下的指定task
	RemoveTaskById(cronName string, id int)
	// RemoveTaskByName 根据taskName删除指定cron下的指定task
	RemoveTaskByName(cronName string, taskName string)
	// Clear 清理掉指定cronName
	Clear(cronName string)
	// Close 停止所有的cron
	Close()
}

// task 任务
type task struct {
	EntryID  cron.EntryID // 任务ID
	Spec     string       // 任务规则
	TaskName string       // 任务名称
}

// taskManager 任务管理
// Every taskManager has many task,task is the atom of taskManager
// Every taskManager just like Constituencies
type taskManager struct {
	corn  *cron.Cron             // 任务的定时调度
	tasks map[cron.EntryID]*task // 任务列表
}

// Timer 定时任务管理
type timer struct {
	cronList map[string]*taskManager // 定时任务列表
	sync.Mutex
}

func (t *timer) FindCronList() map[string]*taskManager {
	return t.cronList
}

func (t *timer) AddTaskByFunc(cronName string, spec string, fun func(), taskName string, option ...cron.Option) (cron.EntryID, error) {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.cronList[cronName]; !ok {
		tasks := make(map[cron.EntryID]*task)
		t.cronList[cronName] = &taskManager{
			corn:  cron.New(option...),
			tasks: tasks,
		}
	}
	id, err := t.cronList[cronName].corn.AddFunc(spec, fun)
	t.cronList[cronName].corn.Start()
	t.cronList[cronName].tasks[id] = &task{
		EntryID:  id,
		Spec:     spec,
		TaskName: taskName,
	}
	return id, err
}

func (t *timer) AddTaskByFuncWithSecond(cronName string, spec string, fun func(), taskName string, option ...cron.Option) (cron.EntryID, error) {
	t.Lock()
	defer t.Unlock()
	option = append(option, cron.WithSeconds())

	if _, ok := t.cronList[cronName]; !ok {
		tasks := make(map[cron.EntryID]*task)
		t.cronList[cronName] = &taskManager{
			corn:  cron.New(option...),
			tasks: tasks,
		}
	}

	id, err := t.cronList[cronName].corn.AddFunc(spec, fun)
	t.cronList[cronName].corn.Start()
	t.cronList[cronName].tasks[id] = &task{
		EntryID:  id,
		Spec:     spec,
		TaskName: taskName,
	}

	return id, err
}

func (t *timer) AddTaskByJobWithSeconds(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error) {
	return t.AddTaskByFuncWithSecond(cronName, spec, job.Run, taskName, option...)
}

func (t *timer) AddTaskByJob(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error) {
	return t.AddTaskByFunc(cronName, spec, job.Run, taskName, option...)
}

func (t *timer) FindCron(cronName string) (*taskManager, bool) {
	t.Lock()
	defer t.Unlock()
	v, ok := t.cronList[cronName]
	return v, ok
}

func (t *timer) StartCron(cronName string) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.cronList[cronName]; ok {
		v.corn.Start()
	}
}

func (t *timer) StopCron(cronName string) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.cronList[cronName]; ok {
		v.corn.Stop()
	}
}

func (t *timer) RemoveTaskById(cronName string, id int) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.cronList[cronName]; ok {
		v.corn.Remove(cron.EntryID(id))
		delete(v.tasks, cron.EntryID(id))
	}
}

func (t *timer) FindTask(cronName string, taskName string) (*task, bool) {
	t.Lock()
	defer t.Unlock()
	v, ok := t.cronList[cronName]
	if !ok {
		return nil, false
	}
	for _, t1 := range v.tasks {
		if t1.TaskName == taskName {
			return t1, true
		}
	}
	return nil, false
}

func (t *timer) RemoveTaskByName(cronName string, taskName string) {
	fk, ok := t.FindTask(cronName, taskName)
	if !ok {
		return
	}
	t.RemoveTaskById(cronName, int(fk.EntryID))
}

func (t *timer) Clear(cronName string) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.cronList[cronName]; ok {
		v.corn.Stop()
		delete(t.cronList, cronName)
	}
}

func (t *timer) Close() {
	t.Lock()
	defer t.Unlock()
	for _, v := range t.cronList {
		v.corn.Stop()
	}
}

func NewTimer() Timer {
	return &timer{
		cronList: make(map[string]*taskManager),
	}
}
