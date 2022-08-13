# study gRPC

## setup
### install protoc
```sh
sudo apt install -y protobuf-compiler
```

### server
```sh
cd api
go build
```

### client
```sh
cd client
buddle install
```

## compile proto
### server 
```sh
protoc -Iproto --go_out=plugins=grpc:api proto/*.proto
```

### client
```sh
cd client
bundle exec grpc_tools_ruby_protoc -I ../proto --ruby_out=app/gen/api/pancake/baker --grpc_out=app/gen/api/pancake/baker ../proto/pancake.proto
```

## run server
```sh
cd api
go run server/server.go
```

