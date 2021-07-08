package main

import (
	"fmt"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func writeFile(value []byte, language string){
	err := ioutil.WriteFile(".gitignore", value, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully created gitignore for %s\n", language)
}



func main() {
	loadingSpinner := wow.New(os.Stdout, spin.Get(spin.Dots), "Loading")
	loadingSpinner.Start()
	resp, err := http.Get("https://www.gitignore.io/api/list")
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	items := strings.Split(sb,",")
	loadingSpinner.Stop()

	prompt := promptui.Select{
		Label: "Select environment",
		Items: items,
		StartInSearchMode : true,
		Searcher : func(input string, index int) bool {
		pepper := items[index]
		name := strings.Replace(strings.ToLower(pepper), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	addingSpinner := wow.New(os.Stdout, spin.Get(spin.Dots), "Adding gitignore")

	giResponse, err := http.Get(fmt.Sprintf("https://www.gitignore.io/api/%s",result))
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	giBody, err := ioutil.ReadAll(giResponse.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	addingSpinner.Stop()
	writeFile(giBody,result)


}
