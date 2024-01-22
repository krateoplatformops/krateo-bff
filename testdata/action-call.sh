#!/bin/bash

curl -X POST 'http://localhost:8080/apis/action?sub=cyberjoker&v' \
    -H 'Content-Type: application/json' \
    -d '{"name":"view","endpointRef":{"name":"typicode-endpoint","namespace":"demo-system"}, "path":"/todos/1"}'
