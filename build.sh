# If you wanna using other randomx fork, change the branch
# branches avaliable: master(=random-x) random-xl random-wow random-arq
echo "Target $*"
if [ ! -d "RandomX" ]; then
  git clone https://github.com/maoxs2/RandomX RandomX
fi

cd RandomX
git checkout $*
git pull origin $*
mkdir build
cd build
cmake -G "Unix Makefiles" ..
make -j`nproc`
mv librandomx.a ../../lib
cd ..
rm -rf build
cd ..
