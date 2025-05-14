@echo off
setlocal enabledelayedexpansion
echo Checking processes on ports 9001 to 9004...

:: Looping through all the ports of the service which has been running.
for /L %%p in (9001,1,9004) do (
    echo Checking port :%%p...

    :: Initially the Process ID is None.
    set PID=

    for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":%%p"') do (

        :: If the process found we kill it.
        set PID=%%a
        echo Process found on port :%%p with PID: !PID!

        echo Killing process with PID: !PID!
        taskkill /PID !PID! /F
        echo Process with PID !PID! killed.
    )

    :: If not found just prints out a message.
    if not defined PID (
        echo No process found on port :%%p.
    )
)

echo Done checking and killing processes.
pause
