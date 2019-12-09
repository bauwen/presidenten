
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
