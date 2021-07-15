package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("开启代理端口，格式—>[ exe 80]")
	defer func() {
		err:=recover()
		if err !=nil{
			fmt.Println("缺少参数使用80端口进行代理")
			Listen,_:=net.Listen("tcp","0.0.0.0:80")
			fmt.Println("成功开启代理端口",80)
			for   {
				conn,_:=Listen.Accept()
				go mid_channel(conn)
			}
		}
	}()
	Port:=os.Args[1]
	Listen,_:=net.Listen("tcp","0.0.0.0:"+Port)
	fmt.Println("成功开启代理端口",Port)
	for   {
		conn,_:=Listen.Accept()
		go mid_channel(conn)
	}

}
func mid_channel(conn net.Conn)  {
	reader := bufio.NewReader(conn)
	var method string
	var address string
	var total string
	for i:=0;i<6 ;i++ {
		msg,err:=reader.ReadString('\n')
		if i==0{
			str:= strings.Split(msg, " ")
			defer func() {
				re:=recover()
				if re!=nil{
					fmt.Println(re)
					fmt.Println("errer str:",msg)
				}
			}()
			address=str[1]
			method=str[0]
			if method!="CONNECT" {
				conn.SetReadDeadline(time.Now().Add(time.Second*1))
				OtherMethod:=msg
				for ;err==nil;{
					if strings.HasPrefix(msg,"Host: "){
						address=strings.Replace(msg,"Host: ","",1)
						address=strings.Replace(address,"\n","",2)
						address=strings.Replace(address,"\r","",2)
					}
					msg,err=reader.ReadString('\n')
					if err!=nil {
						//fmt.Println(err.Error())
					}else{
						OtherMethod+=msg
					}
				}
				fmt.Println(OtherMethod)
				go HTTPConnect(address,OtherMethod,conn)
				return
			}
		}else if(i==2){
			str := strings.Split(msg, ":")
			fmt.Println(i,str[0])
		}
		total+=msg
	}
	go Connect(address,total,conn)
}
func Connect(address string,msg string,connF net.Conn){
	connT,err:=net.Dial("tcp",address)
	if(err!=nil){
		fmt.Println(address,err.Error())
		connF.Close()
	}else{
		//fmt.Println("to target:\n",msg)
		connF.Write([]byte(msg))
		Channel(connF,connT)
	}
}
func HTTPConnect(address ,msg string,connF net.Conn)  {
		connT,err:=net.Dial("tcp", address)
		//connT.SetReadDeadline( time.Now().Add(time.Second*6))
		if(err!=nil){
			fmt.Println(address,err.Error())
			connF.Close()
		}else{
			//fmt.Println("to target:",msg)
			n,err:=connT.Write([]byte(msg))
			if err!=nil{
				fmt.Println(err.Error())
			}else{
				fmt.Println("input:",n)
				defer connT.Close()
				defer connF.Close()
				for{
					buf := [512]byte{}
					n, err := connT.Read(buf[:])
					connF.Write(buf[:n])
					//fmt.Println(string(buf[:n]))
					if err != nil {
						return
					}
				}
			}
		}
}

func Channel(source net.Conn,target net.Conn){
	go func() {
		defer source.Close()
		for{
			buf := [512]byte{}
			n, err := source.Read(buf[:])
			target.Write(buf[:n])
			if err != nil {
				fmt.Println("recv failed, err:", err)
				return
			}
			//fmt.Println(string(buf[:n]))

		}
	}()
	go func() {
		defer target.Close()
		for  {
			buf := [512]byte{}
			n, err := target.Read(buf[:])
			source.Write(buf[:n])
			if err != nil {
				fmt.Println("recv failed, err:", err)
				return
			}
			//fmt.Println(string(buf[:n]))

		}
	}()
}