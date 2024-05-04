package Timer

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type Timer interface {
	// FindCronList 寻找所有的定时任务
	FindCronList() map[string]*taskManager
	// AddTasksByFuncWithSecond 添加Task , 以秒的形式/函数
	AddTasksByFuncWithSecond(cronName string, spec string, fun func(), taskName string, option ...cron.Option) (cron.EntryID, error)
	// AddTasksByJobWithSecond  添加Task , 以秒的形式/接口
	AddTasksByJobWithSecond(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error)
	// AddTaskByFunc 添加Task 函数
	AddTaskByFunc(cronName string, spec string, fun func(), taskName string, option ...cron.Option) (cron.EntryID, error)
	// AddTaskByJob 添加Task 接口
	AddTaskByJob(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error)
	// FindCronByName 通过taskName 获取对应的cron
	FindCronByName(cronName string) (*taskManager, bool)
	//
	StartCron(cronName string)
	StopCron(cronName string)

	//
	RemoveTaskById(cronName string, id int)
	RemoveTaskByName(cronName string)
}

// task 任务
type task struct {
	EntryID  cron.EntryID // 任务ID
	Spec     string       // 任务规则
	TaskName string       // 任务名称
}

// taskManager 任务管理
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

func (t *timer) AddTaskByJob(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error) {
	return t.AddTaskByFunc(cronName, spec, job.Run, taskName, option...)
}

func (t *timer) AddTaskByJobWithSecond(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error) {
	return t.AddTaskByFuncWithSecond(cronName, spec, job.Run, taskName, option...)
}

func (t *timer) FindCron(cronName string) (*taskManager, bool) {
	t.Lock()
	defer t.Unlock()
	v, ok := t.cronList[cronName]
	return v, ok
}

func (t *timer) FindTasks(cronName string, taskName string) (*task, bool) {}

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

func (t *timer) RemoveTaskByName(cronName string) {}
