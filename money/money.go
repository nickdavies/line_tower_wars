package money

import (
    "sync"
)

type PlayerBalance interface {
    Get() uint
    GetIncome() uint

    Spend(amount uint) bool

    IncreaseIncome(amount uint)
    DecreaseIncome(amount uint)

    PayIncome()
}

func NewPlayerBalance(balance, income, min_income uint) PlayerBalance {
    return &playerBalanceStruct{
        balance: balance,
        income: income,
        min_income: min_income,
    }
}

type playerBalanceStruct struct {
    sync.Mutex
    balance uint
    income uint
    min_income uint
}

func (b *playerBalanceStruct) Get() uint {
    b.Lock()
    defer b.Unlock()

    return b.balance
}

func (b *playerBalanceStruct) GetIncome() uint {
    b.Lock()
    defer b.Unlock()

    return b.income
}

func (b *playerBalanceStruct) Spend(amount uint) bool {
    b.Lock()
    defer b.Unlock()

    if b.balance > amount {
        b.balance -= amount
        return true
    }

    return false
}

func (b *playerBalanceStruct) IncreaseIncome(amount uint) {
    b.Lock()
    defer b.Unlock()

    b.income += amount
}

func (b *playerBalanceStruct) DecreaseIncome(amount uint) {
    b.Lock()
    defer b.Unlock()

    if amount > b.income {
        b.income = 0
    } else {
        b.income -= amount
    }

    if b.income < b.min_income {
        b.income = b.min_income
    }
}

func (b *playerBalanceStruct) PayIncome() {
    b.Lock()
    defer b.Unlock()

    b.balance += b.income
}

