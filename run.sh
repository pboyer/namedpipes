go build cmd/writer/writer.go
go build cmd/reader/reader.go
go run cmd/master/master.go -writer=`pwd`/writer -reader=`pwd`/reader -numReaders=3