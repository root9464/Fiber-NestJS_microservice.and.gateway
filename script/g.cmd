@echo off

if "%~3"=="" (
    echo Usage: %0 ^<proto_path^> ^<output_file^> ^<filename^>
    exit /b 1
)

set "proto_path=%~1"
set "output_file=%~2"
set "filename=%~3"

if not exist "%output_file%" mkdir "%output_file%"

protoc --proto_path="%proto_path%" --go_out="%output_file%" --go-grpc_out="%output_file%" %proto_path%/%filename%.proto