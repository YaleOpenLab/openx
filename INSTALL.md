# Installing openx

Openx is a framework on top of which investment platforms that want to use blockchain investments and smart contract validation can be built on. An example of a platform building on top of openx is the [opensolar platform](https://github.com/YaleOpenLab/openx).

One can create their own instance of openx by using the [create-openx-app](https://github.com/YaleOpenLab/create-openx-app) utility. One can also refer to opensolar to see how things are constructed and build their own instance by interfacing with openx's extensive APIs.

Openx right now supports only Stellar blockchain investments and encrypted key storage. After the pilot phase, additional blockchain support will be explored on. The platform has no access to funds at any time during the smart contract execution and if the platform goes rogue (not signing towards the recipient withdrawing funds), one can use the rescue utility to rescue their funds (the platform cannot unilaterally steal your funds in any scenario).

### Operating system

Openx has been tested on Ubuntu 16.04 LTS and macOS 10.13+. The build status on other operating systems is unknown. There are no plans to test on alternate operating systems at the moment.

### Prerequisites

1. Golang

To Save you from a few clicks, use the following script to download and install golang.

Linux:
```
sudo apt-get -y upgrade
wget https://dl.google.com/go/go1.12.4.linux-amd64.tar.gz
sudo tar -xvf go1.12.4.linux-amd64.tar.gz
sudo mv go /usr/local
echo 'GOROOT=/usr/local/go' >> ~/.profile
echo 'GOPATH=$HOME' >> ~/.profile
echo 'PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.profile
source ~/.profile
which go
```

MacOS:
```
brew install golang
```
If you don't have brew installed, its highly recommended:
```
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

2. IPFS

To install IPFS, please follow this excellent installation guide from the IPFS team: [guide](https://docs.ipfs.io/guides/guides/install/)

### Building from Source

Assuming that you have GOPATH set from installing go in step 1, please run this following script:

```
cd $GOPATH/src/github.com/YaleOpenLab/openx
go get -v ./...
go build
```
This gets the necessary dependencies for openx and builds the openx executable.

### Setting up the config file

Duplicate [dummyconfig.yaml](dummyconfig.yaml), rename to config.yaml and replace the relevant values with those desired.

### Downloading a prebuilt version

[The builds website](https://builds.openx.solar/fe) has daily builds for opensolar, openx and the teller. Running them should be as simple as running the executable.
