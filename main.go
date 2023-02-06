package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {
	var region string = "Chinese"
	var userCategory string
	categories := map[string]bool{
		"green":  true,
		"oolong": true,
		"black":  true,
		"puerh":  true,
		"white":  true,
		"yellow": true,
	}

	type Tea struct {
		Category    string
		SubCategory string
		Name        string
		Province    string
		Rating      float64
	}

	var teasSlice []Tea
	var teaSliceByCategory []Tea

	fmt.Printf("Welcome to the %v tea picker application\n", region)
	fmt.Println("Provide a category of tea and you will receive a random tea with brewing instructions!")

	data, _ := ioutil.ReadFile("./teas.json")
	json.Unmarshal(data, &teasSlice)
	fmt.Printf("teasSlice %v \n", teasSlice)

	for {

		fmt.Println("Enter your preferred tea category")
		fmt.Scan(&userCategory)

		_, ok := categories[userCategory]
		if !ok {
			fmt.Printf("%v is not a valid category, options are %v\n", userCategory, categories)
			continue
		}
		for _, teaData := range teasSlice {
			fmt.Printf("teaData %v \n", teaData)
			if teaData.Category == userCategory {
				teaSliceByCategory = append(teaSliceByCategory, teaData)
			}

		}
		fmt.Printf("teaSliceByCategory %v\n", teaSliceByCategory)

		break
	}

}
