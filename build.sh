mkdir -p ./builds/x86
mkdir -p ./builds/x64
GOOS=windows GOARCH=386 go build -o ./builds/386/etcdop.exe
GOOS=windows GOARCH=amd64 go build -o ./builds/amd64/etcdop.exe
GOOS=linux GOARCH=386 go build -o ./builds/386/etcdop
upx ./builds/386/etcdop.exe
GOOS=linux GOARCH=amd64 go build -o ./builds/amd64/etcdop
upx ./builds/amd64/etcdop
