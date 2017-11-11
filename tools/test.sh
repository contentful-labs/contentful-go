#!/usr/bin/env bash
set -e

rm coverage.txt || true
touch coverage.txt

for d in $(go list ./... | grep -v /vendor/); do
    go test -v -coverprofile=profile.out -covermode=count $d

    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done

# to make `go tool cover -html=coverage.txt` happy
# remove the lines starting with mode
# remove the empty lines
sed -i'' -e '/^\s*$/d' coverage.txt
echo "$(awk '!/^mode:/ || !f++' coverage.txt)" > coverage.txt
