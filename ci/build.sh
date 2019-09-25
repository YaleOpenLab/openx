#!/bin/bash

# Trigger this fro mthe go script at routine intervals
echo "Building different packages"

if [ "$WGO" == "" ] ; then
  return
fi

if [ "$GOPATH" == "" ] ; then
  GOPATH=$HOME
fi

echo $GOPATH
cd $GOPATH/go/src/github.com/YaleOpenLab/openx
git pull origin master
go get -v ./...
go build -v ./...
go build
cp openx ci/

go get -v github.com/YaleOpenLab/opensolar
cd $GOPATH/go/src/github.com/YaleOpenLab/opensolar
go get -v ./...
go build -v ./...
cd teller ; go build
cp teller $GOPATH/go/src/github.com/YaleOpenLab/openx/ci/
cd .. ; go build
cp opensolar $GOPATH/go/src/github.com/YaleOpenLab/openx/ci/

cd $GOPATH/go/src/github.com/YaleOpenLab/openx/ci/
tar -cvzf openx.gz openx
tar -cvzf opensolar.gz opensolar
tar -cvzf teller.gz teller
