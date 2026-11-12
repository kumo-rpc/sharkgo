package demo_test

import (
	"fmt"
	"io"

	"github.com/kumo-rpc/sharkgo"
	"github.com/kumo-rpc/sharkgo/internal/demo/books"
)

func Example_book() {
	r, w := io.Pipe()
	go writer(w)
	reader(r)
	// Output:
	// "War and Peace" has 1440 pages
}

func writer(out io.Writer) {
	// Make a brand new empty message.  A Message allocates Shark Proto structs.
	msg, seg := shark.NewSingleSegmentMessage(nil)

	// Create a new Book struct.  Every message must have a root struct.
	book, err := books.NewRootBook(seg)
	if err != nil {
		panic(err)
	}
	book.SetTitle("War and Peace")
	book.SetPageCount(1440)

	// Write the message to stdout.
	err = shark.NewEncoder(out).Encode(msg)
	if err != nil {
		panic(err)
	}
}

func reader(in io.Reader) {
	// Read the message from stdin.
	msg, err := shark.NewDecoder(in).Decode()
	if err != nil {
		panic(err)
	}

	// Extract the root struct from the message.
	book, err := books.ReadRootBook(msg)
	if err != nil {
		panic(err)
	}

	// Access fields from the struct.
	title, err := book.Title()
	if err != nil {
		panic(err)
	}
	pageCount := book.PageCount()
	fmt.Printf("%q has %d pages\n", title, pageCount)
}
