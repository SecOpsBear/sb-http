# Simple Web server

This is a Simple web server written in GO. 

## A simple web server with upload and downloading files

This webserver was designed with penetration testing in mind, where the need for a simple webserver may arise. The code is tested on windows and linux.  

* Simple webserver with custom port assignment.
* Has download and multiple files upload feature.
* List all the file in the current directory.

The upload feature in the server does not perform any security check and it is intentionally done to provide flexibility during the engagement. TLS implementation in progress. 

## How to build   

### Install in linux sb-http

```console
go install github.com/secopsbear/sb-http@latest
```

### Build for linux

```console
go build -o sb-http
```

 > Add **`-ldflags "-s -w"`** flags to reduce the file size by deleting the debug links and info in the binary 

### Reduce the file size

```console
upx --ultra-brute sb-http
```

### Build for window

Generate an executable **`sb-http.exe`** for windows environment.   

```console
env GOOS=windows GOARCH=amd64 go build -o sb-http.exe -ldflags "-s -w"
```

## Example usage

```console
sb-http serve
```


```console
$ sb-http serve -h
Basic server with upload and download feature

Usage:
  sb-http serve [flags]

Flags:
  -h, --help          help for serve
  -p, --port string   Enter the port number (default "8099")
```

## Find a bug?

If you found an issue or would like to submit an improvement to this project, please submit an issue using the issues tab above.