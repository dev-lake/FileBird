<h1 align="center">
  <br>
  <img src="https://github.com/dev-lake/FileBird/blob/main/resources/logo.png?raw=true" alt="FileBird" width="400">
  <!-- <br> -->
  <!-- FileBird -->
  <!-- <br> -->
</h1>

<h4 align="center">A Efficient File Transfer Tool Powered by gRPC</h4>
<p align="center">
    <img alt="GitHub go.mod Go version (subdirectory of monorepo)" src="https://img.shields.io/github/go-mod/go-version/dev-lake/FileBird?filename=client%2Fgo.mod&style=for-the-badge">
</p>

## Features
- Add server information
- Show Server information
- Delete Server
- Upload local file to remote
- Download reomte file to local
- Transmit remote file to remote
- Copy local/remote file
- Move local/remote file
- Del local/remote file
- list local/remote file
- make local/remote dir


## How to Build
### Build Server
```bash
cd server/
go mod tidy
go build -o filebird-server
```
### Build client
```bash
cd client/
go mod tidy
go build -o filebird
```

## How to Use
### Add & Show Server
```bash
# add server
./filebird add_server -n SERVERNAME -a IP -p PORT

# show server
./filebird show_server

# delete server
./filebird del_server
```

### Get File Info
```bash
# check dir file
./filebird ls SERVERNAME:

# check current directory
./filebird pwd SERVERNAME:

# change current dir
./filebird cd SERVERNAME:/home

# Get file info
./bin/filebird-mac-arm64 info vm:filebird-server
filebird-server 10.211.55.4 2000

File Info
---------
Name: filebird-server
Size: 12226427
Owner: parallels parallels
ModTime: 2023-08-20 16:01:02.573767628 +0800 CST
IsDir: false
Mode: -rwxr-xr-x
Path: filebird-server
```

### Copy, Move, Delete
```bash
./filebird cp PATH01 PATH02
./filebird mv PATH01 PATH02
./filebird rm PATH
```

### PATH format example:
- local path: `/root/path...`
- remote path: `ServerName:/root/Path...`


## TODO
- ~~Show Progress bar~~ (complete)
- ~~实现文件夹的递归传输~~ (complete)
- ~~在数据库添加pwd字段，实现 `pwd/cd/ls` 命令和相对目录的传输~~(complete)
- ~~限制 server 名称，不能有 local/localhost/: 等保留字段~~(complete)
- 增加身份验证
- ~~在添加服务器时验证可用性、用户可用性验证~~
- ~~Server 端增加配置文件，配置监听地址~~(complete)
- ~~show_server show server status~~(complete)
- Server 端做成服务
- Hash Check
- 本地、远程外链下载（http, https）
- 实现跨平台编译

## PRO Version TODO
- Transfer file directly from remote to remote 
- Transmission rate limit
- Encryption communication
- User Rights Management
- File Detail (Create time, Mod time, last open， Hash time and so on.)
- 本地、远程外链下载（http, https, ftp, sftp and so on.）
- 远程路径补全
- 用户权限
- show_server 先显示基本信息，然后`loading`加载状态
- GUI Client

## License
![GitHub](https://img.shields.io/github/license/dev-lake/FileBird?style=for-the-badge&color=green&cacheSeconds=3600)