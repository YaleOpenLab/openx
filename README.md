# OpenX

[![Build Status](https://travis-ci.com/YaleOpenLab/openx.svg?branch=master)](https://travis-ci.com/YaleOpenLab/openx)
[![Codecov](https://codecov.io/gh/YaleOpenLab/openx/branch/master/graph/badge.svg)](https://codecov.io/gh/YaleOpenLab/openx)
[![Go Report Card](https://goreportcard.com/badge/github.com/YaleOpenLab/openx)](https://goreportcard.com/report/github.com/YaleOpenLab/openx)

This repo contains a WIP implementation of the OpenX idea in stellar. Broadly, the openx model seeks to implement the paradigm of investing and developing projects without hassles and enabling smart ownership with the help of semi trust-less entities on the blockchain. The openx model can be thought more generally as a platform of platforms and houses multiple platforms within it (in `platforms/`).  The goal is to have a common interface between all parties that relate to a project: investors, investees (i.e. beneficiaries or receivers of the investment) and the family of developers (that include all service providers). Depending upon the country of project origin, investors may be required to complete KYC to be able to invest in assets. The openx model can be adapted to any blockchain but we currently use the Stellar blockchain to create and move project assets.  

## Openx Repositories

Openx, like Go is built on the idea of modularity and reusability of packages.

1. This repo contains the bare bone architecture necessary for the platform  
2. [essentials](https://github.com/Varunram/essentials) contains the code necessary for commonly used packages, crypto and database handlers  
3. [openx-cli](https://github.com/Varunram/openx-cli) contains CLI clients that can interface with openx  
4. [opensolar](https://github.com/YaleOpenLab/opensolar) contains an implementation of the openx idea targeted at solar infrastructure  

While making an issue related to openx however, please open in this repo for easy tracking and referencing in multiple places.

## Documentation

Comprehensive documentation on each platform is available inside each subrepo.

1. [Opensolar](platforms/opensolar/README.md)
2. [Opzones](platforms/ozones/README.md)

## Getting Started

In order to be able to run this, you need to have the latest version of go installed. [Here](https://tecadmin.net/install-go-on-ubuntu/) is a quick tutorial on how to get go installed on a Linux / macOS machine. Older versions of go (more than two versions old according to [the golang wiki](https://github.com/golang/go/wiki/MinorReleases)) have unpatched vulnerabilities and as a result, we will not be backporting versions.

Once you have go installed, you need to get the packages in this repo. Before that, you might need to install the `stellar/go` package separately since it uses a separate dependency manager. Get the `stellar/go` package separately and then run `dep ensure -v` inside `$HOME/stellar/go`. This might take a few minutes to complete.

Once you're done with `stellar/go`, clone the repo using `git clone https://github.com/YaleOpenLab/openx.git` and install the other dependencies using `go get -v ./...`

Now you should be ready to compile (`go build`) and run (`./openx`) the openx executable.

## Installing IPFS

ipfs is used by some parts of the program to store legal contracts, files that the user might want to store permanently. Download a release from https://github.com/ipfs/go-ipfs/releases and run install.sh. In case you face an issue with migration between various ipfs versions, you might need to run [fs-repo-migrations](https://github.com/ipfs/fs-repo-migrations/blob/master/run.md) to migrate to a newer version. If you don't have anything valuable, you can delete the directory and run `ipfs init` again (this will delete the data stored in ipfs prior to deletion)

You need to keep your peer key (`ipfs.key` usually) in a safe place for future reference. Start ipfs using `ipfs daemon` and you can test if it worked by creating a test file `test.txt` and run `ipfs add test.txt` to see if it succeeds. The resultant hash can be decrypted using `curl "http://127.0.0.1:8080/ipfs/hash"` where 8080 is the endpoint of the ipfs server or by doing `cat /ipfs/hash` directly. You can also refer to [this helpful tutorial](https://michalzalecki.com/set-up-ipfs-node-on-the-server/) in order to get started with ipfs.

## Running tests

Running tests is simple with `go test` but the tests have flags since some require running other daemons in the background (`ipfs`). There are two kinds of flags right now - `travis` and `all`. If you need the coverage stats as well, you need to
```
go get golang.org/x/tools/cmd/cover
```
if you already don't have the package. Running `go test --tags="travis" -coverprofile=test.txt ./...` should run all the tests and provide coverage data on each specific package. Running with the tag `travis` will omit the tests in `ipfs/` which requires [a local `go-ipfs` node running](https://michalzalecki.com/set-up-ipfs-node-on-the-server/) as described above.

## Contributing

Please feel free to open Pull Requests and Issues with your changes and suggestions. Before working on a major feature, please describe the same in an issue so everyone can understand what you're building.

## Security

<img src="security/discloseio.png" width="50">  

For security related issues, *DO NOT OPEN A GITHUB ISSUE!*. Please disclose the information responsibly by sending a (preferably PGP Encrypted) email to `contact@varunram.com`. [Our PGP Key](https://pgp.mit.edu/pks/lookup?op=vindex&fingerprint=on&search=0x708C606504A49970) fingerprint is `C98F 0014 9A99 36E4 E56D  2471 708C 6065 04A4 9970`  

In addition to this, openx is fully fully compliant with the [disclose.io](https://disclose.io) core terms followed by [bugcrowd](https://www.bugcrowd.com/resource/what-is-responsible-disclosure/). For more info, see [SECURITY.md](SECURITY.md)

## License
[GPL3](https://github.com/YaleOpenLab/openx/blob/master/LICENSE)
