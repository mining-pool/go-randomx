# If you wanna using other randomx fork, change the branch
# branches avaliable: master(=random-x) random-xl random-wow random-arq
echo "Target $*"
git clone --branch $* https://github.com/maoxs2/RandomX RandomX
cd RandomX

cmake -G "Unix Makefiles" .
make
mv librandomx.a ../lib
cd ..
rm -rf RandomX
