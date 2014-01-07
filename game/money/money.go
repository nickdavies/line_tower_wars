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

    IncomeInterval() int64
}

func NewPlayerBalance(balance, income, min_income uint, income_interval int64) PlayerBalance {
    return &playerBalanceStruct{
        balance: balance,
        income: income,
        min_income: min_income,
        income_interval: income_interval,
    }
}

type playerBalanceStruct struct {
    sync.Mutex
    balance uint
    income uint
    min_income uint
    income_interval int64
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

    if b.balance >= amount {
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

func (b *playerBalanceStruct) IncomeInterval() int64 {
    return b.income_interval
}

