
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o linux .

@REM set CGO_ENABLED=0
@REM set GOOS=windows
@REM set GOARCH=amd64
@REM go build -o win .
@REM
@REM set CGO_ENABLED=0
@REM set GOOS=darwin
@REM set GOARCH=amd64
@REM go build -o darwin .