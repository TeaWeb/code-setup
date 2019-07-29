env GOOS=windows GOARCH=amd64 go build -o coding-windows-x64.exe -ldflags="-s -w"  ../main/main.go
env GOOS=windows GOARCH=386 go build -o coding-windows-386.exe -ldflags="-s -w"  ../main/main.go
env GOOS=linux GOARCH=amd64 go build -o coding-linux -ldflags="-s -w"  ../main/main.go
env GOOS=darwin GOARCH=amd64 go build -o coding-darwin -ldflags="-s -w"  ../main/main.go