package terrain

type Terrain int

const T_Grass Terrain = 0
const T_Wall Terrain = 1
const T_Spawn Terrain= 2
const T_Goal Terrain = 3

const T_NIL Terrain = 999

func (t Terrain) Buildable() bool {
    if t == T_Grass {
        return true
    }

    return false
}

