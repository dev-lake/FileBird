package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ServerInfo struct {
	gorm.Model
	Name string `gorm:"unique"`
	Addr string
	Port int
	User string // user name
	Pass string // password
	Home string // home directory
	Pwd  string // current directory
	Desc string // description
}

var (
	db_path = "data/filebird.db"
)

func InitDB() *gorm.DB {
	exe_path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exe_dir := filepath.Dir(exe_path)
	os.MkdirAll("data", 0755)
	db, err := gorm.Open(sqlite.Open(filepath.Join(exe_dir, db_path)), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&ServerInfo{})
	return db
}

// delete server
func DeleteServer(db *gorm.DB, name string) {
	db.Delete(&ServerInfo{}, "name = ?", name)
}

// add server
func AddServer(db *gorm.DB, server_name string, addr string, port int) error {
	// check server name
	err := CheckServerName(server_name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	server := ServerInfo{Name: server_name, Addr: addr, Port: port}
	user_info, err := GetRemoteUserInfo(server)
	if err != nil {
		fmt.Println("User Check Failed, Continue Add Server? (y/n): ")
		var input string
		fmt.Scanln(&input)
		if input == "y" {
			db.Create(&server)
			fmt.Println("Server Added.")
			return err
		}
		fmt.Println("Server Add Cancled.")
	}
	db.Create(
		&ServerInfo{
			Name: server_name,
			Addr: addr,
			Port: port,
			User: user_info.Username,
			Home: user_info.HomeDir,
			Pwd:  user_info.HomeDir,
		},
	)
	fmt.Println("Server Added.")
	return nil
}

// show server
func ShowServer(db *gorm.DB) []ServerInfo {
	var servers []ServerInfo
	db.Find(&servers)
	return servers
}

// get server pwd field
func GetServerPwd(db *gorm.DB, name string) string {
	var server ServerInfo
	db.Where("name = ?", name).First(&server)
	return server.Pwd
}

// get server
func GetServer(db *gorm.DB, name string) ServerInfo {
	var server ServerInfo
	if name == "localhost" || name == "127.0.0.1" {
		return ServerInfo{Name: "localhost", Addr: "localhost", Port: 2000}
	}
	db.First(&server, "name = ?", name)
	if condition := db.Error; condition != nil {
		log.Fatalln("failed to get server")
	}
	return server
}

// modify server
func ModifyServerAddr(db *gorm.DB, name string, addr string) {
	db.Model(&ServerInfo{}).Where("name = ?", name).Update("addr", addr)
}

// modify server port
func ModifyServerPort(db *gorm.DB, name string, port int) {
	db.Model(&ServerInfo{}).Where("name = ?", name).Update("port", port)
}

// modify server name
func ModifyServerName(db *gorm.DB, name string, new_name string) {
	db.Model(&ServerInfo{}).Where("name = ?", name).Update("name", new_name)
}

// modify server pwd
func ModifyServerPwd(db *gorm.DB, name string, pwd string) {
	db.Model(&ServerInfo{}).Where("name = ?", name).Update("pwd", pwd)
}
