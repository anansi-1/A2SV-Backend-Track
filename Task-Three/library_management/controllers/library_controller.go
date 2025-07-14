package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"github/anansi-1/library_management/models"
	"github/anansi-1/library_management/services"
)

type LibraryController struct {
	service *services.Library
}

func NewLibraryController(service *services.Library) *LibraryController {
	return &LibraryController{service: service}
}

func (c *LibraryController) Run() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nLibrary System")
		fmt.Println("1.Add Book")
		fmt.Println("2.Borrow Book")
		fmt.Println("3.Return Book")
		fmt.Println("4.List Available Books")
		fmt.Println("5.List Borrowed Books")
		fmt.Println("6.Exit")
		fmt.Print("Enter choice: ")

		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			c.handleAddBook(scanner)
		case "2":
			c.handleBorrowBook(scanner)
		case "3":
			c.handleReturnBook(scanner)
		case "4":
			c.handleListAvailableBooks()
		case "5":
			c.handleListBorrowedBooks(scanner)
		case "6":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice, try again.")
		}
	}
}

func (c *LibraryController) handleAddBook(scanner *bufio.Scanner) {
	fmt.Print("Enter Book ID: ")
	scanner.Scan()
	idStr := scanner.Text()
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}

	fmt.Print("Enter Title: ")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter Author: ")
	scanner.Scan()
	author := strings.TrimSpace(scanner.Text())

	book := models.Book{
		ID:     id,
		Title:  title,
		Author: author,
		Status: "Available",
	}

	c.service.AddBook(book)
	fmt.Println("Book added successfully!")
}

func (c *LibraryController) handleBorrowBook(scanner *bufio.Scanner) {
	bookID, memberID, ok := c.readBookAndMemberIDs(scanner)
	if !ok {
		return
	}

	err := c.service.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error borrowing book:", err)
	} else {
		fmt.Println("Book borrowed successfully!")
	}

}

func (c *LibraryController) handleReturnBook(scanner *bufio.Scanner) {
	bookID, memberID, ok := c.readBookAndMemberIDs(scanner)
	if !ok {
		return
	}

	err := c.service.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error returning book:", err)
	} else {
		fmt.Println("Book returned successfully!")
	}
}

func (c *LibraryController) handleListAvailableBooks() {
	books := c.service.ListAvailableBooks()
	if len(books) == 0 {
		fmt.Println("No available books.")
		return
	}

	fmt.Println("Available Books:")
	for _, b := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", b.ID, b.Title, b.Author)
	}
}

func (c *LibraryController) handleListBorrowedBooks(scanner *bufio.Scanner) {
	fmt.Print("Enter Member ID: ")
	scanner.Scan()
	memberIDStr := strings.TrimSpace(scanner.Text())
	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		fmt.Println("Invalid Member ID")
		return
	}

	books := c.service.ListBorrowedBooks(memberID)
	if len(books) == 0 {
		fmt.Println("No borrowed books or member not found.")
		return
	}

	fmt.Printf("Books borrowed by member %d:\n", memberID)
	for _, b := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", b.ID, b.Title, b.Author)
	}
}

func (c *LibraryController) readBookAndMemberIDs(scanner *bufio.Scanner) (bookID int, memberID int, ok bool) {
	fmt.Print("Enter Book ID: ")
	scanner.Scan()
	bookIDStr := strings.TrimSpace(scanner.Text())
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		fmt.Println("Invalid Book ID")
		return 0, 0, false
	}

	fmt.Print("Enter Member ID: ")
	scanner.Scan()
	memberIDStr := strings.TrimSpace(scanner.Text())
	memberID, err = strconv.Atoi(memberIDStr)
	if err != nil {
		fmt.Println("Invalid Member ID")
		return 0, 0, false
	}
	return bookID, memberID, true
}
