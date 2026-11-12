package main

import (
	"os"

	"books"
)

func main() {
	// Make a brand new empty message.  A Message allocates Cap'n Proto structs.
	msg, seg := shark.NewSingleSegmentMessage(nil)

	// Create a new Book struct.  Every message must have a root struct.
	book, err := books.NewRootBook(seg)
	if err != nil {
		panic(err)
	}
	_ = book.SetTitle("War and Peace")
	book.SetPageCount(1440)

	// Write the message to stdout.
	err = shark.NewEncoder(os.Stdout).Encode(msg)
	if err != nil {
		panic(err)
	}
}
