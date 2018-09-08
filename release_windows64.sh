echo Compiling ...
go build -a -ldflags '-s' -o GPPCC14_x64.exe
echo Compiling done!
echo Packing files of GPPCC14 ...
zip GPPCC14_x64.zip GPPCC14_x64.exe assets/fonts/* assets/maps/* assets/textures/* /mingw64/bin/libopenal-1.dll
echo Created GPPCC14_x64.zip
echo Done!