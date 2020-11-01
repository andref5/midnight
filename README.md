# midnight

Golang middleware Kong plugin for study purposes only.

## Build
```
go build -buildmode plugin midnight.go
```

## Test
```
cd sample
docker-compose up

curl -v localhost:8000/change/xy7-54
```