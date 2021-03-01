package main

import "zinx/znet"

func main() {
	s := znet.NewServer("[zinx V0.2]")
	s.Serve()
}
