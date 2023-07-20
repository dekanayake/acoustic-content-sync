dot_env_file_name=$(if [ $1 = "nonprod" ]; then echo ".env_nonprod"; else echo ".env_prod"; fi)
echo "will be used ${dot_env_file_name} env file"
rm -f -R buildScript/build
mkdir  buildScript/build
go build -tags standard
cp acoustic-content-sync buildScript/build
cp -r script/${dot_env_file_name} buildScript/build/.env
cp -r buildScript/configs/*.* buildScript/build
chmod 777 buildScript/build/acoustic-content-sync
chmod 777 buildScript/build/.env
cd buildScript/build
zip  acoustic-content-sync_darvin_386.zip  .env *.yaml acoustic-content-sync