all: protocol server

protocol:
	mkdir ./protocols && thrift --gen go -out . snowflake.thrift	# -thrift_import option

server:
	go build -o server server.go handler.go common.go

client:
	go build -o client client.go common.go

clean:
	rm -rf ./protocols
	rm server client
