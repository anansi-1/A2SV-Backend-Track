package main

import (
	"github/anansi-1/library_management/controllers"
	"github/anansi-1/library_management/services"
	"github/anansi-1/library_management/models"
)

func main() {
	library := &services.Library{
		Books:   make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
	library.Members[1] = models.Member{ID: 1, Name: "a"}
	library.Members[2] = models.Member{ID: 2, Name: "b"}

	controller := controllers.NewLibraryController(library)
	controller.Run()
}
