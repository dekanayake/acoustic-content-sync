rm -f -R buildScript/build
mkdir  buildScript/build
go build -tags standard
cp acoustic-content-sync buildScript/build
cp -r script/.env buildScript/build
cp -r buildScript/configs/*.* buildScript/build
chmod 777 buildScript/build/acoustic-content-sync
chmod 777 buildScript/build/.env
cd buildScript/build
zip  acoustic-content-sync_darvin_386.zip  .env *.yaml acoustic-content-sync