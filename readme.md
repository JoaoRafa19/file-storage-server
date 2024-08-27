# File storage Server

## Testing

Test:
```shell
go test ./... 
```

Coverage:
```shell
go test -coverprofile cover.out && go tool cover -html=cover.out
```

### Diferences betwen the original code


- TCP Error

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

- Check File existance

Original 
```go
_,err := os.Stat(pathKey.FullPath())
if err == fs.ErrNotExist {
	return false
}
return true

```

Simplified
```go 
_, err := os.Stat(Patkey.FullPath())
return errors.Is(err, fs.ErrNotExist)
```