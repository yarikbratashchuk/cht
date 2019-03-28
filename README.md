## cht

### 👋 Hi there!

**cht** is a cli chat for small teams.

### Demo

<p align="center">
<img width="400" src='https://user-images.githubusercontent.com/12980380/55159801-87e7a800-516a-11e9-8a82-ea32f75d96b4.gif' alt="cht usage demo" />
</p>

### Installation

Now you can only build from source:

```bash
$ make install
```

### Usage

#### Step 1. Run `cht-server`

```bash
$ cht-server -h
Usage:
  cht-server [OPTIONS]

Application Options:
  -p, --port=     port to listen on (default: 9090)
      --loglevel= log level for all subsystems {trace, debug, info, error, critical} (default: info)

Help Options:
  -h, --help      Show this help message

$ cht-server # starting chat server
```

#### Step 2. Use `cht` client (in separate session)

```bash
$ cht -h 

NAME:
   cht - chat cli

USAGE:
   cht [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     start    join the chat room and start sending and receiving messages
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --room value      room to join (default: "random")
   --server value    server to connect (default: "127.0.0.1:9090")
   --nickname value  your nickname (default: "noname")
   --help, -h        show help
   --version, -v     print the version

$ cht --nickname="yarik" start # join "random" chat room and start chating
connecting to room random...
Jarvis: papa can you hear me
Caoimhe: hop hey lalaley
yarik connected
  <- write your message here
```


That's it. Enjoy! ❤️


----


MIT License, 2019. Yarik Bratashchuk

