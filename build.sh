VERSION=0.0.4


build () {
    GOOS=$1 GOARCH=$2 go build -ldflags "-X aws-sagemaker-edge-quick-device-setup/distinfo.OS=$1 -X aws-sagemaker-edge-quick-device-setup/distinfo.ARCH=$2 -X aws-sagemaker-edge-quick-device-setup/distinfo.VERSION=$VERSION" -o ./bin/aws-sagemaker-edge-quick-device-setup-$1-$2 main.go
}

if [ $1-$2 != "linux-amd64" -a $1-$2 != "linux-arm64" ]; then
    echo "USAGE: bash build.sh OS ARCH"
    echo "Supported operating system and architecture combinations"
    echo "- linux amd64"
    echo "- linux arm64"
    exit 1
fi


build $1 $2
for algo in sha1sum sha224sum sha256sum sha384sum sha512sum; do
    ${algo} ./bin/aws-sagemaker-edge-quick-device-setup-$1-$2 > ./bin/aws-sagemaker-edge-quick-device-setup-$1-$2.${algo}
done 
