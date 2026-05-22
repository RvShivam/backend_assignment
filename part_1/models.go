package main

type RequestBody struct {
	UserID  string      `json:"user_id"`
	Payload interface{} `json:"payload" binding:"required"`
}