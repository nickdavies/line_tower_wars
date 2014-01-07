package tower

type TowerType struct {
    Id int
    Name string

    Range float64
    FireRate int
    Damage int

    Health uint

    Cost uint
}

type Tower struct {
    Type *TowerType
}

func (t *Tower) Update(deltaTime int64) {
}
