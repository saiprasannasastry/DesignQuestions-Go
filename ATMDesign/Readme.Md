# Design a ATM
* Design an interactive atm to 
    * see balance
    * WithDraw cash
    * Post withdrawl if he wants to see the balance or not 
    
## To Test ##
`go run main.go`
* This creates a new instance everytime, so everytime you run, it reinitilizes the persons bank balance

## Added a new client server model using grpc
* hosting the server in docker
* to run the server run `docker-compose up` from the docker directory
* run the client using `go run cmd/client/main.go` or create a binary and run it
* The client shows the following ops/ interactive way
   * Register  - create a new user into the bank.If successful generates a pin
   * WithDrawCash - remove cash f rom the bank
   * CheckBalance - look at existing balance
   * DepositMoney - add money to the bank
   * Exit - exit the client

