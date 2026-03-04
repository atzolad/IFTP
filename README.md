# Initialize Golang Project


This is a Golang project using the standard net/http library to build a backend API to serve a student sign up portal / admin class management system for IFTP improv classes. Currenlty undergoing a refactor to remove the Gin framework and transition to the standard net/http. 

Data is stored in a Postgres database. 

The Front-end is built using the Bootstrap library and Html/Css/Javascript. 


To initialize a go project, run:
```
mkdir IFTP && cd IFTP
go mod init IFTP
```

Now we can make a main.go file:
```
touch main.go
```

To initialize the server locally run
```
go run main.go
```