version=0.01-$(date +%s)

if [ "$1" == "clean" ] ; then
    rm -rf */out/
fi

echo -----
echo building app
echo -----

cd app
mkdir out
mkdir out/$version
env GOOS=linux go build  -o out/$version/bugle-${version}  -ldflags "-X main.version=$version"
if [ "$?" != "0" ] ;  then
  echo "Failed to build"
  cd ..
  exit
fi
cd ..

echo -----
echo building container
echo -----

cd container
rm -rf bug*
cp ../app/out/$version/* .
docker build -t cminion/bugle:$version .
docker tag  cminion/bugle:$version  cminion/bugle:latest
docker tag  cminion/bugle:$version  cminion/bugle:latest
cd ..
