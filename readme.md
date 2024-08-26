## File storage Server


### Diferences betwen the original code

Originaly the TCP Transport handle the decode error by continue looping
after a error but we can type assert the error.

by `var.(type)`

```go
if err != nil {
	t,isOpError := err.(*net.OpError) 
	if isOpError { // Net Closes
		fmt.Println(t)
		return
	}

	fmt.Printf("TCP ERROR: %s\n", err) // Everything else (failure payload)
	continue
}
```