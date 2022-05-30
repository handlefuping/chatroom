package main

import (
	"fmt"
	"net"
)

type User struct {
	Addr string
	Name string
	Conn net.Conn
	Ch chan string
	Server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()
	return &User{
		Addr: addr,
		Name: addr,
		Conn: conn,
		Ch: make(chan string),
		Server: server,
	}
}

func (receiver *User) OnLine()  {
	// 将当前用户加入服务
	receiver.Server.lock.Lock()
	receiver.Server.OnlineMap[receiver.Name] = receiver
	receiver.Server.lock.Unlock()
	// 用户接受服务气消息
	receiver.ListenMsg()
	// 通知用户上线
	receiver.Server.BroadCast(receiver, "用户上线")

}

func (receiver *User) OffLine()  {
	fmt.Println("offline")
	// 将用户剔除服务
	receiver.Server.lock.Lock()
	delete(receiver.Server.OnlineMap, receiver.Name)
	receiver.Server.lock.Unlock()
	// 关闭消息接受通道
	close(receiver.Ch)
	// 关闭连接
	receiver.Conn.Close()
	// 通知用户下线
	receiver.Server.BroadCast(receiver, "用户下线")


}

func (receiver *User) ListenMsg()  {
	go func() {
		for msg := range receiver.Ch {
			receiver.Conn.Write([]byte(msg))
		}
	}()
}