VERSION=1.0.0

build () {
    GOOS=$1 GOARCH=$2
    go build -ldflags "-X quick-device-setup/distinfo.OS=$1 -X quick-device-setup/distinfo.ARCH=$2 -X quick-device-setup/distinfo.VERSION=$VERSION" -o ./bin/quick-device-setup-$1-$2 main.go
}

case "$1-$2" in
    "linux-amd64")
        build "linux" "amd64"
        ;;
    "linux-amd")
        build "linux" "amd"
        ;;
    "linux-armv8")
        build "linux" "armv8"
        ;;
    *)
        echo "Invalid Command!"
        exit 1
        ;;
esac
