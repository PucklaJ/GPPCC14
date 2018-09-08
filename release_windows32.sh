echo Compiling ...
go build -a -ldflags '-s' -o GPPCC14_x86.exe
echo Compiling done!
echo Packing files of GPPCC14 ...
zip GPPCC14_x86.zip GPPCC14_x86.exe assets/fonts/* assets/maps/* assets/textures/*
echo Created GPPCC14_x86.zip
echo Done!