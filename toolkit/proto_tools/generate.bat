
del ..\..\protocol\* /f /s /q /a
del ..\..\src\protocol\* /f /s /q /a

python ./merge.py

cd ..\..\protocol
..\toolkit\proto_tools\protoc.exe --proto_path=./ --go_out=../src/protocol/  ./*.proto

pause
