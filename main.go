package main

import (
	repository "dynamodb/users"
)

func main() {

	db := repository.New()
	db.GetList()

}
