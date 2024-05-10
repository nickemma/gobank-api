package main

func main() {
	server := NewAPIServer(":5000")
	server.Run()
}
