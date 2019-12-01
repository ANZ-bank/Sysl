@ECHO OFF
SETLOCAL EnableDelayedExpansion

set ROOT=cmd\tests
for /R %ROOT% %%f in (*.sysl) do (
	dist\gosysl.exe --log debug pb --mode textpb --root %ROOT% -o %ROOT%\%%~nf.win.txt  /%%~nf.sysl || exit /b !errorlevel!
)

del cmd\tests\*.win.txt
