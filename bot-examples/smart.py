import sys
from functools import reduce

### Decide what card to exchange with the shit when you are president
# cards: your own cards, as indices from 3 to 15
def onExchange(cards):
    count = {}
    for i in range(3, 16):
        count[i] = 0
    for i in cards:
        count[i] += 1
    
    for k in range(1, 5):
        for i in range(3, 16):
            if count[i] == k:
                return i
         
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
        index = table[0] if table[0] < 15 else 14
    count = {}
    for i in range(3, 16):
        count[i] = 0
    for i in cards:
        count[i] += 1
    jokers = count[15]

    kinds = 0
    for i in range(3, 15):
        if count[i] > 0:
            kinds += 1

    # if only one kind and two's: let's try to play one big hand!
    if kinds == 1:
        for i in range(3, 15):
            c = count[i]
            if c > 0 and i >= index and c + jokers >= len(table):
                return [i] * c + [15] * jokers

    # if two kinds and two's: try to win with big hand to finish afterwards
    if kinds == 2:
        for i in range(14, 2, -1):
            c = count[i]
            if c > 0 and i >= index and c + jokers >= max(4, len(table)):
                return [i] * c + [15] * jokers

    # if in opening of the game: just try to get rid of bad cards
    start_hand = 52 / (len(hands) + 1)
    curr_hand = reduce(lambda x,y: x+y, hands) / len(hands)
    if curr_hand > start_hand/2:
        for k in range(max(1, len(table)), 5):
            for i in range(3, 15):
                c = count[i]
                if c == k and i >= index:
                    return [i] * c

        return None

    # if next opponent is "vulnerable": try to block him
    if hands[0] < 4:
        for j in range(0, jokers + 1):
            for k in range(4, max(hands[0], len(table) - 1), -1):
                for i in range(14, 2, -1):
                    c = count[i]
                    diff = k - c
                    if diff == j and i >= index:
                        return [i] * c + [15] * diff

    # if in the middle/end of the game: no idea what to do
    # although the strategy here should be pretty important
    for k in range(max(1, len(table)), 5):
        for i in range(3, 15):
            c = count[i]
            if c == k and i >= index:
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
