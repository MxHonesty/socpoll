Start-Process -FilePath cmd.exe -ArgumentList "/c title Nsqdlookup Server & nsqlookupd"
Start-Process -FilePath cmd.exe -ArgumentList "/c title Nsqd Client & nsqd --lookupd-tcp-address=localhost:4160"
Start-Process -FilePath cmd.exe -ArgumentList "/c title Mongod & mongod --dbpath ./db"
