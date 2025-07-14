# Library Management System 

## Overview

This is a console-based Library Management System written in Go. It demonstrates core Go programming concepts such as **structs**, **interfaces**, **methods**, **slices**, and **maps**. Users can interact with the system through the terminal to manage books and library members.

---

## Features

- Add a new book  
- Remove an existing book  
- Borrow a book  
- Return a book  
- List all available books  
- View borrowed books by a member  

---

## Data Structures

### Book and Members

Book: Represents a book in the library.

Member: Represents a library member.


    ```go
    type Book struct {
        ID     int
        Title  string
        Author string
        Status string 
    }
    type Member struct {
        ID            int
        Name          string
        BorrowedBooks []Book
    }

## Running the Application

### 1. Navigate to the `library_management` directory:

   ```bash
   cd library_management

### 2. Command-Line Menu

Once the program starts, you'll see the following menu:

1.Library System

2.Add Book

3.Borrow Book

4.Return Book

5.List Available Books

6.List Borrowed Books

7.Exit
Enter choice:


### 3. Select an Action

Type the number corresponding to the action you'd like to perform.

