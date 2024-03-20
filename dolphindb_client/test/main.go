package main

import (
	"context"
	"fmt"
	"github.com/dolphindb/api-go/api"
)

func main() {
	conn, err := api.NewDolphinDBClient(context.TODO(), "localhost:8848", nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	err = conn.Connect()
	if err != nil {
		fmt.Println(err)
	}

	lr := &api.LoginRequest{
		UserID:   "admin",
		Password: "123456",
	}
	err = conn.Login(lr)
	if err != nil {
		fmt.Println(err)
	}

	// Create database
	db, err := conn.Database(&api.DatabaseRequest{
		Directory: "/opt/server/data",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.CreateTable(&api.CreateTableRequest{})
	tb, err := conn.LoadTable(&api.LoadTableRequest{
		Database:  "/opt/server/data",
		TableName: "tb",
	})
	print(tb)

	fmt.Println("Create table")
	//res, err := conn.RunScript(fmt.Sprintf("database(\"/opt/server/data\");select count(id) from tb;", "/opt/server/data"))
	if err != nil {
		fmt.Println("aaaa", err)
		return
	}
	//fmt.Println(res)
	// Disconnect
	conn.Close()
}
