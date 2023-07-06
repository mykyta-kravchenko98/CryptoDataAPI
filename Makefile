PHONY: generate
generate:
		IF NOT EXIST pkg\\cryptodata_v1 mkdir pkg\\cryptodata_v1
		protoc --go_out=pkg/cryptodata_v1 --go_opt=paths=source_relative \
		--go-grpc_out=pkg/cryptodata_v1 --go-grpc_opt=paths=source_relative \
		api/cryptodata_v1/cryptodata.proto
		move pkg\\cryptodata_v1\\api\\cryptodata_v1\\* pkg\\cryptodata_v1\\
		rmdir /s /q pkg\\cryptodata_v1\\api
