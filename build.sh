version=1.0.1-$(date +%s)

if [ "$1" == "clean" ] ; then
    rm -rf */out/
fi

echo -----
echo building app
echo -----

cd app
mkdir out
mkdir out/$version
env GOOS=linux go build  -o out/$version/LogEar-${version}  -ldflags "-X main.version=$version"
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
rm -rf LogEar*
cp ../app/out/$version/* .
docker build -t cminion/logear:$version .
cd ..

if [ "$1" == "release" ] ; then
  docker tag  cminion/logear:$version  cminion/logear:latest
  docker push cminion/logear:$version
  docker push cminion/logear:latest
else
  docker tag  cminion/logear:$version  cminion/logear:alpha
  docker push cminion/logear:alpha
fi


echo -----
echo deploy chart
echo -----
cd helm
helm delete rii --purge || true && helm install . --name rii
sleep 10
echo -----
echo test
echo -----
curl -I http://cminion.cf/LogEar/ -vk
