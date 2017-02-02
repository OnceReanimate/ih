package srv

import (
)

type books struct {
	book	chan uid
	unbook	chan uid
	conns	*conns
	waits	[4]uid
	wait	int
}

func newBooks(conns *conns) *books {
	books := new(books)

	books.book = make(chan uid)
	books.unbook = make(chan uid)
	books.conns = conns
	books.wait = 0

	return books;
}

func (books *books) Loop() {
	for {
		select {
		case uid := <-books.book:
			books.add(uid)
		case uid := <-books.unbook:
			books.sub(uid)
		}
	}
}

func (books *books) Book() chan<- uid {
	return books.book
}

func (books *books) Unbook() chan<- uid {
	return books.unbook
}

func (books *books) add(uid uid) {
	for i := 0; i < books.wait; i++ {
		if books.waits[i] == uid {
			return
		}
	}

	books.waits[books.wait] = uid;
	books.wait++
	if books.wait == 4 {
		books.conns.start <- books.waits
		books.wait = 0
	}
}

func (books *books) sub(uid uid) {
	i := 0
	for i < books.wait && books.waits[i] != uid {
		i++
	}

	if i == books.wait {
		return
	}

	// swap to back, and pop back
	e := books.wait - 1;
	books.waits[i], books.waits[e] = books.waits[e], books.waits[i];
	books.wait--
}



