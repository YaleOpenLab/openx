echo $GOPATH
if [ "$GOPATH" == "" ] ; then
  WGO="$(which go)"
  if [ "$WGO" == "" ] ; then
    if [ "$(uname)" == "Darwin" ]; then
      WBREW="$(which brew)"
      if [ "$WBREW" == "" ] ; then
        # install brew
        echo "installing brew on your mac"
        /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
      fi
      echo "installing golang on your mac"
      brew install golang
    elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
      echo "installing golang on your linux machine"
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
fi

if [ "$GOPATH" == "" ] ; then
  GOPATH=$HOME
fi

sh build.sh
