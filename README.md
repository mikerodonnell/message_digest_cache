## message_digest_cache
demo REST service for key-value-pair storage, optimized for immutable values

### about
stores message strings and returns a SHA 256 hash for future retrieval. uses a combination of local in-memory and distributed Redis caches. since the cache key is a always the SHA hash of the value, there's no updates, expiration, or cache synchronization to worry about across app nodes. so each app node can safely return the locally cached value if one exists, and pull from Redis if not. the local cache won't help in all cases, but never hurts. 

### quickstart
```
git clone git@github.com:mikerodonnell/message_digest_cache

cd message_digest_cache

docker-compose up
```

### usage
#### store a message
```
curl --header "Content-Type: application/json" --data '{"Message": "Be Sure To Drink Your Ovaltine"}' "localhost:8000/messages"

{"digest":"64ec6238b3a34d00e3ce27dd4a0164d80d381479b17fbebc8da382576faa5a06"}
```

#### retrieve the message
```
curl --header "Content-Type: application/json" "localhost:8000/messages/64ec6238b3a34d00e3ce27dd4a0164d80d381479b17fbebc8da382576faa5a06"

{"message":"Be Sure To Drink Your Ovaltine"}
```

### tests
```
docker-compose run app go test -v ./...
```

### notes
runs on Go 1.11 with Go Modules. no GOPATH needed!