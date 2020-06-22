package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

const (
	port = 19999
)

func main() {
	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}

	//wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	//if err != nil {
	//	panic(err)
	//}
	//// test baidu.com
	//if err := wd.Get("http://www.baidu.com"); err != nil {
	//	panic(err)
	//}
	//wd.Quit()

	// test proxy
	//caps.AddProxy(selenium.Proxy{
	//	Type: selenium.Manual,
	//	HTTP: "http://127.0.0.1:108089",
	//	//SOCKS:        "socks5://127.0.0.1:108088",
	//	//SOCKSVersion: 5,
	//	HTTPPort: 108089,
	//	//SocksPort: 108088,
	//})

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}

	// Navigate to the simple playground interface.
	if err := wd.Get("https://google.com"); err != nil {
		panic(err)
	}

	//Get a reference to the text box containing code.
	elem, err := wd.FindElement(selenium.ByCSSSelector, "#code")
	if err != nil {
		panic(err)
	}
	// Remove the boilerplate code already in the text box.
	if err := elem.Clear(); err != nil {
		panic(err)
	}

	// Enter some new code in text box.
	err = elem.SendKeys(`
		package main
		import "fmt"
	
		func main() {
			fmt.Println("Hello WebDriver!\n")
		}
	`)
	if err != nil {
		panic(err)
	}

	// Click the run button.
	btn, err := wd.FindElement(selenium.ByCSSSelector, "#run")
	if err != nil {
		panic(err)
	}
	if err := btn.Click(); err != nil {
		panic(err)
	}

	// Wait for the program to finish running and get the output.
	outputDiv, err := wd.FindElement(selenium.ByCSSSelector, "#output")
	if err != nil {
		panic(err)
	}

	var output string
	for {
		output, err = outputDiv.Text()
		if err != nil {
			panic(err)
		}
		if output != "Waiting for remote server..." {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	fmt.Printf("%s", strings.Replace(output, "\n\n", "\n", -1))

	wd.Quit()
}
