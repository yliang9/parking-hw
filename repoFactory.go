package main

import (
    "errors"
)

const (
    MAP = iota
    DB //as an example, not supported yet
)

//getRepo use Factory pattern to create a repository as it can be easily swap in and out to support different repository
func getRepo(t int) (repository, error){
    switch t {
	case MAP:
		return GetMapRepoInstance(),nil
	case DB:
        //NOT implemented yet
		return nil, errors.New("Not implemented yet")
	default:
		//if type is invalid, return an error
		return nil, errors.New("Invalid Repository Type")
	}
    return nil, nil
}
