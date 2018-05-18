@echo off
cd "../"
echo ----------------- To Lua Pb -------------------
set client_path="../lua/pbnet/pb/"

SETLOCAL ENABLEDELAYEDEXPANSION 
for %%i in (*.proto) do (    
set b=%%~ni
echo toLua: !b!.pb
"bat/"protoc --descriptor_set_out=%client_path%!b!.pb !b!.proto
) 
@rem cd "../"
@rem echo ----------------- CopysLuaProto -------------------
@rem copy /y "Proto" "../lua/pbnet/proto"
@pause