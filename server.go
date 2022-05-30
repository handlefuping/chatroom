package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	OnlineMap map[string]*User
	Ip string
	Port int
	lock sync.RWMutex
	Ch chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Ch: make(chan string),
	}
}
func (receiver *Server) handleConnection(conn net.Conn)  {

	user := NewUser(conn, receiver)
	user.OnLine()


		for {
			buf := make([]byte, 10 * 1024)
			err := user.Conn.SetReadDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				fmt.Println("用户超时", err)
				user.OffLine()
				break
			}
			n, err := user.Conn.Read(buf)
			if err != nil && err != io.EOF{
				fmt.Println("接受客户端数据错误", err)
				break
			}
			if n == 0 {
				user.OffLine()
				break
			}
			fmt.Println("接收客户端字节长度：", n)
			receiver.BroadCast(user, string(buf))
		}

}
func (receiver *Server) BroadCast(user *User, msg string)  {
	sendMsg := user.Name + ":" + msg + "\n"
	receiver.Ch <- sendMsg
}

func (receiver *Server) ListenMsg()  {
	go func() {
		for  {
			msg := <- receiver.Ch
			receiver.lock.Lock()
			// todo 增加消息协议
			for _, user := range receiver.OnlineMap {
				user.Ch <- msg
			}

			receiver.lock.Unlock()
		}
	}()

}

func (receiver *Server) Run()  {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", receiver.Ip, receiver.Port))
	if err != nil {
		fmt.Println("聊天室启动失败：", err)
		return
	}
	defer ln.Close()
	receiver.ListenMsg()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("进入聊天室失败：", err)
		}
		fmt.Println("进入聊天室")
		go receiver.handleConnection(conn)
	}

	fmt.Println("main end")
}