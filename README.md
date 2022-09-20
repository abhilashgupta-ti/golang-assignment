Solution to Golang Assignment - https://www.notion.so/playment/Golang-assignment-6e6df82a03024d9c92d44821b666001c

To run,  change directory to `src` and then execute -
```
go mod init //One time step
```
followed by
```
export EXPIRY=<No of seconds you are willing to let the data be stale for.> //change according to requirements
go run .
``` 

Now change terminal and run -
```
curl http://localhost:9999/prices
```
To get the results on the terminal.