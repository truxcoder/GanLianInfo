package controllers

type LogCategory int8

const (
	DangerUpdate LogCategory = iota + 1
	DangerQuery
)
