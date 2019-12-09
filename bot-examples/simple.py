import sys

### Decide what card to exchange with the shit when you are president
# cards: your own cards, as indices from 3 to 15
def onExchange(cards):
    return cards[0]


### Decide what cards to play when it is your turn
# cards: your own cards, as indices from 3 to 15
# table: the cards on the table you need to beat
# hands: the number of cards of each subsequent player
# history: all the cards that have already been played
def onTurn(cards, table, hands, history):
    if len(table) == 0:
        index = -1
    else:
        index = table[0]
        if index == 15:
            index == 14

    count = {}
    for i in range(3, 16):
        count[i] = 0
    for i in cards:
        count[i] += 1
    
    for i in range(3, 16):
        c = count[i]
        if c > 0 and i >= index and c >= len(table):
            return [i] * c

    return None


# You should not have to alter anything below
# It is just an abstraction for the interface

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
