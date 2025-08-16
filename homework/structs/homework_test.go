package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(p *GamePerson) {
		copy(p.name[:], name)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.params1 |= uint32(mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.params1 |= uint32(health) << 10
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.params1 |= uint32(respect) << 20
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.params1 |= uint32(strength) << 24
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.params1 |= uint32(experience) << 28
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.params2 |= uint16(level)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.params2 |= 0b1_0000
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.params2 |= 0b10_0000
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.params2 |= 0b100_0000
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.params2 |= uint16(personType) << 8
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x    int32
	y    int32
	z    int32
	gold uint32

	/*
		params encoding bits order:
		mana 		[0; 1000] 10 bits
		health 		[0; 1000] 10 bits
		respect		[0; 10] 4 bits
		strength	[0; 10] 4 bits
		expereince 	[0; 10] 4 bits
		total 32 bits
	*/
	params1 uint32

	/*
		level		[0; 10] 4 bits
		house 		flag	1 bit
		weapon		flag	1 bit
		family		flag	1 bit
		type		[0, 1, 2] -> [builder, farrier, warrior] 2 bits
		total 9 bits, other 7 are unused
	*/
	params2 uint16
	name    [42]byte
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}

	for _, option := range options {
		option(&person)
	}

	return person
}

func (p *GamePerson) Name() string {
	n := 0
	for n < len(p.name) && p.name[n] != 0 {
		n++
	}
	return string(p.name[:n])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	manaBitMask := uint32(0b1111111111)
	return int(p.params1 & manaBitMask)
}

func (p *GamePerson) Health() int {
	healthBitMask := uint32(0b1111111111 << 10)
	return int(p.params1 & healthBitMask >> 10)
}

func (p *GamePerson) Respect() int {
	respectBitMask := uint32(0b1111 << 20)
	return int(p.params1 & respectBitMask >> 20)
}

func (p *GamePerson) Strength() int {
	strengthBitMask := uint32(0b1111 << 24)
	return int(p.params1 & strengthBitMask >> 24)
}

func (p *GamePerson) Experience() int {
	expereinceBitMask := uint32(0b1111 << 24)
	return int(p.params1 & expereinceBitMask >> 24)
}

func (p *GamePerson) Level() int {
	levelBitMask := uint16(0b1111)
	return int(p.params2 & levelBitMask)
}

func (p *GamePerson) HasHouse() bool {
	return p.params2&0b1_0000 == 0b1_0000
}

func (p *GamePerson) HasGun() bool {
	return p.params2&0b10_0000 == 0b10_0000
}

func (p *GamePerson) HasFamily() bool {
	return p.params2&0b100_0000 == 0b100_0000
}

func (p *GamePerson) Type() int {
	return int(p.params2 & (0b11 << 7))
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
