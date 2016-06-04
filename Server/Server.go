package main

import "net"
import "fmt"
import "os"
import "time"
import "strconv"
import "io/ioutil"
import "encoding/binary"

var MESSAGESIZE int = 1024
var MESSAGES int = 10000
var CLIENTS int = 128

func main(){
    ln, err := net.Listen("tcp", ":8080")  
    if err != nil{
        fmt.Println(err)
        os.Exit(1)
    }
    id := 0
    
    for {
        if id == CLIENTS {
            break
        }
        conn, err := ln.Accept()
        if err != nil{
            fmt.Println(err)
            continue
        }
        id++
        go handleConnection(conn, id)       
    }
    
    time.Sleep(10 * time.Second);
    defer ln.Close()
    
}

func handleConnection(c net.Conn, id int){
    
    var oneWayLatencies []byte
    
    for i:=0; i < MESSAGES; i++ {
        buf := make([]byte,MESSAGESIZE)
        //receive message
        n, err := c.Read(buf)
        if err != nil || n != MESSAGESIZE {
            fmt.Println("HANDLE ERROR")
            c.Close()
            break
        }else {
           currentTime := time.Now().UnixNano()
	       sentTime, _ := binary.Varint(buf)
           latency := currentTime - sentTime
           x := strconv.FormatInt(latency,10)
           oneWayLatencies = append(oneWayLatencies,x...)
           oneWayLatencies = append(oneWayLatencies,"\n"...)
        }  
        
        //send message
        buffer2 := make([]byte,MESSAGESIZE)
        binary.PutVarint(buffer2, time.Now().UnixNano())
        c.Write(buffer2)
    }
    ioutil.WriteFile("client_server",oneWayLatencies,0777)
    fmt.Println("END")
}