# Code Challenge - Authorizer
You are tasked with implementing an application that authorizes a transaction for a specific account following a set of
predefined rules.

# Why Go? or Golang?
I really love this langauge, and I think it offers a great tradeoff between performance, simplicity and readability.

# How to build?
- Run `make build` to compile the project directly on your OS, binary will be added to build directory, it requires go 1.16+.
- Run `make docker-build` to create a container, it requires docker.

# How to run an example?
- Run `make run` to execute an example using the data from `testdata/operations`.
- Run `make docker-run` to execute an example on a container using the data from `testdata/operations`.

The output in both cases must contain like this:

```
{"account":{"activeCard":true,"availableLimit":1000},"violations":[]}
{"account":{"activeCard":true,"availableLimit":900},"violations":[]}
{"account":{"activeCard":true,"availableLimit":800},"violations":[]}
{"account":{"activeCard":true,"availableLimit":700},"violations":[]}
{"account":{"activeCard":true,"availableLimit":600},"violations":[]}
{"account":{"activeCard":true,"availableLimit":500},"violations":[]}
{"account":{"activeCard":true,"availableLimit":400},"violations":[]}
{"account":{"activeCard":true,"availableLimit":400},"violations":["insufficient-limit"]}
{"account":{"activeCard":true,"availableLimit":400},"violations":["insufficient-limit"]}
{"account":{"activeCard":true,"availableLimit":400},"violations":["insufficient-limit"]}
```

# How to run with custom data files?
- Run `./build/authorizer < testdata/sample` to execute directly on your OS.
- Run `cat testdata/sample |  docker run -i authorizer` to execute on a docker.

On both cases `testdata/sample` represents the file which contains your data.

# How to run tests?
Tests run on local OS, so you require go 1.16+.
- `make unit-test` executes unit tests using golang testing package, shows coverage percentage after execution and packages tested (Some packages are being skipped because they don't contain functions to test).
- `make integration-test` executes integration tests, starts DB and Service to execute examples and compare them to the expected output.

# Makefile
I created a Makefile to build and run the project, but it has other commands available, you can run `make help` to list available commands:

```
  $  make help

 Choose a command run in :

  build                 Compile project.
  run                   Build and execute project using testdata/operations.
  unit-test             Execute unit tests.
  integration-test      Execute integration tests.
  lint                  Execute linter using the rules in .golangci.yml.
  clean                 Delete build directory and log.
  docker-build          Build project and create container with it.
  docker-run          Execute the application within a container with testdata.
  update-dependencies   Update all golang dependencies.
```

# Code design choices

## Project structure

Project structure and a small description of the most important directories and go files.

There are comments in most of the packages if you need more info or something in particular.

```
.
|-- cmd
|   |-- authorizer ------------ Main package
|   |   |-- integration_test.go - Integration tests, similar to main initializes dependencies and tests application
|   |   |-- main.go ------------- main() func initializes dependencies and runs the application
|   |   `-- testdata ------------ Testdata used by integration tests
|-- Dockerfile
|-- go.mod
|-- go.sum
|-- internal -------------------- All the code in internal directory cannot be imported by other applications
|   |-- app
|   |   |-- model --------------- Declaration of the model or structs needed across the application
|   |   |   `-- model.go 
|   |   |-- service ------------- Implements most of the logic of the operations createAccount and Transaction
|   |   |   |-- parser.go ------- Parses the stdin to get the json required by the application
|   |   |   |-- parser_test.go
|   |   |   |-- rules------------ Business rules implemented as funcs and used by service package
|   |   |   |   |-- rules.go
|   |   |   |   `-- rules_test.go
|   |   |   |-- service.go ------- Service implements most of the logic used to execute the operations
|   |   |   `-- service_test.go
|   |   |-- storage -------------- Implements the database logic
|   |   |   |-- inmemory.go
|   |   |   `-- inmemory_test.go
|   |   `-- violations ----------- Violations declared as constants
|   |       `-- violations.go
|   `-- common ------------------- Common functions not directly related to this application
|   |    `-- logfile
|   |        `-- logfile.go
|   `-- root --------------------- Package that controls the flow of the application, reads the lines from stdin and decide which service operation to execute
|       |-- reader --------------- Gets the string and unmarshals it to a struct for both operations
|       |   |-- parser.go
|       |   `-- parser_test.go
|       |-- root.go
|       `-- root_test.go
|-- Makefile
|-- README.md
|-- scripts ---------------------- All the scripts used by Makefile
`-- testdata --------------------- Data used for execution or as an example
    |-- operations
    `-- sample
```

## Interfaces
### Simpler and better unit tests
I used interfaces across several parts of the code because with them you can simulate dependencies a lot easier and simplify unit testing a lot, for example using the Storage interface we can simulate the operations of any kind of storage and concentrate on the logic of the function without worrying about starting a DB. Interfaces can also help us to add new Storage types without worrying about the caller, as long as we respect the methods, everything can work smoothly as always.

This choice was really useful for testing Service package because it contains the business logic of the operations CreateAccount and Transaction, and we can test them without really connecting to a DB, we just create a mockup that fulfills the Storage interface, and we can simulate any business case needed for the case.
```
type Storage interface {
    CreateAccount(a model.Account) error
    GetAccount(aID int) model.Account
    ExecuteTransaction(a model.Account, t model.Transaction) (model.Account, error)
    GetTransactions(accountID int) []model.Transaction
    Close() error
}
```

## Function Execute() and main()
### Unit test "main" flow
The function Execute controls the flow of the application, receiving the input from an io.Reader and calling the service operations as needed.
The Main() was kept really simple to be able to test most of the logic on the Execute function.

```
func Execute(auth Authorizer, reader io.Reader, writer io.Writer) {
```

## Business Rules as functions
### Implementing Business Rules
To solve the problem of adding the business rules I added a function `func (br *BusinessRule) ExecuteRules() (bool, string)` that calls all the existing business rules, so the callers don't have to change anything in case we add new business rules.

All the functions called by `ExecuteRules` share the same input and outputs to be consistent, I'm sure there must be a more dynamic, fast and user-friendly way to create the business rules but this was the best approach I could think of in the given time.
Something more configurable would be ideal, maybe even something with a front-end and selectable parameters and conditions.

Adding new rules require to compile the project again.

```
func (br *BusinessRule) ExecuteRules() (bool, string) {
    response, violation := br.isActive()
       if !response {
        return response, violation
    }

	response, violation = br.sufficientLimit()
	if !response {
		return response, violation
	}

	response, violation = br.doubleTransaction()
	if !response {
		return response, violation
	}

	response, violation = br.highFrequency()
	if !response {
		return response, violation
	}

	return true, ""
}

```

## Database as Maps
### Simulating a DB with go structures

For the Database I tried to keep it as simple as possible, so I used maps to simulate the tables, one for the account and another one for the transactions, the transactions uses as value a slice of transactions to store all the transactions that were executed.
Both maps use the accountID as key, and in this example we used only one account (id = 1) but with this database design we can add several accounts as needed.

```
// Account in this package represents the table of Accounts in the simulated DB
type Account struct {
	Id             int
	ActiveCard     bool
	AvailableLimit int
}

// Transaction in this package represents the table of Transactions in the simulated DB
type Transaction struct {
	Id       uuid.UUID
	Merchant string
	Amount   int
	Time     time.Time
}
```