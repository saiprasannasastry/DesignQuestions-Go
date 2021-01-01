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

	c := make(chan os.Signal, 2)
	// using ticker to see if something was added to the queue, making assumption
	// that we should not run cleaner until something gets added in queue
	ticker := time.NewTicker(3 * time.Second)

	//err := make(chan err)

	rand.Seed(time.Now().UnixNano())
	// adder keeps adding tasks to the taskQueue in the default state
	go Adder()

	//watching to see every x seconds if something is being added to queue to start the cleaner
	go func() {
		for {
			select {
			case <-ticker.C:

				Monitor()
			}
		}
	}()

	//Executor picks the task from the Taskqueue does further processing and adds status and completed status
	go Executor(taskQueue)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case <-c:
		log.Println(" Received signal Interrupt stopping the server")
		fmt.Println(len(list))
		ticker.Stop()
		break
		//this scenario is not valid because of assumptions , otherwise both go routines
		// i.e adder and cleaner can send err to the error channel and we can exit
		// case <-err:
		// 	log.Println("programme crashed")
		// 	break

	}

}

//Monitor runs every x seconds and we start the cleaner to clean only if something gets added to queue
//This way this allows both go routines to run independently and run their task forever
func Monitor() {
	if len(list) > 0 {
		Cleaner()
	}
}

//Executor picks the tasks from the added queue and assigns status based on processing
// here we are assiging random status but can have our logic to add status as whatever
func Executor(taskQueue chan *Task) {
	// could take an err chan as an argument, based on processing if some error we could
	// put it in error channel to crash the service

	for {
		task := <-taskQueue
		// generating it randomly, but we can have actual scenario here to updated the status
		randmomStatus := status[rand.Intn(len(status))]
		if randmomStatus == "completed" {
			task.IsCompleted = true
			completedTask++
		}
		// generating a timeout for a timeout case with values between 0-10
		if randmomStatus == "timeout" {
			task.Timeout = rand.Intn(10)
		}
		task.Status = randmomStatus
		mutex.Lock()
		list = append(list, task)
		mutex.Unlock()
	}
}

//Adder just adds tasks to TaskQueue, here we have just implemented a for loop
// we can keep a list of items such as {like,comment,share} and generate radom values
// and add to taskQueue if only liked or commented.
func Adder() {
	action := [3]string{"like", "comment", "share"}
	id := 1
	for {
		randomAction := action[rand.Intn(len(action))]
		if randomAction == "like" || randomAction == "share" {
			work := &Task{Id: id, TaskData: "data" + fmt.Sprint(id), Status: "untouched", CreationTime: time.Now()}
			id++
			taskQueue <- work
		}
	}
}

//Cleaner runs after something gets added into the queue, if the task is completed
func Cleaner() {
	// could take an err chan as an argument, based on processing if some error we could
	// put it in error channel to crash the service
	for {

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
		} else {

			if element.Status == "timeout" && element.Timeout > 3 {
				removeIndex(k)
				log.Printf("removing the Task %v with description %v because timeout was %v", element.Id, element.TaskData, element.Timeout)

			} else {

				removeIndex(k)
				// sudo case just making this status as completed
				//asumption that after retry task got completed, but the same logic which
				//was used in adder and executor could be used to generate
				element.IsCompleted = true
				element.Status = "completed"
				completedTask++
				mutex.Lock()
				list = append(list, element)
				mutex.Unlock()

			}
		}
	}

}

func removeIndex(index int) {
	mutex.Lock()
	list = append(list[:index], list[index+1:]...)
	mutex.Unlock()

}
