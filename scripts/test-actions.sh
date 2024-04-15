#!/bin/bash

curl -X POST "http://localhost:8090/apis/actions" \
    --url-query "name=demo" \
    --url-query "namespace=demo-system" \
    --url-query "group=composition.krateo.io" \
    --url-query "version=v1alpha1" \
    --url-query "plural=fireworksapps" \
    --url-query "kind=Fireworksapp" \
    --url-query "verbose=true" \
    --url-query "sub=local" \
    --url-query "orgs=devs" \
    -H "Content-Type: application/json" \
    --data-binary @./testdata/fireworksapp.fake.json
