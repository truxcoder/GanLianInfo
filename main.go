package main

import "GanLianInfo/router"

func main() {
	r := router.Register()
	router.Start(r)
}
