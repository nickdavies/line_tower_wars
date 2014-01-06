package game

import (
    "runtime"
    "fmt"
    "time"
    "sync"
)

import (
    "github.com/nickdavies/line_tower_wars/game/stage"
    "github.com/nickdavies/line_tower_wars/game/player"
)

type PlayerControls interface {
    Tick() (int64, bool)
}

type playerController struct {
    game *Game

    my_player *player.Player
}

func (p *playerController) Tick() (int64, bool) {
    p.game.resetBarrier.Done()
    p.game.resetBarrier.Wait()

    p.game.tickBarrier.Done()
    p.game.tickBarrier.Wait()

    return p.game.deltaTime, p.game.running
}

type Game struct {
    sync.Mutex

    NumPlayers int
    running bool
    end bool
    deltaTime int64

    resetBarrier sync.WaitGroup
    tickBarrier sync.WaitGroup

    players []*player.Player
    stage *stage.Stage
}

func NewGame(cfg GameConfig, NumPlayers int) (*Game, []PlayerControls) {
    if NumPlayers != 2 {
        panic("Only two players are supported atm")
    }

    g := &Game{
        NumPlayers: NumPlayers,
        running: true,
    }

    g.stage = stage.NewStage(NumPlayers)

    var prev_player *player.Player

    g.players = make([]*player.Player, NumPlayers)
    controls := make([]PlayerControls, NumPlayers)
    for i := 0; i < NumPlayers; i++ {
        g.players[i] = player.NewPlayer(g.stage.GetPlayer(i))
        controls[i] = &playerController{
            game: g,
            my_player: g.players[i],
        }

        if prev_player != nil {
            prev_player.NextPlayer = g.players[i]
            g.players[i].PrevPlayer = prev_player
        }
        prev_player = g.players[i]
    }
    prev_player.NextPlayer = g.players[0]
    g.players[0].PrevPlayer = g.players[NumPlayers - 1]

    g.resetBarrier.Add(NumPlayers + 1)
    g.tickBarrier.Add(NumPlayers + 1)

    return g, controls
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

    update_barrier_1.Add(g.NumPlayers + 1)
    update_barrier_2.Add(g.NumPlayers + 1)

    for i := 0; i < g.NumPlayers; i++ {
        go func(p_id int) {
            player := g.players[p_id]
            for g.running {
                update_barrier_1.Done()
                update_barrier_1.Wait()

                fmt.Printf("player %d update tick\n", p_id)

                player.Update(g.deltaTime)

                update_barrier_2.Done()
                update_barrier_2.Wait()
            }
        }(i)
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

        // let the player structs update
        update_barrier_1.Done()
        update_barrier_1.Wait()
        update_barrier_1.Add(g.NumPlayers + 1)

        if living_players == 1 {
            winner = living
            g.running = false
        }

        if g.end {
            winner = -1
            g.running = false
        }

        // wait for player structs to finish updating
        update_barrier_2.Done()
        update_barrier_2.Wait()
        update_barrier_2.Add(g.NumPlayers + 1)

        // release tick barrier and reset it
        // while people are waiting for reset barrier
        g.tickBarrier.Done()
        g.tickBarrier.Wait()
        g.tickBarrier.Add(g.NumPlayers + 1)

        runtime.Gosched()

        fmt.Println("game tick")
    }

    return winner
}
