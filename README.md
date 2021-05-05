# study gRPC

## proto compile
`protoc -Iproto --go_out=plugins=grpc:api proto/*.proto`

## run server
```sh
cd api
go run server/server.go
```

