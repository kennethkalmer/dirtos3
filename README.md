# dirtos3

Upload all files (recursively) in a directory to an S3 bucket.

## Requirements

* [Go](https://golang.org)
* [git](http://git-scm.com)
* An [S3 Bucket](https://aws.amazon.com/S3)

## Building

```sh
$ go get github.com/kennethkalmer/dirtos3
$ cd $GOPATH/src/kennethkalmer/dirtos3
$ go build dirtos3.go
```

## Configuring

Set the following environment variables:

* `AWS_ACCESS_KEY_ID`
* `AWS_SECRET_ACCESS_KEY`
* `S3_BUCKET_NAME`
* `S3_PREFIX`

These should all be blatantly obvious

## Running

```sh
$ ./dirtos3 -source="."
```

## TODO

1. Parallelize uploads

