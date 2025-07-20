package main


import(
	"github/anansi-1/Task-Four-Task-Manager/router"
	"log"
)

func main() {
		if err := router.Run(); err != nil {
		log.Fatal(err)
	}
	
}