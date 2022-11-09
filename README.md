# File Bird
**A Efficient File Transfer Software**

### Features
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

### TODO
- ~~Show Progress bar~~ (complete)
- ~~实现文件夹的递归传输~~ (complete)
- ~~在数据库添加pwd字段，实现 `pwd/cd/ls` 命令和相对目录的传输~~(complete)
- ~~限制 server 名称，不能有 local/localhost/: 等保留字段~~(complete)
- 增加身份验证
- ~~在添加服务器时验证可用性、用户可用性验证~~
- Server 端增加配置文件，配置监听地址、用户权限等
- Server 端做成服务
- Hash Check
- 本地、远程外链下载（http, https）

### PRO Version TODO
- Transfer file directly from remote to remote 
- Transmission rate limit
- Encryption communication
- User Rights Management
- File Detail (Create time, Mod time, last open， Hash time and so on.)
- 本地、远程外链下载（http, https, ftp, sftp and so on.）
- 远程路径补全
- GUI Client

## Compile
### Server
```bash
cd server/
go mod tidy
go build -o filebird-server
```
### client
```bash
cd client/
go mod tidy
go build -o filebird
```

## Usage
### Add & Show Server
```bash
# add server
./filebird add_server -n SERVERNAME -a IP -p PORT

# show server
./filebird show_server
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

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=dev-lake/FileBird&type=Date)](https://star-history.com/#dev-lake/FileBird&Date)
