# File Bird
**A Efficient File Transfer Software**

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

