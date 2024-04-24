
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o linux .

@REM
@REM set CGO_ENABLED=0
@REM set GOOS=windows
@REM set GOARCH=amd64
@REM go build -o win .
@REM
@REM @REM
@REM @REM set CGO_ENABLED=0
@REM @REM set GOOS=darwin
@REM @REM set GOARCH=amd64
@REM @REM go build -o darwin .