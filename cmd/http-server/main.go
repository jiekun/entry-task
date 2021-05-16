package main

import "github.com/2014bduck/entry-task/internal/routers"

func main() {
	r := routers.NewRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}