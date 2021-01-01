* sample.go just defines a sample cases with 4 or 5 tasks added to give the user a better understanding of the problem.
* in ideal world our programe keeps running and only in case of error we will exit.for that we can have a err case in select loop. it is not a valid scenario in sample because of assumptions made to generate taskStatus

* to run sampler just run `go run sample/sample.go`
* to run actual problem `go run main.go`

main.go never dies until a signal interrupt `ctl+c` is given, main.go has only the log statment that the problem statement mentions

There are other alternative ways to achieve if there was some mode of communication between both the go routines, or we used external packages but since the problem statement explicitely mentioned that the go routines should no communicate, had to add a ticker to monitor if something was added to queue , to start the cleaner so that we don't panic

TaskQueue in our code is defined as `var list[]*Task`

## OTHER solutions ##
* we can have Executor and Cleaner running as seperate microsevices  and have a database as a microservice
  * Executor microservice just adds to the database
  * in the database we can make the table to  notifiy clener service
  * the cleaner service can then check the databse , and based on event status in database we can update the Task in DB or clear it
  
* We can have a messaging Queue such as rabbitMq /zmq /PUBSUB
  * Addder adds to workerqueue
  * executor picks tasks from queue and executs the tasks and pushs the task data into TaskQueue
  * cleaner can subscribe on that topic and based on event status we can readd to queue or remove the task from TaskQueue
  
 * we can have both go routines communicating 
   * executor pushes data into which cleaner channel is listening, after fetching from cleaner channel , we can delete from taskQueue(list) or updateTaskQueue(list)
  


