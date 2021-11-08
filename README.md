# Start Dependencies

```shell
docker-compose up -d
```

# Start Server

```shell
go run main.go
```

# Demo

```shell
curl 'http://localhost:8091/hello?name=world'
# output
{"result":"Hello: world!"}
```
