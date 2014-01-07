package pathing

import (
    "math"
)

import (
    "github.com/nickdavies/go-astar/astar"
)

type Loc struct {
    Row uint16
    Col uint16
}

func (l Loc) ToFloat(offset float64) Locf {
    return Locf{
        Row: float64(l.Row) + offset,
        Col: float64(l.Col) + offset,
    }
}

type Locf struct {
    Row float64
    Col float64
}

func (l Locf) ToInt() Loc {
    return Loc{
        Row: uint16(l.Row),
        Col: uint16(l.Col),
    }
}

func (l Locf) Dist(other Locf) float64 {
    return math.Sqrt(math.Pow(l.Row - other.Row, 2) + math.Pow(l.Row - other.Row, 2))
}

type Path struct {
    aStar astar.AStar
    root *astar.PathPoint

    astar_path []astar.Point

    row_off uint16
    col_off uint16

    path []Loc
    path_map map[Loc]int
}

func NewPath(aStar astar.AStar, root *astar.PathPoint, row_off, col_off uint16) *Path {

    p := &Path{
        aStar: aStar,
        root: root,
        row_off: row_off,
        col_off: col_off,

        path: make([]Loc, 0),
        path_map: make(map[Loc]int),
    }

    i := 0
    for root != nil {
        point := Loc{uint16(root.Row) + row_off, uint16(root.Col) + col_off}

        p.path = append(p.path, point)
        p.astar_path = append(p.astar_path, root.Point)
        p.path_map[point] = i

        root = root.Parent
        i++
    }

    return p
}
func (p *Path) RouteTo(source Loc) *Path {

    astar_source := []astar.Point{
        astar.Point{
            Row: int(source.Row - p.row_off),
            Col: int(source.Col - p.col_off),
        },
    }

    path := p.aStar.FindPath(astar.NewListToPoint(true), p.astar_path, astar_source)

    return NewPath(p.aStar, path, p.row_off, p.col_off)
}

func (p *Path) Move(currentf Locf, distance float64) (Locf, bool) {
    current := currentf.ToInt()
    _, ok := p.path_map[current]
    if !ok {
        panic("You are not on that path!")
    }

    for {
        current_index := p.path_map[current]
        if current_index == len(p.path) - 1 {
            return currentf, true
        }
        next := p.path[current_index + 1]
        nextf := next.ToFloat(0.5)

        next_dist := math.Sqrt(math.Pow(nextf.Row - currentf.Row, 2) + math.Pow(nextf.Col - currentf.Col, 2))

        if next_dist <= distance {
            current = next
            currentf = nextf

            distance -= next_dist
        } else {
            dist_percent := distance / next_dist

            currentf.Row += dist_percent * (nextf.Row - currentf.Row)
            currentf.Col += dist_percent * (nextf.Col - currentf.Col)

            return currentf, false
        }
    }
}

func (p *Path) Start() Loc {
    return p.path[0]
}

func (p *Path) End() Loc {
    return p.path[len(p.path)-1]
}

func (p *Path) Startf() Locf {
    return p.path[0].ToFloat(0.5)
}

func (p *Path) Endf() Locf {
    return p.path[len(p.path)-1].ToFloat(0.5)
}

func (p *Path) On(l Loc) bool {
    _, ok := p.path_map[l]
    return ok
}

func (p *Path) GetPathArray() []Loc {
    return p.path
}

