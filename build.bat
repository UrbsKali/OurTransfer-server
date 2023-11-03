@echo off
cmd /C "set GOOS=windows&set GOARCH=amd64&go build -o ./build/file-win-amd64.exe urbskali/file"
cmd /C "set GOOS=windows&set GOARCH=arm64&go build -o ./build/file-win-arm64.exe urbskali/file"
cmd /C "set GOOS=linux&set GOARCH=amd64&go build -o ./build/file-deb-amd64 urbskali/file"
cmd /C "set GOOS=linux&set GOARCH=arm64&go build -o ./build/file-deb-arm64 urbskali/file"
cmd /C "set GOOS=darwin&set GOARCH=amd64&go build -o ./build/file-mac-amd64 urbskali/file"
cmd /C "set GOOS=darwin&set GOARCH=arm64&go build -o ./build/file-mac-arm64 urbskali/file"
echo Build complete.