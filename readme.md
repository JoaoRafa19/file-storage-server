# Distributed File Storate 


## TCP transport
Test the connection 
```bash 
go test ./p2p -run TestTCPTransport
```
or run main 
```go 
tr := p2p.NewTCPTransport(":3000")

if err := tr.ListenAndAccept(); err != nil {
	log.Fatal(err)
}
select {}
```
```bash
nc  localhost 4000
```
