./goal network create -r $PWD/private -n private -t ../genesis.json
./goal network start -r private/
./goal network status -r private/
./goal wallet list -d private/
./goal wallet new blah -d net1/Primary/
./goal wallet -f blah -d net1/Primary
./goal account new -d net1/Primary/
