package main

import(
	"github.com/it234/gowebssh/internal/app/manageweb"
)

func main(){
	manageweb.Run("")
	select{}
}