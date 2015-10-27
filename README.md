# stay
Golang storage and query engine implementing bitmap indexing. Focused on providing fast runtime queries for realtime analytics with a simple REST interface for access.

Start the server using `go run main.go`

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
