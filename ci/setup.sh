echo $GOPATH
if [ "$GOPATH" == "" ] ; then
  WGO="$(which go)"
  if [ "$WGO" == "" ] ; then
    echo "installing golang on your machine"
    sudo apt-get -y upgrade
    wget https://dl.google.com/go/go1.12.4.linux-amd64.tar.gz
    sudo tar -xvf go1.12.4.linux-amd64.tar.gz
    sudo mv go /usr/local
    echo 'GOROOT=/usr/local/go' >> ~/.profile
    echo 'GOPATH=$HOME' >> ~/.profile
    echo 'PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.profile
    source ~/.profile
  fi
fi

# At this instance, we assume golang is installed on the user's machine
go get -v github.com/YaleOpenLab/openx
cd $GOPATH/go/src/github.com/YaleOpenLab/openx
go get -v ./...
go build -v ./...
go build
cp openx ~/

go get -v github.com/YaleOpenLab/opensolar
cd $GOPATH/go/src/github.com/YaleOpenLab/opensolar
go get -v ./...
go build -v ./...
cp teller/teller ~/
go build
cp opensolar ~/

tar -cvzf openx.gz openx
tar -cvzf opensolar.gz opensolar
tar -cvzf teller.gz teller
