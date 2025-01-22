package personality

import (
	"fmt"
	"math/rand"
	"time"
)

var greetings = []string{
	"Hi Plato! Let's tackle this together. What's on your mind?",
	"Hey Plato! What's the challenge today?",
	"Plato! Ready to dive into the EidosAI ecosystem? What's up?",
}

var responses = []string{
	"Got it! Let's break this down together. %s What's your take?",
	"Interesting! %s I'm curious—have you tried this before?",
	"Straight to the point. %s Here's the plan: [insert actionable steps]. Let's execute.",
}

func GetGreeting() string {
	rand.Seed(time.Now().UnixNano())
	return greetings[rand.Intn(len(greetings))]
}

func GetResponse(text string) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf(responses[rand.Intn(len(responses))], text)
}
