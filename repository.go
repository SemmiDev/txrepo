package main

type Repository interface {
	Creator
	Finder
	Updater
}

var _ Repository = (*Queries)(nil)
