Solution to Golang Assignment - https://www.notion.so/playment/Golang-assignment-6e6df82a03024d9c92d44821b666001c

## Prerequisite

Change to the `src` directory and execute once -
```
go mod init //One time step
```

## Execution
To set the `EXPIRY` variable -
```
export EXPIRY=<Expiry time in seconds>
```
Run the server -
```
go run .
``` 

Run the client on a different terminal -
```
curl http://localhost:9999/prices
```
To get the results on the terminal.