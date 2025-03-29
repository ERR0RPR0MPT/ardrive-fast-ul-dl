@echo off
set GOOS=linux
set GOARCH=amd64
go build -o .\build\ardrive_stream_linux_amd64 .\cmd\ardrive-stream
set GOOS=linux
set GOARCH=arm64
go build -o .\build\ardrive_stream_linux_arm64 .\cmd\ardrive-stream
set GOOS=windows
set GOARCH=amd64
go build -o .\build\ardrive_stream_windows_amd64.exe .\cmd\ardrive-stream
pause