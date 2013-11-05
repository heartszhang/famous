@echo off
set mongo_ip=127.0.0.1
call mongo.exe %mongo_ip%/admin mongo_stop.js
