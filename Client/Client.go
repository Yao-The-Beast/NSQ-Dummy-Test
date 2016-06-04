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

func main(){
    
    var roundTripLatencies []byte
    var oneWayLatencies []byte
    
    address := "127.0.0.1:8080"
    conn, err := net.Dial("tcp",address)
    if err != nil{
        println("CONNECTION ERROR ", err.Error())
        os.Exit(1)
    }
  
    for i:= 0; i < MESSAGES; i++ {
        
        buffer := make([]byte, MESSAGESIZE)
        binary.PutVarint(buffer, time.Now().UnixNano())
        
        //write to the connection
        conn.Write(buffer);
        currentTime1 := time.Now().UnixNano()
        
        //listen for the callback
        bufferRec := make([] byte, MESSAGESIZE)
        n, err := conn.Read(bufferRec)
        currentTime2 := time.Now().UnixNano()
        
        if err != nil || n != MESSAGESIZE {
            fmt.Println("HANDLE ERROR")
            conn.Close()
            break
        }else {
	       serverSentTime, _ := binary.Varint(bufferRec)
           //append one way latency
           oneWayLatency := currentTime2 - serverSentTime
           x := strconv.FormatInt(oneWayLatency,10)
           oneWayLatencies = append(oneWayLatencies,x...)
           oneWayLatencies = append(oneWayLatencies,"\n"...)
           
           //append round trip latency  
           roundTripLatency := currentTime2 - currentTime1
           y := strconv.FormatInt(roundTripLatency,10)
           roundTripLatencies = append(roundTripLatencies,y...)
           roundTripLatencies = append(roundTripLatencies,"\n"...)         
        }  
        time.Sleep(100 * time.Microsecond)
    }
    //write output to files
    ioutil.WriteFile("server_client", oneWayLatencies, 0777)
    ioutil.WriteFile("roundtrip", roundTripLatencies, 0777)
    
    fmt.Println("END");
}