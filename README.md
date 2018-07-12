# Quip: A Simple Microservice Example for go-kit

## Requirements
----

* each day, post a quote randomly selected from a list
* be able to occasionally add new quotes to the list of quotes
* post new quotes from the queue in order, ahead of the existing list
* when a new quote has been posted, move it to the main list for later reuse in the random selection
* ensure new quotes can only be posted to the list by the maintainer


## Approach
----
* long-running service with scheduler to post daily quote
* service includes http listenter to accept new quotes from client
* use strong cryptographic practices for client/server exchange. Hashed and signed payload
* use existing list of quotes stored in AWS SimpleDB

## Dependencies
----

[Go kit](https://gokit.io/)
    
    go get github.com/go-kit/kit

### Other dependences you may find you need to install

* github.com/go-logfmt/logfmt
* github.com/go-stack/stack
* go get golang.org/x/net/context/ctxhttp
