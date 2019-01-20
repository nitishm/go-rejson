/*
Go-ReJSON is a Go client for ReJSON redis module (https://github.com/RedisLabsModules/rejson)

	ReJSON is a Redis module that aims to provide full support for ECMA-404
	The JSON Data Interchange Standard as a native data type.

	It allows storing, updating and fetching JSON values from Redis keys (documents).

	Primary features of ReJSON Module:
		* Full support of the JSON standard
		* JSONPath-like syntax for selecting element inside documents
		* Documents are stored as binary data in a tree structure, allowing fast access
		  to sub-elements
		* Typed atomic operations for all JSON values types

Go-ReJSON implements all the features of ReJSON Module, without any dependency on the client used for Redis in GoLang.


Installation


To install and use ReJSON module, one must have the pre-requisites installed and setup. Follow the following steps :
	go get github.com/nitishm/go-rejson
	cd $GOPATH/src/github.com/nitishm/go-rejson
	./install-redis-rejson.sh


Examples



*/
package rejson
