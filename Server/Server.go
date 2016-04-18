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
    var latencies []int64
    buf := make([]byte,MESSAGESIZE)
    
    for i:=0; i < MESSAGES; i++ {
        n, err := c.Read(buf)
        if err != nil || n != MESSAGESIZE {
            fmt.Println("HANDLE ERROR")
            c.Close()
            break
        }else {
           currentTime := time.Now().UnixNano()
	       sentTime, _ := binary.Varint(buf)
           latency := currentTime - sentTime
           latencies = append(latencies,latency)
        }  
    }
    if id == 1 {
        sum := int64(0)
        for _, latency := range latencies {
                sum += latency
        }
        averageLatency := int64(sum) / int64(len(latencies)) / 1000
        fmt.Printf("Latency is %d \n", averageLatency)
    }
}