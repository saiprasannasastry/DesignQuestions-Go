//sample.go just defines a sample cases with 4 or 5 tasks added to give the user a better understand
//in ideal world our programe keeps running and only in case of error we will exit
//for that we can have a err case in select loop. is not a valid scenario in sample because
// of assumptions made to generate taskStatus

// hi from Executor
// &{1 false timeout 2021-01-01 21:16:08.835477 +0530 IST m=+0.000213537 data1 4}
// &{2 false timeout 2021-01-01 21:16:08.835486 +0530 IST m=+0.000222231 data2 6}
// &{3 false completed 2021-01-01 21:16:08.835487 +0530 IST m=+0.000222905 data3 5}
// &{4 false timeout 2021-01-01 21:16:08.835487 +0530 IST m=+0.000223358 data4 1}
// &{5 false failed 2021-01-01 21:16:08.835488 +0530 IST m=+0.000223875 data5 2}

//for timeout >3 we delete the items from queue by logging the timestamp when the item was removed

//hi from cleaner
// 2021/01/01 21:16:19 Removed data data1 with id 1
// 2021/01/01 21:16:19 Removed data data2 with id 2
// 2021/01/01 21:16:19 Removed data data3 with id 4

//2021/01/01 21:16:19 done cleaning
// 0 //len(list after cleaning)

//Another output where cleaner is waiting when cleaner is running and adder is not adding but running
// hi from E
// &{1 true completed 2021-01-01 22:02:09.137094 +0530 IST m=+0.000259095 data1 0}
// &{2 true completed 2021-01-01 22:02:09.137104 +0530 IST m=+0.000268878 data2 0}
// &{3 true completed 2021-01-01 22:02:09.137105 +0530 IST m=+0.000269284 data3 0}
// &{4 true completed 2021-01-01 22:02:09.137105 +0530 IST m=+0.000269562 data4 0}
// &{5 false timeout 2021-01-01 22:02:09.137105 +0530 IST m=+0.000269843 data5 7}
// hi from clean
// 2021/01/01 22:02:20 removing the Task 5 with description data5 because timeout was 7
// waiting for added to add tasks
// 2021/01/01 22:02:20 done cleaning

package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

//Task defines the data types the struct holds
type Task struct {
	Id           int
	IsCompleted  bool
	Status       string    // untouched, completed, failed, timeout
	CreationTime time.Time // when was the task created
	TaskData     string    // field containing data about the task
	Timeout      int
}

var list []*Task
var completed = [2]bool{true, false}
var status = [3]string{"completed", "failed", "timeout"}
var taskQueue = make(chan *Task, 50)
var k = 0
var completedTask int
var mutex = &sync.Mutex{}

func main() {
	done := make(chan bool)
	c := make(chan os.Signal, 2)
	ticker := time.NewTicker(3 * time.Second)

	rand.Seed(time.Now().UnixNano())
	// adder keeps adding tasks to the taskQueue in the default state
	go Adder()

	//watching to see every x seconds if something is being added to queue to start the cleaner
	go func() {
		for {
			select {
			case <-ticker.C:

				Monitor(done)
			}
		}
	}()

	//Executor picks the task from the Taskqueue does further processing and adds status and completed status
	go Executor(taskQueue)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case <-c:
		log.Println(" Received signal Interrupt stopping the server")
		break
	case <-done:
		ticker.Stop()
		log.Printf("done cleaning ")
		fmt.Println(len(list))
		break
	}

}

//Monitor runs every x seconds and we start the cleaner to clean only if something gets added to queue
//This way this allows both go routines to run independently and run their task forever
func Monitor(done chan bool) {
	if len(list) > 0 {
		Cleaner(done)
	}
}

//Executor picks the tasks from the added queue and assigns status based on processing
// here we are assiging random status but can have our logic to add status as whatever
func Executor(taskQueue chan *Task) {
	fmt.Println("hi from E")
	for i := 0; i <= 5; i++ {
		task := <-taskQueue
		// generating it randomly, but we can have actual scenario here to updated the status
		randmomStatus := status[rand.Intn(len(status))]
		if randmomStatus == "completed" {
			task.IsCompleted = true
			completedTask++
		}
		if randmomStatus == "timeout" {
			task.Timeout = rand.Intn(10)
		}
		task.Status = randmomStatus
		fmt.Println(task)
		mutex.Lock()
		list = append(list, task)
		mutex.Unlock()

	}
}

//Adder just adds tasks to TaskQueue, here we have just implemented a for loop
// we can keep a list of items such as {like,comment,share} and generate radom values
// and add to taskQueue if only liked or commented.
func Adder() {

	for id := 1; id <= 5; id++ {

		work := &Task{Id: id, TaskData: "data" + fmt.Sprint(id), Status: "untouched", CreationTime: time.Now()}
		taskQueue <- work
	}
}

//Cleaner runs after something gets added into the queue, if the task is completed
func Cleaner(done chan bool) {
	time.Sleep(8 * time.Second)
	fmt.Println("hi from clean")
	//worst case scenario where all are not completed
	for i := 0; i < 10; i++ {

		if len(list) == 0 || k == len(list) {
			fmt.Println("waiting for added to add tasks")
			break
		}
		// never going to happen, just a check to mitigate panic
		if k > len(list) {
			k--
		}
		element := list[k]
		if element.IsCompleted == true {
			removeIndex(k)
			//k = k + 1
		} else {
			if element.Status == "timeout" && element.Timeout > 3 {
				removeIndex(k)
				log.Printf("removing the Task %v with description %v because timeout was %v", element.Id, element.TaskData, element.Timeout)
			} else {
				removeIndex(k)
				// sudo case just making this status as completed
				//asumption that after retry task got completed
				element.IsCompleted = true
				element.Status = "completed"
				completedTask++
				mutex.Lock()
				list = append(list, element)
				mutex.Unlock()
				//log.Printf("Removed data %v with id %v ", element.TaskData, element.Id)
			}
		}
	}
	done <- true

}

func removeIndex(index int) {
	mutex.Lock()
	list = append(list[:index], list[index+1:]...)
	mutex.Unlock()

}
