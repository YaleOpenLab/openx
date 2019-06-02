./goal network create -r $PWD/private -n private -t ../genesis.json
./goal network start -r private/
./goal network status -r private/
./goal wallet list -d private/
./goal wallet new blah -d net1/Primary/
./goal wallet -f blah -d net1/Primary
./goal account new -d net1/Primary/
./goal clerk send -a 10000000 -f CL4XO7BXWQ7ZZRWGZIJPZW5MHEEBOFJWKM3LMOURX2BOPXVTIIIL2DMIUQ -t YXU3MTTKV74UAGED6ROTHVVPEY5646WI3N5FLLQZWFV66AFKVQ5PMMYDZE -d net1/Primary -w unencrypted-default-wallet
