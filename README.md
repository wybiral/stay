# stay
Golang storage and query engine implementing bitmap indexing.

Start the server using `go run server/main.go`

## Update keys with Python
To update keys use the /update handler. Data is sent as a JSON object in the form of {key: {feature: True/False}}. Setting features to True will add them for that key and False will remove them. You can update multiple keys at once.

Example:
```
import json
import requests
update = {
    'user:1': {'likes:a': True, 'likes:b': True, 'likes:c': True},
    'user:2': {'likes:b': True, 'likes:c': True, 'likes:d': True},
    'user:2': {'likes:c': True, 'likes:d': True, 'likes:e': True},
}
requests.post('http://localhost:8080/update', data=json.dumps(update))
```

## Query with Python
Querys are done from the /query handler.

Return all keys with feature "likes:b":
```
import json
import requests
query = 'likes:b'
r = requests.post('http://localhost:8080/update', data=json.dumps(query))
print r.json()
```

Return all keys with features 'likes:b' and 'likes:c':
```
import json
import requests
query = ['and', 'likes:b', 'likes:c']
r = requests.post('http://localhost:8080/update', data=json.dumps(query))
print r.json()
```

Available query features: and, or, xor, not
