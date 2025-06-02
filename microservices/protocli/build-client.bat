@echo off
setlocal enabledelayedexpansion

if not exist bin (
    mkdir bin
)

docker build -t protocli-builder .

for /f "delims=" %%i in ('docker create protocli-builder') do set CONTAINER_ID=%%i

docker cp %CONTAINER_ID%:/app/dist/protocli-linux bin\protocli-linux
docker cp %CONTAINER_ID%:/app/dist/protocli.exe bin\protocli.exe

docker rm %CONTAINER_ID%

echo Бинарники скопированы в bin
pause
