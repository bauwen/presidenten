package server

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
    //"time"
    
    "../engine"
    "../messenger"
)

var server *messenger.Server
var playerSockets map[int]*messenger.Socket
var playerCount int
var playerNames map[int]string
var game *engine.Game
var state string
var playerDone map[int]bool
var playerPass map[int]bool
var playerTurn int
var history []string
var lastPlayer int

var callbacks map[string]func(*messenger.Socket, []string) = map[string]func(*messenger.Socket, []string){
    "connection": onConnection,
    "disconnection": onDisconnection,
    "name": onName,
    "play": onPlay,
    "exchange": onExchange,
}

func StartServer(port string) {
    var err error
    playerNames = make(map[int]string)
    playerSockets = make(map[int]*messenger.Socket)
    playerDone = make(map[int]bool)
    playerPass = make(map[int]bool)
    lastPlayer = -1
    state = "lobby"
    
    server, err = messenger.CreateServer(port)
    if err != nil {
        fmt.Println("Error: kan server niet starten op poortnummer " + port)
        return
    }
    fmt.Println("\nWelkom bij (c) The Presidenten Server!")
    fmt.Println("De server is aan het wachten op spelers op poort " + port + "!")
    
    fmt.Println("\nTyp 'start' om te starten met de huidige spelers in de lobby")
    fmt.Println("Typ 'run' om te starten zonder vertraging bij het spelen met bots")
    fmt.Println("\nLobby:\n")
    go handleInput()
    
    server.Run(callbacks)
}

func handleInput() {
    r := bufio.NewReader(os.Stdin)
    for {
        input, _ := r.ReadString('\n')
        input = strings.TrimSpace(input)
        if input == "start" && state == "lobby" {
            state = "play"
            startGame(true)
            startMatch()
        }
        if input == "run" && state == "lobby" {
            state = "play"
            startGame(false)
            startMatch()
        }
    }
}

func startGame(delayed bool) {
    game = engine.NewGame(playerCount)

    d := "1"
    if !delayed { d = "0" }
    server.Broadcast("settings", d)
}

func startMatch() {
    history = []string{}
    fmt.Println("\nDe kaarten worden uitgedeeld...")
    for i := 0; i < playerCount; i++ {
        parts := []string{ fmt.Sprint(i) }
        for j := 0; j < playerCount; j++ {
            parts = append(parts, playerNames[j])
        }
        socket := playerSockets[i]
        socket.Send("start", parts...)
        playerDone[i] = false
        playerPass[i] = false
    }
    state = "dummy"
    
    //time.Sleep(1 * time.Second)
    
    receiver := game.StartMatch()
    originalReceiver := receiver
    extraCount := 52 % playerCount
    extraReceivers := make([]int, extraCount)
    for i := 0; i < extraCount; i++ {
        extraReceivers[i] = receiver
        receiver = (receiver + 1) % playerCount
    }
    exchange := game.CurrentPresident >= 0
    if exchange {
        state = "exchange"
    } else {
        state = "play"
    }
    for i := 0; i < playerCount; i++ {
        parts := []string{}
        if exchange {
            parts = append(parts, fmt.Sprint(game.CurrentShit), fmt.Sprint(game.CurrentPresident))
        }
        cards := game.Cards[i]
        for j := 0; j < len(cards); j++ {
            parts = append(parts, fmt.Sprint(cards[j]))
        }
        parts = append(parts, "#")
        for j := 0; j < len(extraReceivers); j++ {
            parts = append(parts, fmt.Sprint(extraReceivers[j]))
        }
        socket := playerSockets[i]
        if exchange {
            socket.Send("distribution-exchange", parts...)
        } else {
            socket.Send("distribution-only", parts...)
        }
    }
    
    if !exchange {
        playerTurn = (originalReceiver + playerCount - 1) % playerCount
        nextTurn()
    } else {
        playerTurn = game.MatchLosers[len(game.MatchLosers) - 1]
        playerTurn = (playerTurn + playerCount - 1) % playerCount
    }
}

func getPlayerFromSocket(socket *messenger.Socket) int {
    for i, s := range playerSockets {
        if s == socket {
            return i
        }
    }
    return -1
}

func onConnection(socket *messenger.Socket, args []string) {
    playerSockets[playerCount] = socket
    playerCount += 1
}

func onDisconnection(socket *messenger.Socket, args []string) {
    id := getPlayerFromSocket(socket)
    name := playerNames[id]
    socket.Broadcast("left", name)
    fmt.Println("Iemand heeft zich uitgelogd: " + name)
    os.Exit(1)
}

func onName(socket *messenger.Socket, args []string) {
    name := args[0]
    id := getPlayerFromSocket(socket)
    playerNames[id] = name
    fmt.Println("\t" + name)
    var names []string
    for i := 0; i < playerCount; i++ {
        names = append(names, playerNames[i])
    }
    socket.Send("lobby", names...)
    socket.Broadcast("joined", name)
}

func onPlay(socket *messenger.Socket, args []string) {
    id := getPlayerFromSocket(socket)
    if args[0] == "pas" {
        playerPass[id] = true
        server.Broadcast("passed", fmt.Sprint(id), fmt.Sprint(len(game.Cards[id])))
        nextTurn()
    } else {
        lastPlayer = id
        var cards []int
        for i := 0; i < len(args); i++ {
            index, _ := strconv.Atoi(args[i])
            cards = append(cards, index)
        }
        parts := []string{ fmt.Sprint(id), fmt.Sprint(len(game.Cards[id]) - len(args)) }
        for _, arg := range args {
            parts = append(parts, arg)
            history = append(history, arg)
        }
        server.Broadcast("play", parts...)
        gameOver, done, first := game.Play(id, cards)
        if done {
            playerDone[id] = true
            s := "burger"
            if first {
                s = "president"
            }
            server.Broadcast("done", fmt.Sprint(id), s)
        }
        if gameOver {
            stats := []string{}
            for i := 0; i < playerCount; i++ {
                p := game.GetPresidentCount(i)
                c := game.GetCivilianCount(i)
                s := game.GetShitCount(i)
                stats = append(stats, fmt.Sprint(p), fmt.Sprint(c), fmt.Sprint(s))
            }
            printStats(stats)
            server.Broadcast("over", stats...)

            //time.Sleep(2 * time.Second)
            startMatch()
        } else {
            nextTurn()
        }
    }
}

func nextTurn() {
    for {
        playerTurn = (playerTurn + 1) % playerCount
        if playerPass[playerTurn] {
            continue
        }
        if playerDone[playerTurn] {
            if playerTurn != lastPlayer {
                continue
            }
            for {
                playerTurn = (playerTurn + 1) % playerCount
                if !playerDone[playerTurn] {
                    lastPlayer = playerTurn
                    break
                }
            }
            break
        }
        break
    }

    if playerTurn == lastPlayer {
        lastPlayer = -1
        server.Broadcast("sweep", fmt.Sprint(playerTurn))
        for i := 0; i < playerCount; i++ {
            playerPass[i] = false
        }
    }

    /*
    for {
        playerTurn = (playerTurn + 1) % playerCount
        if !playerPass[playerTurn] && !playerDone[playerTurn] {
            break
        }
    }
    playing := 0
    for i := 0; i < playerCount; i++ {
        if !playerPass[i] && !playerDone[i] {
            playing += 1
        }
    }
    if playing == 1 {
        server.Broadcast("sweep", fmt.Sprint(playerTurn))
        for i := 0; i < playerCount; i++ {
            playerPass[i] = false
        }
    }
    */
    parts := []string{ fmt.Sprint(playerTurn) }
    for i := 0; i < playerCount; i++ {
        parts = append(parts, fmt.Sprint(len(game.Cards[i])))
    }
    for _, h := range history {
        parts = append(parts, h)
    }
    parts = append(parts, "#")
    for _, card := range game.Cards[playerTurn] {
        parts = append(parts, fmt.Sprint(card))
    }
    server.Broadcast("turn", parts...)
}

func onExchange(socket *messenger.Socket, args []string) {
    index, _ := strconv.Atoi(args[0])
    badCard, bestCard := game.ExchangeCards(index)
    server.Broadcast("exchange", fmt.Sprint(badCard), fmt.Sprint(bestCard), fmt.Sprint(game.CurrentShit), fmt.Sprint(game.CurrentPresident))
    nextTurn()
}

func printStats(args []string) {
    fmt.Println("\n=======================================================")
    fmt.Println("De match is voorbij! Hier zijn de huidige statistieken:\n")
    fmt.Print(strings.Repeat(" ", 17) + "\t")
    fmt.Print("President" + strings.Repeat(" ", 3))
    fmt.Print("Burger" + strings.Repeat(" ", 6))
    fmt.Print("Shit" + strings.Repeat(" ", 8))
    fmt.Println("\n")
    k := 0
    for i := 0; i < playerCount; i++ {
        name := playerNames[i]
        name = padName(name)
        fmt.Print(name + ":\t")
        p := args[k]
        k += 1
        fmt.Print(p + strings.Repeat(" ", 12 - len(p)))
        c := args[k]
        k += 1
        fmt.Print(c + strings.Repeat(" ", 12 - len(c)))
        s := args[k]
        k += 1
        fmt.Print(s + strings.Repeat(" ", 12 - len(s)))
        fmt.Println("")
    }
    fmt.Println("=======================================================\n")
}

func padName(name string) string {
    return strings.Repeat(" ", 16 - len(name)) + name
}
