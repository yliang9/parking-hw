package main

import (
    "fmt"
    "errors"
)

const (
    MAP = iota
    DB
)

func getRepo(t int) (repository, error){
    switch t {
	case MAP:
		return GetMapRepoInstance(),nil
	case DB:
        //NOT implemented yet
		return nil, errors.New(fmt.Sprintf("Not implemented yet"))
	default:
		//if type is invalid, return an error
		return nil, errors.New("Invalid Repository Type")
	}
    return nil, nil
}
