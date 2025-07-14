package services

import (
	"errors"
	"fmt"
	"github/anansi-1/library_management/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book 
	ListBorrowedBooks(memberID int) []models.Book

}

type Library struct {

	Books map[int]models.Book
	Members map[int]models.Member

}

func (l *Library) AddBook(book models.Book) {

	l.Books[book.ID] = book
	

}
func (l *Library) RemoveBook(bookID int) {

	delete(l.Books,bookID)
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, ok := l.Books[bookID]
	if !ok {
		return errors.New("No such book")
	}

	member, ok := l.Members[memberID]
	if !ok {
		return errors.New("No such member")
	}

	if book.Status == "Borrowed" {
		return errors.New("The book is currently unavailable.")
	}

	book.Status = "Borrowed"
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member
l.Books[bookID] = book

	return nil
}


func (l *Library) ReturnBook(bookID int, memberID int) error {
	
	book, ok := l.Books[bookID]
    if !ok {
        return errors.New("No such book")
    }
    if book.Status == "Available" {
        return errors.New("the book is not borrowed")
    }
    member, ok := l.Members[memberID]
    if !ok {
        return errors.New("No such member")
    }
    book.Status = "Available"
    l.Books[bookID] = book
    newBorrowed := make([]models.Book, 0, len(member.BorrowedBooks))
for _, b := range member.BorrowedBooks {
    if b.ID != bookID {
        newBorrowed = append(newBorrowed, b)
    }
}
member.BorrowedBooks = newBorrowed

    l.Members[memberID] = member
    return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	var Books []models.Book
	for _, book := range l.Books {
		if book.Status == "Available" {
			Books = append(Books, book)
		}
	}
	return Books
}



func (l *Library) ListBorrowedBooks(memberID int) []models.Book {

	member,ok :=l.Members[memberID]
	if !ok {
		fmt.Println("Member doesn't exist")
		return nil
	}

	return member.BorrowedBooks

}

