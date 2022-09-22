Solution to Golang Assignment - https://www.notion.so/playment/Golang-assignment-6e6df82a03024d9c92d44821b666001c

## Prerequisite

Change to the `src` directory and execute once -
```
go mod init //One time step
```

## Execution
To set the `EXPIRY_TIME` variable -
```
export EXPIRY_TIME=<Expiry time in seconds>
```
Optionally, also set the `PRICE_TRACKER` variable.
This variable sets which API to call to get the exchange rates.
For now, it can be set to `coindesk` (also the default if this variable is not set), 
new values will be added with support for other APIs.
```
export PRICE_VARIABLE=coindesk
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