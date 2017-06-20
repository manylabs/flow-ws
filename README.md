# Flow Websocket

Flow websocket sender/reciever for messages

Currently in development

To set up:

extract project to $GOPATH/src/flow-ws. Ex: $GOPATH/src/flow-ws/main.go
```
cd $GOPATH/src/flow-ws
go get
```

for local development, I recommend rerun to autobuild when changes are made.  https://github.com/skelterjohn/rerun

Once rerun is installed
```
rerun flow-ws
```

Alternatively: use full path to rerun if $GOPATH/bin is not set in environment variables

```
/path/to/rerun.exe flow-ws
```


Configuration file path must be defined. FLOWCONFIG environment variable must be set to config.json. Ex: 

```
export FLOWCONFIG=$GOPATH/config.json
```

config.json contains:

{
  "vars": {
    "gorpDatabaseType":  "sqlite3",
    "gorpDatabaseUri":  "E:/projects/rhizo-server/main/rhizo.db",
    "debugMode":  "false",
    "listenPort": ":8000"
  }
}

// gorpDatabaseType = "mysql", "sqllite3", "postgresql".  This is the driver type
// gorpDatabaseUri = authenticated string to connect to database, or path to database if run on a socket connection.  This is specific to the DB driver used.  See https://github.com/go-sql-driver/mysql#examples for mysql examples of TCP and socket connections.

Included "database/sql" drivers are "mysql", "sqllite3", "postgresql".
