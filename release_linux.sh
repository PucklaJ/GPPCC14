echo Compiling ...
go build -a -ldflags '-s' -o GPPCC14
echo Compiling done!
echo Packing files of GPPCC14 ...
tar -cf GPPCC14.tar GPPCC14 assets/fonts/* assets/maps/* assets/textures/*
echo Created GPPCC14.tar
echo Done!