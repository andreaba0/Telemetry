protoc --go_out=Agent/proto --go_opt=paths=source_relative --go-grpc_out=Agent/proto --go-grpc_opt=paths=source_relative metrics.proto

protoc --go_out=Datastore/proto --go_opt=paths=source_relative --go-grpc_out=Datastore/proto --go-grpc_opt=paths=source_relative metrics.proto
