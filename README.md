## diaginfra

## Program for checking the operation services
## Create conf/config.yaml file checking services and run main file

## Directory project
- build react compiling directory
- conf yaml configurate for testing stend
- internal go module
  - chport check port status
  - churl check url status
  - getlog get and parse logs
  - handl parse result shell command 
  - sshcmd run shell command in ssh
  - webapi gin web server
- public react directory
- src react directory


## Build programm
 ```
 npm init
 npm run build
 go mod init mymain
 go mod tidy
 CGO_ENABLED=0 go build -o diaginfra main.go
```

- Directories needed to work - conf and build.

- Run ./diaginfra 
- Open in browser http://localhost:3000

 Sergey Kukrin 21.11.2022