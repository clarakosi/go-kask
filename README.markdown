go-kask
=======

Prototype of Kask in Golang.

First
-----

    $ cqlsh -f cassandra_schema.cql
    $ go run kask.go storage.go

If necessary, you can pass environment variables for any of `CASSANDRA_HOST`,
`CASSANDRA_PORT`, `CASSANDRA_KEYSPACE`, or `CASSANDRA_TABLE`.

Then
----

    $ curl -D - -X POST http://localhost:8080/sessions/v1/foo -d 'bar'
    HTTP/1.1 200 OK
    Content-Type: application/octet-stream
    Date: Tue, 11 Dec 2018 22:50:46 GMT
    Content-Length: 0
    
    $ curl -D - -X GET  http://localhost:8080/sessions/v1/foo; echo
    HTTP/1.1 200 OK
    Content-Type: application/octet-stream
    Date: Tue, 11 Dec 2018 22:51:10 GMT
    Content-Length: 3
    
    bar
    $
