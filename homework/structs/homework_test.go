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

const healthOffset = 10

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.params1 |= uint32(health) << healthOffset
	}
}

const respectOffset = 20

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.params1 |= uint32(respect) << respectOffset
	}
}

const strengthOffset = 24

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.params1 |= uint32(strength) << strengthOffset
	}
}

const experienceOffset = 28

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.params1 |= uint32(experience) << experienceOffset
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.params2 |= uint16(level)
	}
}

const oneBitMask = 0b1
const houseOffset = 4

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.params2 |= oneBitMask << houseOffset
	}
}

const gunOffset = 5

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.params2 |= oneBitMask << gunOffset
	}
}

const familyOffset = 6

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.params2 |= oneBitMask << familyOffset
	}
}

const personTypeOffset = 7

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.params2 |= uint16(personType) << personTypeOffset
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
	return unsafe.String(unsafe.SliceData(p.name[:n]), n)
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

const tenBitsMask = 0b1111111111

func (p *GamePerson) Mana() int {
	manaBitMask := uint32(tenBitsMask)
	return int(p.params1 & manaBitMask)
}

func (p *GamePerson) Health() int {
	healthBitMask := uint32(tenBitsMask << healthOffset)
	return int(p.params1 & healthBitMask >> healthOffset)
}

const fourBitsMask = 0b1111

func (p *GamePerson) Respect() int {
	respectBitMask := uint32(fourBitsMask << respectOffset)
	return int(p.params1 & respectBitMask >> respectOffset)
}

func (p *GamePerson) Strength() int {
	strengthBitMask := uint32(fourBitsMask << strengthOffset)
	return int(p.params1 & strengthBitMask >> strengthOffset)
}

func (p *GamePerson) Experience() int {
	expereinceBitMask := uint32(fourBitsMask << experienceOffset)
	return int(p.params1 & expereinceBitMask >> experienceOffset)
}

func (p *GamePerson) Level() int {
	levelBitMask := uint16(fourBitsMask)
	return int(p.params2 & levelBitMask)
}

func (p *GamePerson) HasHouse() bool {
	const houseMask = uint16(oneBitMask << houseOffset)
	return p.params2&houseMask == houseMask
}

func (p *GamePerson) HasGun() bool {
	const gunMask = uint16(oneBitMask << gunOffset)
	return p.params2&gunMask == gunMask
}

func (p *GamePerson) HasFamily() bool {
	const familyMask = uint16(oneBitMask << familyOffset)
	return p.params2&familyMask == familyMask
}

const twoBitsMask = 0b11

func (p *GamePerson) Type() int {
	const personTypeMask = uint16(twoBitsMask << personTypeOffset)
	return int(p.params2 & personTypeMask >> personTypeOffset)
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
