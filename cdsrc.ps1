$project_root = split-path $myinvocation.mycommand.path
set-location "$project_root/src/github.com/heartszhang"

function cd-root{
  set-location $project_root
}