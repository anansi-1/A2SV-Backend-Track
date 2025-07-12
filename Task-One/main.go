package main

import (
	"fmt"
)

func main() {
	
	var userName string
	var numberOfSubjects int

	fmt.Println("Enter your name and number of subjects you have taken: \n eg. anansi 3")
	fmt.Scanf("%s %d",&userName,&numberOfSubjects)

	
	subjectGrade := make(map[string]int)
	// subjects := 0
	fmt.Printf("Enter the name of the Subject and your Grade %d times for each subject: \n eg. Math 100 \n",numberOfSubjects)

	i := 0
	for i < numberOfSubjects {
		var subject string
		var grade int
			fmt.Scanf("%s %d",&subject,&grade) 
			if grade < 0 || grade > 100 {
				fmt.Printf("%d is not a valid grade \n Enter a vaild grade again between 0-100 \n",grade)
				continue
			}
			subjectGrade[subject] = grade
			i++
	}
	fmt.Printf("Student Name: %s \n",userName)
	calculateAverage(subjectGrade)

}

func calculateAverage(grades map[string]int){
	sum := 0
	fmt.Println("Student Result:")
	for subject,grade := range grades{
		sum += grade
		fmt.Printf("%s : %d \n",subject,grade)
	}
	average := float64(sum)/float64(len(grades))
	fmt.Printf("Average Grade:%.2f \n",average)
}