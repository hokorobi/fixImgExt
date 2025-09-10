go vet
go vet -vettool=C:\Users\syuhe\bin\shadow.exe
go vet -vettool=C:\Users\syuhe\bin\defers.exe
staticcheck .
go build -ldflags -s