Solution to Golang Assignment - https://www.notion.so/playment/Golang-assignment-6e6df82a03024d9c92d44821b666001c

## Prerequisite
Set the optional environment variables `EXPIRY_TIME` (unset default is 0 seconds) and `PRICE_TRACKER` (unset default 
is `coindesk`).
```
export EXPIRY_TIME=<Expiry time in seconds>
export PRICE_VARIABLE=coindesk
```
The `EXPIRY_TIME` variable determines how long the fetched exchange data is valid.
The `PRICE_TRACKER` variable sets which online API to call to get the exchange rates.
For now, it can be set to `coindesk`, new values will be added with support for other APIs.

## Execution
Change to the `src` directory and execute the server -
```
go run .
``` 

Run the client on a different terminal -
```
curl http://localhost:9999/prices
```
To get the results on the terminal.