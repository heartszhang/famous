$project_root = split-path $MyInvocation.MyCommand.Path
$backend = "$project_root\src\github.com\heartszhang\backend\run\run.exe"
$backend_script="$project_root\backend.bat"
$backend_stop="$project_root\backend_stop.bat"
$mongod = "$project_root\mongodb-win32-x86_64-2008plus-2.4.6-rc1\bin\mongod.exe"
$mongo="$project_root\mongodb-win32-x86_64-2008plus-2.4.6-rc1\bin\mongo.exe"

$dist_dir="$project_root\dist"
$dist_dir_debug="$project_root\famousfront\bin\Debug"

#mkdir.exe -v -p $dist_dir_debug
#mkdir.exe -p $dist_dir_debug\mongodb
#mkdir.exe -p $dist_dir_debug\mongodb\db
New-Item -path $dist_dir_debug\mongodb\db -type directory -ErrorAction SilentlyContinue
cp.exe -u -v $backend "$dist_dir_debug\run.exe"
cp.exe -u -v $mongod  "$dist_dir_debug\mongod.exe"
cp.exe -u -v $backend_script "$dist_dir_debug\backend.bat"
cp.exe -u -v $mongo "$dist_dir_debug\mongo.exe"
cp.exe -u -v "$project_root\mongo_stop.js" "$dist_dir_debug\mongo_stop.js"
cp.exe -u -v $backend_stop "$dist_dir_debug\backend_stop.bat"