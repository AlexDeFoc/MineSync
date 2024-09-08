@echo off


del /Q /S export\ > NUL 2>&1
mkdir export\ > NUL 2>&1


GOTO :MAIN

:: WINDOWS BUILDs
:windows_start
    echo --== WINDOWS
EXIT /B 0

:windows_i386
    set GOOS=windows
    set GOARCH=386
    echo -= windows_i386
    go build -o mineSync.exe -tags noaudio -ldflags "-w -s" main.go
    tar -cf windows_i386.tar -C .\ mineSync.exe > NUL 2>&1
    7z a -t7z -mx=9 windows_i386.tar.xz windows_i386.tar > NUL 2>&1
    move windows_i386.tar.xz export\ > NUL 2>&1
    del /Q *.tar > NUL 2>&1
    del /Q mineSync.exe > NUL 2>&1

EXIT /B 0

:windows_amd64
    set GOOS=windows
    set GOARCH=amd64
    echo -= windows_amd64
    go build -o mineSync.exe -tags noaudio -ldflags "-w -s" main.go
    tar -cf windows_amd64.tar -C .\ mineSync.exe > NUL 2>&1
    7z a -t7z -mx=9 windows_amd64.tar.xz windows_amd64.tar > NUL 2>&1
    move windows_amd64.tar.xz export\ > NUL 2>&1
    del /Q *.tar > NUL 2>&1
    del /Q mineSync.exe > NUL 2>&1

EXIT /B 0


:padding
    echo ---------------------
EXIT /B 0


:MAIN
call :windows_start
call :windows_i386
call :windows_amd64
call :padding
