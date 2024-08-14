package main

import (
	"dynamodb/repository"
)

func main() {

	db := repository.New()

	db.GetList()

}
