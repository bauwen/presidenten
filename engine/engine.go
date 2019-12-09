package engine

// The engine for the Presidenten-game (cards 3, 4, ..., A, 2 have index 3, 1, ..., 14, 15)
// Note that this code does not in any way verify the correctness of the game play
// Validation of correct game play should be done externally

import (
    "math/rand"
    "sort"
    "time"
)

type Game struct {
    PlayerCount      int
    CurrentShit      int
    CurrentPresident int
    Cards            [][]int
    ExtraReceiver    int
    MatchWinners     []int
    MatchLosers      []int
}

func NewGame(playerCount int) *Game {
    return &Game{
        PlayerCount: playerCount,
        CurrentShit: -1,
        CurrentPresident: -1,
        Cards: make([][]int, playerCount),
    }
}

func (g *Game) StartMatch() int {
    deck := make([]int, 52)
    for i := 3; i <= 15; i++ {
        for j := 0; j < 4; j++ {
            deck[(i - 3) * 4 + j] = i
        }
    }
    
    shuffled := make([]int, 52)
    rand.Seed(time.Now().UnixNano())
    for i, j := range rand.Perm(52) {
        shuffled[i] = deck[j]
    }
    
    amount := 52 / g.PlayerCount
    cards := make([][]int, g.PlayerCount)
    for i := 0; i < g.PlayerCount; i++ {
        from := amount * i
        to := amount * (i + 1)
        for j := from; j < to; j++ {
            cards[i] = append(cards[i], shuffled[j])
        }
    }
    
    receiver := g.ExtraReceiver
    for i := amount * g.PlayerCount; i < 52; i++ {
        cards[receiver] = append(cards[receiver], shuffled[i])
        receiver = (receiver + 1) % g.PlayerCount
    }

    kCards := make([][]int, len(cards))
    for i := 0; i < len(cards); i++ {
        kParts := make([]int, len(cards[i]))
        for j := 0; j < len(cards[i]); j++ {
            kParts[j] = cards[i][j]
        }
        kCards[i] = kParts
    }
    
    for i := 0; i < len(kCards); i++ {
        sort.Ints(kCards[i])
    }
    
    receiver = g.ExtraReceiver
    g.ExtraReceiver = (g.ExtraReceiver + 1) % g.PlayerCount
    g.Cards = kCards
    
    return receiver
}

func (g *Game) ExchangeCards(index int) (int, int) {
    badCard := g.Cards[g.CurrentPresident][index]
    bestCard := g.Cards[g.CurrentShit][len(g.Cards[g.CurrentShit]) - 1]
    g.Cards[g.CurrentShit][len(g.Cards[g.CurrentShit]) - 1] = badCard
    g.Cards[g.CurrentPresident][index] = bestCard
    
    sort.Ints(g.Cards[g.CurrentShit])
    sort.Ints(g.Cards[g.CurrentPresident])
    
    return badCard, bestCard
}

func (g *Game) Play(player int, cards []int) (bool, bool, bool) {
    for _, value := range cards {
        for i := 0; i < len(g.Cards[player]); i++ {
            if g.Cards[player][i] == value {
                g.Cards[player] = append(g.Cards[player][:i], g.Cards[player][i+1:]...)
                break
            }
        }
    }
    
    first := false
    if len(g.Cards[player]) == 0 {
        first = true
        for i := 0; i < g.PlayerCount; i++ {
            if i != player && len(g.Cards[i]) == 0 {
                first = false
                break
            }
        }
        if first {
            g.CurrentPresident = player
            g.MatchWinners = append(g.MatchWinners, player)
        }
    }
    
    playersLeft := g.PlayerCount
    for i := 0; i < g.PlayerCount; i++ {
        if len(g.Cards[i]) == 0 {
            playersLeft -= 1
        }
    }
    if (playersLeft == 1) {
        for i := 0; i < g.PlayerCount; i++ {
            if len(g.Cards[i]) > 0 {
                g.CurrentShit = i
                g.MatchLosers = append(g.MatchLosers, i)
                break
            }
        }
    }
    
    gameOver := playersLeft <= 1
    playerDone := len(g.Cards[player]) == 0
    return gameOver, playerDone, first
}

func (g *Game) GetPresidentCount(player int) int {
    count := 0
    for _, i := range g.MatchWinners {
        if i == player {
            count += 1
        }
    }
    return count
}

func (g *Game) GetCivilianCount(player int) int {
    count := len(g.MatchWinners)
    count -= g.GetPresidentCount(player)
    count -= g.GetShitCount(player)
    return count
}

func (g *Game) GetShitCount(player int) int {
    count := 0
    for _, i := range g.MatchLosers {
        if i == player {
            count += 1
        }
    }
    return count
}
