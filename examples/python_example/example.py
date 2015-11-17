'''
Uses the requests module to interact with a local StayDB server.
'''

def stay(type, data, host='localhost', port=8080):
    import json
    import requests
    url = 'http://%s:%i/%s' % (host, port, type)
    r = requests.post(url, data=json.dumps(data))
    return r.json()

stay('add', {
    'nina':   ['species:human', 'sex:female'],
    'elaine': ['species:cat',   'sex:female'],
    'davy':   ['species:human', 'sex:male'],
    'percy':  ['species:cat',   'sex:male'],
})

print stay('query',
    ['or',
        ['and', 'sex:male',   'species:human'],
        ['and', 'sex:female', 'species:cat'],
    ]
)
