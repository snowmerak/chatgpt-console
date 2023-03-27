package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	isDesigner := flag.Bool("designer", false, "designer mode")
	isAdvisor := flag.Bool("adviser", false, "adviser mode")
	isArchitect := flag.Bool("architect", false, "architect mode")
	isExpert := flag.Bool("expert", false, "expert mode")
	isDeveloper := flag.Bool("developer", false, "developer mode")

	flag.Parse()

	systemPrompt := ""
	if *isDesigner {
		systemPrompt = designer
	} else if *isAdvisor {
		systemPrompt = advisor
	} else if *isArchitect {
		systemPrompt = architect
	} else if *isExpert {
		systemPrompt = expert
	} else if *isDeveloper {
		systemPrompt = developer
	}

	client := &http.Client{}

	for {
		fmt.Print("> ")
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			log.Println(err)
			continue
		}

		req, err := NewSimpleChatGPTRequest(systemPrompt, input)
		if err != nil {
			log.Println(err)
			continue
		}

		post, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(req))
		if err != nil {
			log.Println(err)
			continue
		}

		post.Header.Set("Content-Type", "application/json")
		post.Header.Set("Authorization", "bearer "+chatGPTKey)

		func() {
			resp, err := client.Do(post)
			if err != nil {
				log.Println(err)
				return
			}
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				return
			}

			res, err := ContentsFromResponseChatGPT3(respBody)
			if err != nil {
				log.Println(err)
				return
			}

			fmt.Println(">>", strings.Join(res, "\n"))
		}()
	}
}
