package main

import (
	"encoding/xml"
)

type EmailData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type XMLData struct {
	XMLName xml.Name `xml:"data"`
	JSON    EmailData `xml:"json"`
	Base64  string   `xml:"base64"`
}