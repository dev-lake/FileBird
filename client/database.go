package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ServerInfo struct {
	gorm.Model
	Name string
	Addr string
	Port int
	Pass string // password
	Desc string // description
}

var (
	db_path = "filebird.db"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
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
func AddServer(db *gorm.DB, name string, addr string, port int) {
	db.Create(&ServerInfo{Name: name, Addr: addr, Port: port})
}

// show server
func ShowServer(db *gorm.DB) []ServerInfo {
	var servers []ServerInfo
	db.Find(&servers)
	return servers
}

// get server
func GetServer(db *gorm.DB, name string) ServerInfo {
	var server ServerInfo
	if name == "localhost" || name == "127.0.0.1" {
		return ServerInfo{Name: "localhost", Addr: "localhost", Port: 2000}
	}
	db.First(&server, "name = ?", name)
	if condition := db.Error; condition != nil {
		panic("failed to get server")
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
