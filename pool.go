package main

import (
	"fmt"
	"time"
)

//定义任务类型Task
type Task struct {
	f func() error
}

//创建一个Task任务
func NewTask(argF func() error) *Task {
	t := Task{
		f : argF,
	}

	return &t
}

//Task也需要一个执行业务的方法
func (t *Task) Execute(){
	//调用任务中绑定好的业务方法
	t.f()
}


/*
	创建携程池相关业务
*/

//定义一个Pool携程池
type Pool struct {
	//对外的Task入口 EntryChannel
	EntryChannel chan *Task

	//内部的Task队列 JobsChannel
	JobChannel chan *Task

	//协程池最大worker数量
	workerNum int
}

//创建goroutine Pool
func NewPool(cap int) *Pool {
	//创建一个Pool
	p := Pool {
		EntryChannel: make(chan *Task),
		JobChannel: make(chan *Task),
		workerNum: cap,
	}

	return &p
}

//协程池创建一个worker,并且让这个worker去工作
func (p *Pool) worker(workerID int) {
	//一个worker具体的工作

	//1 永久的从JobsChannel中去取任务
	for task := range p.JobChannel{
		//task 就是当前worker从JobsChannel
		task.Execute()
		fmt.Println("workerID",workerID,"执行完了一个任务")
	}
}

//协程池启动一个方法
func (p *Pool) run() {
	//1 根据workerNUM来创建worker去工作
	for i := 0; i < p.workerNum; i++ {
		//为每个worker分配一个goroutine
		go p.worker(i)
	}

	//2 从EntryChannel中取任务，取到的任务发送给JobChannel
	for task := range p.EntryChannel{
		//循环读取外部队列的任务到内部队列中
		p.JobChannel <- task
	}
}

//测试协程池工作
func main(){
	//1 创建一些任务
	t := NewTask(func() error {
		//当前任务的业务，打印出当前的系统时间
		fmt.Println(time.Now())
		return nil
	})

	//2 创建一个pool协程池，协程池最大的worker数量是4
	p := NewPool(4)

	taskNum := 0

	//3 将这些任务，提交给协程池处理
	go func() {
		for {
			//不断的向p中去写入任务t，每个任务就是打印当前的时间
			p.EntryChannel <- t
			taskNum += 1
			fmt.Println("当前一共执行了", taskNum, "任务")
		}
	}()

	p.run()
}
