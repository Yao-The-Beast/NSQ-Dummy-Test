package main

import "net"
import "fmt"
import "os"
import "time"
import "encoding/binary"

var MESSAGESIZE int = 1000
var MESSAGES int = 10000
var CLIENTS int = 128

func main(){
    
    for i:= 0; i < CLIENTS; i++ {
        go setupConnection()
    } 
    time.Sleep(time.Second * 10);
    fmt.Println("ALL END");
}


func setupConnection(){
    address := "127.0.0.1:8080"
    conn, err := net.Dial("tcp",address)
    if err != nil{
        println("CONNECTION ERROR ", err.Error())
        os.Exit(1)
    }
  
    for i:= 0; i < MESSAGES; i++ {
        
        buffer := make([]byte, MESSAGESIZE)
        binary.PutVarint(buffer, time.Now().UnixNano())
        conn.Write(buffer);
        
        time.Sleep(100 * time.Microsecond);
    }
    fmt.Print("!");
}