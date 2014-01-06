package game

import (
    "time"
    "sync"
)

import (
    "github.com/nickdavies/line_tower_wars/game/stage"
    "github.com/nickdavies/line_tower_wars/game/player"
)

type PlayerControls interface {
    Tick() int64
}

type playerController struct {
    game *Game
}

func (p *playerController) Tick() (int64, bool) {
    p.game.resetBarrier.Done()
    p.game.resetBarrier.Wait()

    p.game.tickBarrier.Done()
    p.game.tickBarrier.Wait()

    return p.game.deltaTime, p.game.running
}

type Game struct {
    NumPlayers int
    running bool
    end bool
    deltaTime int64

    resetBarrier sync.WaitGroup
    tickBarrier sync.WaitGroup

    players []*player.Player
    stage *stage.Stage
}

func NewGame(cfg GameConfig, NumPlayers int) *Game {
    if NumPlayers != 2 {
        panic("Only two players are supported atm")
    }

    s := stage.NewStage(NumPlayers)

    var prev_player *player.Player

    players := make([]*player.Player, NumPlayers)
    for i := 0; i < NumPlayers; i++ {
        players[i] = player.NewPlayer(s.GetPlayer(i))

        if prev_player != nil {
            prev_player.NextPlayer = players[i]
        }
        prev_player = players[i]
    }
    prev_player.NextPlayer = players[0]

    g := &Game{
        NumPlayers: NumPlayers,
        running: true,
        deltaTime: 0,

        players: players,
        stage: s,
    }

    g.resetBarrier.Add(NumPlayers + 1)
    g.tickBarrier.Add(NumPlayers + 1)

    return g
}

func (g *Game) GetStage() *stage.Stage {
    return g.stage
}

func (g *Game) GetPlayer(id int) *player.Player {
    return g.players[id]
}

func (g *Game) Running() bool {
    return g.running
}

func (g *Game) End() {
    g.end = true
}

func (g *Game) Run() int {
    var prevTime int64
    var winner int

    var update_barrier_1 sync.WaitGroup
    var update_barrier_2 sync.WaitGroup

    update_barrier_1.Add(1)
    update_barrier_2.Add(1)

    for i := 0; i < g.NumPlayers; i++ {
        player := g.players[i]
        go func() {
            for g.running {
                update_barrier_1.Done()
                update_barrier_1.Wait()

                if !g.running {
                    return
                }

                player.Update(g.deltaTime)

                update_barrier_2.Done()
                update_barrier_2.Wait()
            }
        }()
    }

    for g.running {
        // release reset barrier and reset it
        // while people are waiting for tick barrier
        g.resetBarrier.Done()
        g.resetBarrier.Wait()
        g.resetBarrier.Add(g.NumPlayers + 1)

        // Everyone else is now waiting for the tick barrier

        now := time.Now().UnixNano()
        g.deltaTime = now - prevTime
        prevTime = now

        // Check to see if anyone has won
        living_players := 0
        living := 0
        for i := 0; i < g.NumPlayers; i++ {
            if g.players[i].GetLives() <= 0 {
                g.players[i].Die()
            } else {
                living_players++
                living = i
            }
        }

        if living_players == 1 {
            winner = living
            g.running = false
        }

        if g.end {
            winner = -1
            g.running = false
        }


        // let the player structs update
        update_barrier_1.Done()
        update_barrier_1.Wait()
        update_barrier_1.Add(g.NumPlayers + 1)

        // wait for player structs to finish updating
        update_barrier_2.Done()
        update_barrier_2.Wait()
        update_barrier_2.Add(g.NumPlayers + 1)

        // release tick barrier and reset it
        // while people are waiting for reset barrier
        g.tickBarrier.Done()
        g.tickBarrier.Wait()
        g.tickBarrier.Add(g.NumPlayers + 1)
    }

    return winner
}
