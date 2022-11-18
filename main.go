package main

import (
	"encoding/json"
	"fmt"
	"log"
	rand2 "math/rand"
	"net/http"
	"os"
)

//var count int = 0

func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	http.HandleFunc("/", handler)

	log.Printf("starting server on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("http listen error: %v", err)
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Fprint(w, "Let the battle begin!")
		return
	}

	var v ArenaUpdate
	defer req.Body.Close()
	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&v); err != nil {
		log.Printf("WARN: failed to decode ArenaUpdate in response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := play(v)
	fmt.Fprint(w, resp)
}

func play(input ArenaUpdate) string {
	log.Printf("IN: %#v", input)

	commands := []string{"F", "R", "L", "T"}
	// var rand int = 3
	// if count == 3 {
	// 	rand = rand2.Intn(4)
	// }
	// //TODO add your implementation here to replace the random response
	// count++
	// count = count % 4

	//return commands[rand]

	res := process(input)
	if res == "" {
		log.Println("Status", "got empty sending random")
		res = commands[rand2.Intn(4)]
	}
	log.Println("Status", "res:", res)

	return res
}

func increment(dirMap map[string]int, dir string) {
	val, ok := dirMap[dir]
	if !ok {
		dirMap[dir] = 1
	}
	dirMap[dir] = val + 1
}

func process(input ArenaUpdate) string {
	me := input.Arena.State[input.Links.Self.Href]
	if me.WasHit {
		log.Println("Status", "got hit from someone")
		return ""
	}
	dirMap := make(map[string]int, 4)
	diffX, diffY := 0, 0
	switch me.Direction {
	case "N":
		for _, pl := range input.Arena.State {
			diffX, diffY = me.X-pl.X, me.Y-pl.Y

			if diffX == 0 && diffY > 0 && diffY <= 3 {
				log.Println("Status", "got someone in current direction")
				return "T"
			}
			if diffX > 0 && diffX <= 3 {
				increment(dirMap, "R")
			} else if diffX < 0 && diffX >= -3 {
				increment(dirMap, "L")
			} else {
				increment(dirMap, "F")
			}
		}

	case "E":
		for _, pl := range input.Arena.State {
			diffX, diffY = me.X-pl.X, me.Y-pl.Y

			if diffY == 0 && diffX > 0 && diffX <= 3 {
				log.Println("Status", "got someone in current direction")
				return "T"
			}
			if diffY > 0 && diffY <= 3 {
				increment(dirMap, "L")
			} else if diffY < 0 && diffY >= -3 {
				increment(dirMap, "R")
			} else {
				increment(dirMap, "F")
			}
		}

	case "S":
		for _, pl := range input.Arena.State {
			diffX, diffY = me.X-pl.X, me.Y-pl.Y

			if diffX == 0 && diffY < 0 && diffY >= -3 {
				log.Println("Status", "got someone in current direction")
				return "T"
			}
			if diffX > 0 && diffX <= 3 {
				increment(dirMap, "L")
			} else if diffX < 0 && diffX >= -3 {
				increment(dirMap, "R")
			} else {
				increment(dirMap, "F")
			}
		}
	case "W":
		for _, pl := range input.Arena.State {
			diffX, diffY = me.X-pl.X, me.Y-pl.Y

			if diffY == 0 && diffX < 0 && diffX >= -3 {
				log.Println("Status", "got someone in current direction")
				return "T"
			}
			if diffY > 0 && diffY <= 3 {
				increment(dirMap, "L")
			} else if diffY < 0 && diffY >= -3 {
				increment(dirMap, "R")
			} else {
				increment(dirMap, "F")
			}
		}
	}
	var (
		max     int    = 0
		nextDir string = "U"
	)
	for k, v := range dirMap {
		if v > max {
			max = v
			nextDir = k
		}
	}
	if max == 0 && nextDir == "F" {
		nextDir = "" // will denote go for random
	}
	log.Println("Status", dirMap, max, nextDir)
	return nextDir
}
