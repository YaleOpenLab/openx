#!/bin/bash

# Trigger this fro mthe go script at routine intervals
echo "Building different packages"

WGO="$(which go)"
if [ "$WGO" == "" ] ; then
  exit
fi

if [ "$GOPATH" == "" ] ; then
  GOPATH=$HOME
fi

echo $GOPATH
cd $GOPATH/go/src/github.com/YaleOpenLab/openx
git pull origin master
go get -v ./...

env GOOS=darwin GOARCH=amd64 go build -o openx-darwinamd64
env GOOS=linux GOARCH=amd64 go build -o openx-linuxamd64
env GOOS=linux GOARCH=386 go build -o openx-linux386
env GOOS=linux GOARCH=arm64 go build -o openx-arm64
env GOOS=linux GOARCH=arm go build -o openx-arm
mv openx-linuxamd64 openx-linux386 openx-arm64 openx-arm openx-darwinamd64 ci/

go get -v github.com/YaleOpenLab/opensolar
cd $GOPATH/go/src/github.com/YaleOpenLab/opensolar
go get -v ./...

cd teller
env GOOS=darwin GOARCH=amd64 go build -o teller-darwinamd64
env GOOS=linux GOARCH=amd64 go build -o teller-linuxamd64
env GOOS=linux GOARCH=386 go build -o teller-linux386
env GOOS=linux GOARCH=arm64 go build -o teller-arm64
env GOOS=linux GOARCH=arm go build -o teller-arm
mv teller-* $GOPATH/go/src/github.com/YaleOpenLab/openx/ci/

cd ..

env GOOS=darwin GOARCH=amd64 go build -o opensolar-darwinamd64
env GOOS=linux GOARCH=amd64 go build -o opensolar-linuxamd64
env GOOS=linux GOARCH=386 go build -o opensolar-linux386
env GOOS=linux GOARCH=arm64 go build -o opensolar-arm64
env GOOS=linux GOARCH=arm go build -o opensolar-arm
mv opensolar-* $GOPATH/go/src/github.com/YaleOpenLab/openx/ci/

cd $GOPATH/go/src/github.com/YaleOpenLab/openx/ci/

tar -cvzf openx-darwinamd64.gz openx-darwinamd64
tar -cvzf openx-linuxamd64.gz openx-linuxamd64
tar -cvzf openx-linux386.gz openx-linux386
tar -cvzf openx-arm64.gz openx-arm64
tar -cvzf openx-arm.gz openx-arm

tar -cvzf opensolar-darwinamd64.gz opensolar-darwinamd64
tar -cvzf opensolar-linuxamd64.gz opensolar-linuxamd64
tar -cvzf opensolar-linux386.gz opensolar-linux386
tar -cvzf opensolar-arm64.gz opensolar-arm64
tar -cvzf opensolar-arm.gz opensolar-arm

tar -cvzf teller-darwinamd64.gz teller-darwinamd64
tar -cvzf teller-linuxamd64.gz teller-linuxamd64
tar -cvzf teller-linux386.gz teller-linux386
tar -cvzf teller-arm64.gz teller-arm64
tar -cvzf teller-arm.gz teller-arm

cd $GOPATH/go/src/github.com/YaleOpenLab/
cp -r openx openx-temp
cd openx-temp
rm -rf .git/ ci/
cd ..
tar -cvzf openx.gz openx-temp/
mv openx.gz $GOPATH/go/src/github.com/YaleOpenLab/openx/ci

cd $GOPATH/go/src/github.com/YaleOpenLab/
cp -r opensolar opensolar-temp
cd opensolar-temp
rm -rf .git/
cd ..
tar -cvzf opensolar.gz opensolar-temp/
mv opensolar.gz $GOPATH/go/src/github.com/YaleOpenLab/openx/ci

cd $GOPATH/go/src/github.com/YaleOpenLab/opensolar/
cp -r teller teller-temp
tar -cvzf teller.gz teller-temp/
mv teller.gz $GOPATH/go/src/github.com/YaleOpenLab/openx/ci

cd $GOPATH/go/src/github.com/YaleOpenLab/
rm -rf *-temp/
