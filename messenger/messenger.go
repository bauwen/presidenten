package messenger

import (
    "bufio"
    "container/list"
    "net"
    "strings"
)

var DELIMITER byte = '\x1f'

type Server struct {
    ln    net.Listener
    conns *list.List
}

type Socket struct {
    Server    *Server
    element   *list.Element
    conn      net.Conn
    callbacks map[string]func(*Socket, []string)
}

func (server *Server) Connections() []net.Conn {
    var conns []net.Conn
    for e := server.conns.Front(); e != nil; e = e.Next() {
        conns = append(conns, e.Value.(net.Conn))
    }
    return conns
}

func (server *Server) Count() int {
    return server.conns.Len()
}

func (server *Server) Broadcast(id string, parts ...string) {
    message := []byte(id + " " + strings.Join(parts, " ") + string(DELIMITER))
    for e := server.conns.Front(); e != nil; e = e.Next() {
        conn := e.Value.(net.Conn)
        conn.Write(message)
    }
}

func (server *Server) Run(callbacks map[string]func(*Socket, []string)) {
    for {
        conn, err := server.ln.Accept()
        if err != nil {
            continue
        }
        element := server.conns.PushBack(conn)
        socket := &Socket{
            Server: server,
            element: element,
            conn: conn,
            callbacks: callbacks,
        }
        if callback, ok := socket.callbacks["connection"]; ok {
            callback(socket, []string{})
        }
        go handleSocketConnection(socket)
    }
}

func (socket *Socket) Send(id string, parts ...string) {
    message := []byte(id + " " + strings.Join(parts, " ") + string(DELIMITER))
    socket.conn.Write(message)
}

func (socket *Socket) Broadcast(id string, parts ...string) {
    server := socket.Server
    if server == nil {
        return
    }
    message := []byte(id + " " + strings.Join(parts, " ") + string(DELIMITER))
    for e := server.conns.Front(); e != nil; e = e.Next() {
        if e != socket.element {
            conn := e.Value.(net.Conn)
            conn.Write(message)
        }
    }
}

func CreateServer(port string) (*Server, error) {
    ln, err := net.Listen("tcp", ":" + port)
    if err != nil {
        return nil, err
    }
    server := &Server{ ln, list.New() }
    return server, nil
}

func Connect(address string, callbacks map[string]func(*Socket, []string)) (*Socket, error) {
    conn, err := net.Dial("tcp", address)
    if err != nil {
        return nil, err
    }
    socket := &Socket{
        Server: nil,
        element: nil,
        conn: conn,
        callbacks: callbacks,
    }
    if callback, ok := socket.callbacks["connection"]; ok {
        callback(socket, []string{})
    }
    go handleSocketConnection(socket)
    return socket, nil
}

func handleSocketConnection(socket *Socket) {
    server := socket.Server
    r := bufio.NewReader(socket.conn)
    for {
        message, err := r.ReadString(DELIMITER)
        // TODO: lock socket.callbacks (mutex) and defer unlock
        if err != nil {
            if server != nil {
                server.conns.Remove(socket.element)
            }
            if callback, ok := socket.callbacks["disconnection"]; ok {
                callback(socket, []string{})
            }
            break
        }
        message = strings.Trim(message, string(DELIMITER))
        parts := strings.Split(message, " ")
        id := parts[0]
        if callback, ok := socket.callbacks[id]; ok {
            callback(socket, parts[1:])
        }
    }
}
