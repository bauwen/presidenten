package client

import (
    "bufio"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "os/exec"
    "sort"
    "strconv"
    "strings"

    "time"
    
    "../messenger"
)

var name string
var myId int
var playerCount int
var playerNames map[int]string
var currentCards []int
var lastCards []int
var delayed bool

var botPlay bool
var lastPlay []string
var pipe io.Reader
var stdin io.Writer
var stderr io.Reader
var botProgram *exec.Cmd

var callbacks map[string]func(*messenger.Socket, []string) = map[string]func(*messenger.Socket, []string){
    "disconnection": onDisconnection,
    "lobby": onLobby,
    "joined": onJoined,
    "left": onLeft,
    "start": onStart,
    "distribution-only": onDistributionOnly,
    "distribution-exchange": onDistributionExchange,
    "turn": onTurn,
    "passed": onPassed,
    "play": onPlay,
    "done": onDone,
    "over": onOver,
    "sweep": onSweep,
    "exchange": onExchange,
    "settings": onSettings,
}

func StartGame(port string, isBot bool, botCommand string, botFile string) {
    delayed = true
    botPlay = isBot
    pipe = os.Stdin
    if botPlay {
        cmd := exec.Command(botCommand, botFile)
        botProgram = cmd
        stdin, _ = cmd.StdinPipe()
        pipe, _ = cmd.StdoutPipe()
        stderr, _ = cmd.StderrPipe()
        cmd.Start()
    }

    // show intro
    fmt.Println("\nWelkom bij (c) The Presidenten Game!\n")
    
    // ask server address
    address := prompt("Geef het IP-adres op van de server (laat leeg voor 'localhost'): ")
    if len(address) == 0 {
        address = "localhost"
    }
    if !strings.Contains(address, ":") {
        address += ":" + port
    }
    
    // connect to server
    socket, err := messenger.Connect(address, callbacks)
    if err != nil {
        fmt.Println("Error: kan niet verbinden met server op adres " + address)
        return
    }
    fmt.Println("Je bent succesvol verbonden met de server!\n")
    
    // ask player name
    for {
        name = prompt("Geef je naam op (max. 10 karakters, geen spaties): ")
        if len(name) > 10 {
            name = name[:10]
        }
        name = strings.Split(name, " ")[0]
        name = strings.TrimSpace(name)
        if len(name) > 0 {
            break
        }
        fmt.Println("Dat is een ongeldige naam! Probeer opnieuw...\n")
    }
    socket.Send("name", name)
    //fmt.Println("Dag " + name + "!\n")
    
    // block main thread (callbacks determine flow from here on)
    select {
    }
}

func prompt(text string) string {
    fmt.Print(text)
    r := bufio.NewReader(os.Stdin)
    input, _ := r.ReadString('\n')
    input = strings.TrimSpace(input)
    return input
}

func promptVia(text string) string {
    fmt.Print(text)
    r := bufio.NewReader(pipe)
    input, _ := r.ReadString('\n')
    input = strings.TrimSpace(input)
    return input
}

func printExtraReceivers(receivers []string) {
    if len(receivers) == 0 {
        return
    }
    s := ""
    if len(receivers) == 1 {
        s += receivers[0] + " kreeg "
    } else {
        for i := 0; i < len(receivers); i++ {
            s += receivers[i]
            if i < len(receivers) - 2 {
                s += ", "
            } else if i == len(receivers) - 2 {
                s += " en "
            } 
        }
        s += " kregen "
    }
    fmt.Println(s + "een extra kaart.")
}

func printCurrentCards() {
    if len(currentCards) == 0 {
        fmt.Println("| ")
        return
    }
    s := ""
    prevCard := currentCards[0]
    for i := 0; i < len(currentCards); i++ {
        card := currentCards[i]
        if card != prevCard {
            s += "| "
        }
        s += cardIndexToName(card) + " "
        prevCard = card
    }   
    fmt.Println(s)
}

func onDisconnection(socket *messenger.Socket, args []string) {
    fmt.Println("\n\nError: de verbinding met de server is verbroken")
    if botPlay {
        botProgram.Process.Kill()
    }
    os.Exit(1)
}

func onLobby(socket *messenger.Socket, args []string) {
    fmt.Println("\nLobby:\n")
    for _, name := range args {
        fmt.Println("\t" + name)
    }
}

func onJoined(socket *messenger.Socket, args []string) {
    fmt.Println("\t" + args[0])
}

func onLeft(socket *messenger.Socket, args []string) {
    name := args[0]
    fmt.Println("\n\nError: " + name + " is uitgelogd. Het spel is beëindigd")
    if botPlay {
        botProgram.Process.Kill()
    }
    os.Exit(1)
}

func onStart(socket *messenger.Socket, args []string) {
    //fmt.Println("\nDe kaarten worden uitgedeeld...")
    myId, _ = strconv.Atoi(args[0])
    playerCount = len(args) - 1
    playerNames = make(map[int]string)
    for i := 1; i < len(args); i++ {
        playerNames[i-1] = args[i]
    }
    lastCards = []int{ 1, -1 }
    lastPlay = []string{}
}

func onDistributionOnly(socket *messenger.Socket, args []string) {
    currentCards = []int{}
    i := 0
    for ; i < len(args); i++ {
        arg := args[i]
        if arg == "#" {
            i += 1
            break
        }
        card, _ := strconv.Atoi(arg)
        currentCards = append(currentCards, card)
    }
    receivers := []string{}
    for ; i < len(args); i++ {
        id, _ := strconv.Atoi(args[i])
        receivers = append(receivers, playerNames[id])
    }

    fmt.Println("\nDe kaarten worden uitgedeeld...")
    printExtraReceivers(receivers)
    fmt.Println("Dit zijn jouw kaarten:")
    printCurrentCards()
    fmt.Println("")
}

func onDistributionExchange(socket *messenger.Socket, args []string) {
    currentCards = []int{}
    shit, _ := strconv.Atoi(args[0])
    pres, _ := strconv.Atoi(args[1])
    i := 2
    for ; i < len(args); i++ {
        arg := args[i]
        if arg == "#" {
            i += 1
            break
        }
        card, _ := strconv.Atoi(arg)
        currentCards = append(currentCards, card)
    }
    receivers := []string{}
    for ; i < len(args); i++ {
        id, _ := strconv.Atoi(args[i])
        receivers = append(receivers, playerNames[id])
    }
    fmt.Println("De kaarten worden uitgedeeld...")
    printExtraReceivers(receivers)
    fmt.Println("Dit zijn jouw kaarten:")
    printCurrentCards()
    fmt.Println("")

    if myId == shit {
        name := playerNames[pres]
        index := currentCards[len(currentCards) - 1]
        best := cardIndexToName(index)
        fmt.Println("Jij bent de shit. Je beste kaart wordt geruild met " + name + ": " + best)
    } else if myId == pres {
        name := playerNames[shit]
        fmt.Println("Jij bent de president!")
        x := -1
        index := -1
        for {
            if botPlay {
                s := "exchange"
                for _, c := range currentCards {
                    s += " " + fmt.Sprint(c)
                }
                stdin.Write([]byte(s + "\n"))
            }
            input := promptVia("Geef op welke kaart je met " + name + " wilt ruilen: ")
            input = strings.TrimSpace(input)
            x = cardNameToIndex(input)
            if x > 2 {
                has := false
                for i, card := range currentCards {
                    if card == x {
                        has = true
                        index = i
                        break
                    }
                }
                if has {
                    break
                } else {
                    fmt.Println("Gelieve een kaart op te geven die je in je bezit hebt...")
                }
            } else {
                fmt.Println("Gelieve een geldige kaart op te geven...")
                if botPlay {
                    bytes, _ := ioutil.ReadAll(stderr)
                    fmt.Println("")
                    fmt.Println(string(bytes))
                    botProgram.Process.Kill()
                    os.Exit(1)
                }
            }
        }
        if botPlay {
            if delayed {
                time.Sleep(5 * time.Second)
            }
        }
        socket.Send("exchange", fmt.Sprint(index))
    } else {
        nameShit := playerNames[shit]
        namePres := playerNames[pres]
        fmt.Println("De shit " + nameShit + " en de president " + namePres + " zijn een kaart aan het ruilen...")
    }
}

func onTurn(socket *messenger.Socket, args []string) {
    id, _ := strconv.Atoi(args[0])
    if id != myId {
        //name := playerNames[id]
        //fmt.Println(name + " is aan zet...")
        return
    }

    cardCountAll := []string{}
    for i := 0; i < playerCount; i++ {
        cardCountAll = append(cardCountAll, args[1 + i])
    }
    cardCount := []string{}
    ii := id
    for i := 0; i < playerCount - 1; i++ {
        ii = (ii + 1) % playerCount
        cardCount = append(cardCount, cardCountAll[ii])
    }

    history := []string{}
    k := 1 + playerCount
    for {
        token := args[k]
        if token == "#" {
            break
        }
        history = append(history, token)
        k += 1
    }

    fmt.Println("")
    //fmt.Println("Jij bent aan zet!")
    currentCards = []int{}
    for _, arg := range args[k+1:] {
        index, _ := strconv.Atoi(arg)
        currentCards = append(currentCards, index)
    }
    printCurrentCards()

    var cards []string
    for {
        if botPlay {
            s := "play"
            for _, h := range history {
                s += " " + h
            }
            s += " #"
            for _, l := range lastPlay {
                s += " " + l
            }
            s += " #"
            for _, n := range cardCount {
                s += " " + n
            }
            s += " #"
            for _, c := range currentCards {
                s += " " + fmt.Sprint(c)
            }
            stdin.Write([]byte(s + "\n"))
            //fmt.Println("\n(DEBUG)PIPED: " + s)
        }
        var err int
        input := promptVia("Jij bent aan zet! Geef op welke kaarten je wilt spelen: ")
        cards, err = parseCardInput(input)
        if err == 0 {
            break
        }
        if err == 1 {
            fmt.Println("Daar zit een kaart bij die je niet hebt...")
        }
        if err == 2 {
            fmt.Println("Dat is ongeldige invoer! Geef geldige kaarten op, gescheiden door spaties...")
        }
        if err == 3 {
            fmt.Println("Die kaartencombinatie kan niet gelegd worden op de laatste zet...")
        }
        if err == 4 {
            fmt.Println("Je hebt geen kaarten opgegeven. Typ 'pas' indien je wilt passen...")
        }
        if err == 5 {
            fmt.Println("Dat is een ongeldige kaartencombinatie! Probeer opnieuw...")
        }
        if err == 6 {
            fmt.Println("Je kan niet passen voor de eerste zet...")
        }
        if botPlay {
            bytes, _ := ioutil.ReadAll(stderr)
            fmt.Println("")
            fmt.Println(string(bytes))
            botProgram.Process.Kill()
            os.Exit(1)
        }
    }
    if botPlay {
        fmt.Println("")
        if delayed {
            time.Sleep(1 * time.Second)
        }
    }
    if len(cards) == 0 {
        //fmt.Println("Je hebt gepast!")
        fmt.Println("")
        socket.Send("play", "pas")
    } else {
        //fmt.Println("Je kaarten zijn gelegd!")
        fmt.Println("")
        socket.Send("play", cards...)
    }
}

func parseCardInput(input string) ([]string, int) {
    var output []string
    var cards []int

    input = strings.TrimSpace(input)
    if len(input) == 0 {
        return output, 4
    }
    if input == "pas" || input == "pass" || input == "p" {
        if lastCards[1] == -1 {
            return output, 6
        }
        return output, 0
    }

    //
    var tokens []string
    for i := 0; i < len(input); i++ {
        c := input[i]
        switch (c) {
        case '3', '4', '5', '6', '7', '8', '9', 'J', 'j', 'Q', 'q', 'K', 'k', 'A', 'a', '2':
            tokens = append(tokens, string(c))
        case '1':
            i += 1
            if i == len(input) || input[i] != '0' {
                return output, 2
            }
            tokens = append(tokens, "10")

        case ' ', '\r', '\n', '\t':
            continue

        default:
            return output, 2
        }
    }
    //

    x := -1
    //tokens := strings.Split(input, " ")
    for _, token := range tokens {
        //token = strings.TrimSpace(token)
        if len(token) == 0 {
            continue
        }
        card := cardNameToIndex(token)
        if card < 3 {
            return output, 2
        }
        if x != -1 && card != x && card != 15 {
            return output, 5
        }
        if x == -1 && card != 15 {
            x = card
        }
        cards = append(cards, card)
    }

    cardCount := make(map[int]int)
    for i := 3; i <= 15; i++ {
        cardCount[i] = 0
    }
    for _, card := range currentCards {
        cardCount[card] += 1
    }
    for _, card := range cards {
        if cardCount[card] == 0 {
            return output, 1
        }
        cardCount[card] -= 1
    }

    if x == -1 {
        x = 14  // A
    }

    if len(cards) < lastCards[0] || x < lastCards[1] {
        return output, 3
    }

    sort.Ints(cards)
    for _, card := range cards {
        output = append(output, fmt.Sprint(card))
    }
    return output, 0
}

func onPassed(socket *messenger.Socket, args []string) {
    id, _ := strconv.Atoi(args[0])
    name := playerNames[id]
    name = padName(name)
    n, _ := strconv.Atoi(args[1])
    space := " "
    if len(args[1]) > 1 {
        space = ""
    }
    fmt.Println("(" + fmt.Sprint(n) + ")" + space + name + ":\tpas")
}

func onPlay(socket *messenger.Socket, args []string) {
    id, _ := strconv.Atoi(args[0])
    name := playerNames[id]
    name = padName(name)
    n, _ := strconv.Atoi(args[1])
    space := " "
    if len(args[1]) > 1 {
        space = ""
    }
    icons := []string{}
    lastPlay = []string{}
    for i := 2; i < len(args); i++ {
        lastPlay = append(lastPlay, args[i])
        index, _ := strconv.Atoi(args[i])
        icon := cardIndexToName(index)
        icons = append(icons, icon)
    }
    s := strings.Join(icons, " ")
    fmt.Println("(" + fmt.Sprint(n) + ")" + space + name + ":\t" + s)

    var cards []int
    for _, arg := range args[2:] {
        index, _ := strconv.Atoi(arg)
        cards = append(cards, index)
    }
    x := cards[0]
    if x == 15 {
        x = 14
    }
    lastCards[0] = len(cards)
    lastCards[1] = x
}

func onDone(socket *messenger.Socket, args []string) {
    id, _ := strconv.Atoi(args[0])
    name := playerNames[id]
    //status := args[1]
    if id != myId {
        fmt.Println("" + name + " is uit!")// Status: " + status)
    }
}

func onOver(socket *messenger.Socket, args []string) {
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

func onSweep(socket *messenger.Socket, args []string) {
    lastCards = []int{ 1, -1 }
    lastPlay = []string{}
    fmt.Println("\n-------------------------------------------------------------------------")
    fmt.Println("Iedereen heeft gepast, kaarten worden van tafel geveegd.. Volgende ronde!")
    fmt.Println("-------------------------------------------------------------------------\n")

    id, _ := strconv.Atoi(args[0])
    if id != myId {
        name := playerNames[id]
        //fmt.Println("Het is aan " + name + "...\n")
        fmt.Println(name + " mag beginnen...\n")
    }
}

func onExchange(socket *messenger.Socket, args []string) {
    badCard, _ := strconv.Atoi(args[0])
    bestCard, _ := strconv.Atoi(args[1])
    badPlayer, _ := strconv.Atoi(args[2])
    bestPlayer, _ := strconv.Atoi(args[3])
    badName := cardIndexToName(badCard)
    bestName := cardIndexToName(bestCard)
    shit := playerNames[badPlayer]
    pres := playerNames[bestPlayer]
    fmt.Println("\nDe president " + pres + " heeft met de shit " + shit + " een " + badName + " geruild voor een " + bestName + ".\n")
}

func onSettings(socket *messenger.Socket, args []string) {
    delayed = args[0] == "1"
}

func padName(name string) string {
    return strings.Repeat(" ", 16 - len(name)) + name
}

func cardIndexToName(index int) string {
    switch index {
    case 3:  return "3"
    case 4:  return "4"
    case 5:  return "5"
    case 6:  return "6"
    case 7:  return "7"
    case 8:  return "8"
    case 9:  return "9"
    case 10: return "10"
    case 11: return "J"
    case 12: return "Q"
    case 13: return "K"
    case 14: return "A"
    case 15: return "2"
    }
    return "?"
}

func cardNameToIndex(name string) int {
    switch name {
    case "3":  return 3
    case "4":  return 4
    case "5":  return 5
    case "6":  return 6
    case "7":  return 7
    case "8":  return 8
    case "9":  return 9
    case "10": return 10
    case "J", "j":  return 11
    case "Q", "q":  return 12
    case "K", "k":  return 13
    case "A", "a":  return 14
    case "2":  return 15
    }
    return -1
}
