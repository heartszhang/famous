$project_root = split-path $MyInvocation.MyCommand.Path

$env:GOPATH = "$env:GOPATH;$project_root"
$env:PATH = "$env:PATH;$project_root/bin"
# mongod --config /etc/mongodb.conf
$mongo_conf = (Get-ChildItem mongod.conf).FullName

## mongo.config
#dbpath=e:\mongodb
#logpath=e:\mongodb\mongo.log
#diaglog=3
#bind_ip=127.0.0.1
##

#mongod -v --bind_ip 127.0.01 --logpath xxxxx --dbpath xxxxx
$mongo_ip = "127.0.0.1"
$mongo_logpath = "$project_root\mongodb\mongo.log"
$mongo_dbpath = "$project_root\mongodb\db\"
#"--bing_ip $mongo_ip -- logpath $mongo_logpath --dbpath $mongo_dbpath"

function db_start(){
    Start-Process "mongod.exe" "--bind_ip $mongo_ip --logpath $mongo_logpath --dbpath $mongo_dbpath"
}
function db_stop(){
    Start-Process "mongo.exe" "$mongo_ip/admin mongo_stop.js"
}

# https://launchpad.net/bzr/2.6/2.6b1/+download/bzr-2.6b1-1-setup.exe
#hg clone https://code.google.com/p/go.net/ src/code.google.com/p/go.net
#go get github.com/robfig/revel/revel
#go get labix.org/v2/mgo
#wget http://downloads.mongodb.org/win32/mongodb-win32-x86_64-2008plus-2.4.6.zip
# go get -u github.com/djimenez/iconv-go
# configure cgo 
#CGO_CFLAGS=-IC:\MINGW64/include
#CGO_LDFLAGS=-LC:\MINGW64/lib
#download libiconv for windows 64bits
#go get -u github.com/qiniu/iconv
# // #cgo windows LDFLAGS: -liconv // add to iconv.go

$OutputEncoding = New-Object -typename System.Text.UTF8Encoding