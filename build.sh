VERSION=1.0.0

build () {
    GOOS=$1 GOARCH=$2
    go build -ldflags "-X aws-sagemaker-edge-quick-device-setup/distinfo.OS=$1 -X aws-sagemaker-edge-quick-device-setup/distinfo.ARCH=$2 -X aws-sagemaker-edge-quick-device-setup/distinfo.VERSION=$VERSION" -o ./bin/aws-sagemaker-edge-quick-device-setup-$1-$2 main.go
}

if [ $1 != "linux" -a $1 != "windows" ]; then
    echo "Invalid Operating System!"
fi

if [ $2 != "amd64" -a $2 != "386" -a $2 != "arm64" ]; then
    echo "Invalid Architecture!"
fi


build $1 $2
for algo in sha1sum sha224sum sha256sum sha384sum sha512sum; do
    ${algo} ./bin/aws-sagemaker-edge-quick-device-setup-$1-$2 > ./bin/aws-sagemaker-edge-quick-device-setup-$1-$2.${algo}
done 