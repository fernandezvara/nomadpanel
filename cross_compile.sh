#!/bin/bash
mkdir bin/
gox -arch="amd64" -os="linux darwin" -verbose -output="bin/{{.OS}}_{{.Arch}}/{{.Dir}}"
