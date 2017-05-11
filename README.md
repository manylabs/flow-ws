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

Alternatively: if development on a Windows environment

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
    "SQLLitePath":  "/path/to/rhizo.db",
    "debugMode":  "false",
    "listenPort": ":8000"
  }
}