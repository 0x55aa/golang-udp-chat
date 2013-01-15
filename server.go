package main

import(
    "fmt"
    "net"
    "os"
    "time"
    "strings"
    "strconv"
)

const(
    port string = ":1200"
)

type Server struct{
    conn *net.UDPConn //
    messages chan string //接受到的消息
    clients map [int]Client
}

type Client struct{
    userID int
    userName string
    userAddr *net.UDPAddr

}
type Message struct{
    status int
    userID int
    userName string
    content string
}

func (s *Server) handleMessage(){
    var buf [512]byte

    fmt.Println("1\n")
    n, addr, err := s.conn.ReadFromUDP(buf[0:])
    if err != nil{
        return
    }
    fmt.Println("2\n")
    //分析消息
    msg := string(buf[0:n])
    fmt.Println(msg)
    m := s.analyzeMessage(msg)
    switch m.status{
        //进入聊天室消息
        case 1:
            var c Client
            c.userAddr = addr
            c.userID = m.userID
            c.userName = m.userName
            s.clients[c.userID] = c //添加用户
            s.messages <- msg
            fmt.Println("3")
        //用户发送消息
        case 2:
            s.messages <- msg
        //client发来的退出消息
        case 3:
            delete(s.clients, m.userID)
            s.messages <- msg
        default:
            fmt.Println("未识别消息", msg)
    }

    //fmt.Println(n,addr,string(buf[0:n]))


}
//这里还要判断一下数组的长度，
func (s *Server) analyzeMessage(msg string) (m Message) {
    //var m Message
    s1 := strings.Split(msg,"###")

        s2 := strings.Split(s1[1],"##")
        //fmt.Println(s2)
        switch s2[0]{
            case "1":
                m.status,_ = strconv.Atoi(s2[0])
                fmt.Println("44")
                m.userID,_ = strconv.Atoi(s2[1])
                m.userName = s2[2]
                fmt.Println(m)
                return
            case "2":
                m.status,_ = strconv.Atoi(s2[0])
                m.userID,_ = strconv.Atoi(s2[1])
                m.content = s2[2]
                return
            case "3":
                m.status,_ = strconv.Atoi(s2[0])
                m.userID,_ = strconv.Atoi(s2[1])
                return
            default:
                fmt.Println("未识别消息", msg)
                return
        }
    return
}
func (s *Server) sendMessage() {
    for{
        msg := <- s.messages
        daytime := time.Now().String()
        sendstr := msg + daytime
        fmt.Println(00,sendstr)
        for _,c := range s.clients {
            fmt.Println(c)
            n,err := s.conn.WriteToUDP([]byte(sendstr),c.userAddr)
            fmt.Println(n,err)
        }
    }

}

func checkError(err error){
    if err != nil{
        fmt.Fprintf(os.Stderr,"Fatal error:%s",err.Error())
        os.Exit(1)
    }
}

func main(){
    udpAddr, err := net.ResolveUDPAddr("udp4",port)
    checkError(err)

    var s Server
    s.messages = make(chan string,20)
    s.clients =make(map[int]Client,0)

    s.conn,err = net.ListenUDP("udp",udpAddr)
    checkError(err)

    go s.sendMessage()

    for{
        s.handleMessage()
    }
}
