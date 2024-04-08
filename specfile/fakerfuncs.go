package specfile

import (
	"github.com/brianvoe/gofakeit/v7"
)

type fakerFuncs struct {
	faker *gofakeit.Faker
}

func (f *fakerFuncs) New(seed uint64) *fakerFuncs {
	return &fakerFuncs{
		faker: gofakeit.New(seed),
	}
}

func (f *fakerFuncs) UUID() string {
	return f.faker.UUID()
}

func (f *fakerFuncs) Name() string {
	return f.faker.Name()
}

func (f *fakerFuncs) FirstName() string {
	return f.faker.FirstName()
}

func (f *fakerFuncs) LastName() string {
	return f.faker.LastName()
}

func (f *fakerFuncs) Email() string {
	return f.faker.Email()
}

func (f *fakerFuncs) GamerTag() string {
	return f.faker.Gamertag()
}
