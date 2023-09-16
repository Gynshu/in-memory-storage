# In-Memory Storage

This is a simple in-memory storage service that allows clients to add, delete, and retrieve key-value pairs via a REST API. The service also includes rate limiting functionality to limit the number of requests per second from a single IP address.

## API

The API includes the following endpoints:

- `POST /set`: Add a new key-value pair to the storage. The request body should include a JSON object with the key and value fields. An optional `expiration` field can be included to set a time-to-live value for the key in seconds.
- `DELETE /delete?key=`: Delete the key-value pair with the specified key from the storage.
- `GET /get?key=`: Retrieve the value for the key with the specified key from the storage.
- `GET /all`: Retrieve all key-value pairs from the storage.

Object should be in the following format:

```json
{
    "key": "key",
    "value": "value",
    "expiration": 60
}
```
as in this Entity struct:
```go
// Entity represents a key-value pair in the in-memory storage.
type Entity struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	// Expiration is the time in nanoseconds when the key-value pair will expire.
	Expiration int64 `json:"expiration"`
}
```

If a client exceeds the rate limit, the service will return a `429 Too Many Requests` HTTP status code.

## Configuration

The service can be configured using the `environment` variables listed below:<br>
`SERVER_PORT`  server port, default 8080 <br>
`RATE_LIMIT`  rate limit in requests per second, default 10 <br>
This value is used to calculate how many nanoseconds to wait between requests from a single IP address. <br>
## Building and Running

To build the service, run the following command:

```
go build -o in-memory-storage cmd/in-memory-storage/main.go
```

To run the service, use the following command:

```
./in-memory-storage
```

The service can also be run using Docker. To build the Docker image, use the following command:

```
docker build -t in-memory-storage .
```

To run the Docker container, use the following command:

```
docker run -p 8080:8080 in-memory-storage
```

## Testing

To run the tests for the service, use the following command:

```
go test ./...
```