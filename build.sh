VERSION="1.2"
function buildMac() {
    go build -o bin/$VERSION/picnic-tui-mac .
    echo "Mac built"
}

function buildWin64() {
    GOOS=windows GOARCH=amd64 go build -o bin/$VERSION/picnic-tui-win64.exe .
    echo "Win64 built"
}

function buildWin32() {
    GOOS=windows GOARCH=386 go build -o bin/$VERSION/picnic-tui-win32.exe .
    echo "Win32 built"
}

function buildLinux64() {
    GOOS=linux GOARCH=amd64 go build -o bin/$VERSION/picnic-tui-x64 .
    echo "Linux 64 built"
}

function buildLinux32() {
    GOOS=linux GOARCH=386 go build -o bin/$VERSION/picnic-tui-x32 .
    echo "Linux 32 built"
}


buildMac
buildLinux32
buildLinux64
buildWin32
buildWin64