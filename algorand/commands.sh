./goal network create -r $PWD/private -n private -t ../genesis.json
./goal network start -r private/
./goal network status -r private/
./goal wallet list -d private/
