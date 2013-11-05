@echo off
set project_root=%~dp0
set mongo_ip=127.0.0.1
set mongo_logpath=%project_root%\mongodb\mongo.log
set mongo_dbpath=%project_root%\mongodb\db\

call mongod.exe --bind_ip %mongo_ip% --logpath %mongo_logpath% --dbpath %mongo_dbpath%
call run.exe
