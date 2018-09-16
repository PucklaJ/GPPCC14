echo Compiling ...
go build -a -ldflags '-s' -o GPPCC14_x86.exe
echo Compiling done!
echo Packing files of GPPCC14 ...
cp /mingw32/bin/libopenal-1.dll libopenal-1.dll
zip GPPCC14_x86.zip GPPCC14_x86.exe assets/fonts/* assets/maps/* assets/textures/* assets/sounds/* libopenal-1.dll
rm libopenal-1.dll
echo Created GPPCC14_x86.zip
echo Done!
