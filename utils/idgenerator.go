package utils

type IdGenerator interface {
	NextId() uint64
}

type idGeneratorImpl struct {
	ch chan uint64
}

func (en *idGeneratorImpl) NextId() uint64 {
	return <- en.ch
}

func NewIdGenerator(startWith uint64) IdGenerator {
	en := &idGeneratorImpl{make(chan uint64)}
	go func() {
		for {
			en.ch <- startWith
			startWith++
		}
	}()

	return en
}
