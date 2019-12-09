package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strconv"
    
    "./server"
    "./client"
)

var defaultPort string = "3009"

func main() {
    if len(os.Args) < 2 {
        client.StartGame(defaultPort, false, "", "")
        return
    }
    
    switch os.Args[1] {
    case "server":
        port := defaultPort
        if len(os.Args) > 2 {
            port, err := strconv.Atoi(os.Args[2])
            if err != nil || port < 1 || 65535 < port {
                fmt.Println("Error: ongeldig poortnummer opgegeven")
                break
            }
        }
        server.StartServer(fmt.Sprint(port))
        
    case "bot":
        if len(os.Args) < 3 {
            fmt.Println("Error: de bestandsnaam van de bot is niet opgegeven")
            return
        }
        filename := os.Args[2]
        command := ""
        switch filepath.Ext(filename) {
        case ".py": command = "python"
        case ".js": command = "node"
        }
        if len(os.Args) > 3 {
            command = os.Args[3]
        }
        client.StartGame(defaultPort, true, command, filename)

    case "pybot":
        filename := "bot.py"
        if len(os.Args) >= 3 {
            filename = os.Args[2]
            if filepath.Ext(filename) != ".py" {
                filename += ".py"
            }
        }
        err := ioutil.WriteFile(filename, []byte(pybotCode()), 0644)
        if err != nil {
            fmt.Println("Error: het schrijven van het bestand '" + filename + "' is mislukt")
            return
        }
        fmt.Println("Een voorbeeld van een bot-programma is als '" + filename + "' aangemaakt")

    case "help":
        showHelp()

    default:
        fmt.Println("Error: onbekend commando '" + os.Args[1] + "' opgegeven")
        fmt.Println("Gebruik 'help' voor een lijst van ondersteunde commando's")
    }
}

func showHelp() {
    fmt.Println("Presidenten is een command-line versie van het gelijknamige kaartspel.")
    fmt.Println("\nGebruik:\n")
    fmt.Println("\tpresidenten commando [argumenten]\n")
    fmt.Println("De beschikbare combinaties zijn:")

    fmt.Println("\n\tpresidenten\n")
    fmt.Println("\t\tStart het spel (als mens) door te verbinden met een server.")
    fmt.Println("\t\tGeef na het starten het IP-adres op van de server, optioneel")
    fmt.Println("\t\tgevolgd door een poortnummer, bv. '192.168.0.100:8080'.")

    fmt.Println("\n\tpresidenten server [poort]\n")
    fmt.Println("\t\tStart een server voor een presidentenspel. Het poortnummer")
    fmt.Println("\t\tkan optioneel worden meegegeven (zo niet wordt 3009 gebruikt).")
    fmt.Println("\t\tOm vervolgens te starten met de spelers in de lobby, typ 'start'.")
    fmt.Println("\t\tIndien de spelers enkel bots zijn, kan ook 'run' worden getypt,")
    fmt.Println("\t\twat de artificiële vertragingen voor bot-beslissingen uitzet.")

    fmt.Println("\n\tpresidenten bot <bestandsnaam> [commando]\n")
    fmt.Println("\t\tStart het spel als bot. Dit betekent dat alle beslissingen")
    fmt.Println("\t\tvoor deze speler worden genomen door een computerprogramma.")
    fmt.Println("\t\tGeef de bestandsnaam op van het programma, optioneel gevolgd")
    fmt.Println("\t\tdoor het commando om dat programma uit te voeren. Voor een")
    fmt.Println("\t\tPython-programma kan dit bijvoorbeeld 'bot.py python' zijn.")
    fmt.Println("\t\tHet computerprogramma zal communiceren met het spel via zijn")
    fmt.Println("\t\tstandaardinvoer en standaarduitvoer, zie 'presidenten pybot'.")
    fmt.Println("\t\tHet IP-adres om te verbinden met een server en de spelersnaam")
    fmt.Println("\t\tmoeten bij de start wel eerst nog manueel worden ingegeven.")

    fmt.Println("\n\tpresidenten pybot [bestandsnaam]\n")
    fmt.Println("\t\tGenereert een voorbeeld van een computerprogramma dat kan")
    fmt.Println("\t\tworden gebruikt als functionele bot in een presidentenspel.")
    fmt.Println("\t\tOm het programma te kunnen gebruiken is Python 3 vereist.")
    fmt.Println("\t\tMerk echter op dat een bot in gelijk welke programmeertaal kan")
    fmt.Println("\t\tworden geschreven, zolang het protocol wordt geïmplementeerd.")

    fmt.Println("\n\tpresidenten help\n")
    fmt.Println("\t\tToon dit overzicht van de beschikbare commando's.")
}

func pybotCode() string {
    return `
"""
Voorbeeld van een simpele bot voor het presidentenspel.

Enkel de functies 'onExchange' en 'onTurn' moeten aangepast
worden om het programma te laten spelen zoals je wil.
Onderaan bevindt zich code voor het communicatieprotocol.

In beide functies zijn alle argumenten steeds een lijst
van natuurlijke getallen. De kaarten zelf worden voorgesteld
als indices van 3 tot 15 (met 3 = 3, ..., A = 14, 2 = 15).

Het is belangrijk dat het programma correct is. Indien je
kaarten kiest die je niet bezit, lagere kaarten oplegt, ...
dan zal de bot vastlopen tijdens het spelen van het spel!
"""

import sys

### Bepaal welke kaart je wil ruilen met de shit als president
# Argumenten:
#   cards: de kaarten die je momenteel in je hand hebt
# Return:
#   de kaart die je wilt ruilen (één getal)
def onExchange(cards):

    # geef de laagste kaart terug
    return cards[0]


### Bepaal welke kaarten je aflegt op jouw beurt
# Argumenten:
#   cards: de kaarten die je momenteel in je hand hebt
#   table: de kaarten op de tafel waarop jij moet spelen
#   hands: het aantal kaarten dat elke volgende speler bezit
#   history: alle kaarten die al zijn afgelegd geweest
# Return:
#   de kaarten die je wilt opleggen (lijst van getallen)
#   of None om te passen
def onTurn(cards, table, hands, history):

    # bepaal de rang van de huidige kaarten op tafel
    if len(table) == 0:
        index = -1
    else:
        index = table[0] if table[0] < 15 else 14

    # hou voor elke kaart bij hoeveel je er van hebt
    count = {}
    for i in range(3, 16):
        count[i] = 0
    for i in cards:
        count[i] += 1
    
    # leg alle kaarten af van de laagst mogelijke kaart
    for i in range(3, 16):
        c = count[i]
        if c > 0 and c >= len(table) and i >= index:
            return [i] * c

    # pas indien niet mogelijk (2's worden hier genegeerd)
    return None


# Je zou aan de code hieronder niets moeten wijzigen.

def cardIndexToName(index):
    if index == 11: return "J"
    if index == 12: return "Q"
    if index == 13: return "K"
    if index == 14: return "A"
    if index == 15: return "2"
    return str(index)

while True:
    parts = sys.stdin.readline().strip().split(" ")

    if parts[0] == "exchange":
        cards = list(map(int, parts[1:]))
        el = onExchange(cards)
        sys.stdout.write(cardIndexToName(el) + "\n")
        sys.stdout.flush()

    if parts[0] == "play":
        k = 1
        history = []
        while True:
            token = parts[k]
            if token == "#": break
            history.append(int(token))
            k += 1
        k += 1
        table = []
        while True:
            token = parts[k]
            if token == "#": break
            table.append(int(token))
            k += 1
        k += 1
        hands = []
        while True:
            token = parts[k]
            if token == "#": break
            hands.append(int(token))
            k += 1
        cards = list(map(int, parts[k+1:]))
        lst = onTurn(cards, table, hands, history)
        if lst is None:
            output = ["pas"]
        else:
            output = list(map(cardIndexToName, lst))
        sys.stdout.write(" ".join(output) + "\n")
        sys.stdout.flush()
`
}