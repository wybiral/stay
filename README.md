# StayDB
StayDB is an in-memory storage engine written in [Go](https://golang.org/) implementing [bitmap indexing](https://en.wikipedia.org/wiki/Bitmap_index) and querying over a RESTful API. Focused on supporting fast real time analytics of large sets of data.

## Installing StayDB

StayDB is written in Go, so make sure you have that installed first [(see here)](https://golang.org/doc/install). Once Go is installed use Go's package manager to install StayDB by typing the following in your command line:

`go get github.com/wybiral/stay`

Then you can build StayDB by typing:

`go build github.com/wybiral/stay`

This should produce an executable in your current directory. You can run this executable to start StayDB (running it with the `-help` flag will explain the command line options).

## Update keys with Python
You can add features to a key using the /add handler by supplying a JSON mapping of {key: [...features to add...]}

Example:
```
import json
import requests
update = {
    'user:1': ['likes:a', 'likes:b', 'likes:c'],
    'user:2': ['likes:b', 'likes:c', 'likes:d'],
    'user:3': ['likes:c', 'likes:d', 'likes:e'],
}
requests.post('http://localhost:8080/add', data=json.dumps(update))
```

## Query with Python
Querys are done from the /query handler.

Return all keys with feature "likes:b":
```
import json
import requests
query = 'likes:b'
r = requests.post('http://localhost:8080/query', data=json.dumps(query))
print r.json()
```

Return all keys with features 'likes:b' and 'likes:c':
```
import json
import requests
query = ['and', 'likes:b', 'likes:c']
r = requests.post('http://localhost:8080/query', data=json.dumps(query))
print r.json()
```

Available query features: and, or, xor, not
