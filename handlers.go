package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"time"
	"log"

	"github.com/gin-gonic/gin"
)

func handleRequest(c *gin.Context) {
	log.Println("handleRequest started")

	var data EmailData
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Printf("Error binding JSON: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received data: %+v", data)
	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempt %d to insert or update record", i+1)

		_, err := db.Exec("INSERT INTO email_data (id, name, email) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE name = ?, email = ?", data.ID, data.Name, data.Email, data.Name, data.Email)
		if err == nil {
			log.Printf("Successfully inserted or updated record with ID: %d, Name: %s, Email: %s\n", data.ID, data.Name, data.Email)
			c.Status(http.StatusCreated)
			c.String(http.StatusOK, "POST Request Completed Successfully.")
			return
		}

		log.Printf("Error inserting or updating record: %s", err)
		time.Sleep(time.Second * time.Duration(i*i))
	}
	log.Println("Failed to insert or update record after multiple attempts")
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert or update record after multiple attempts"})
}

func handleBinaryRequest(c *gin.Context) {
	log.Println("handleBinaryRequest started")

	id := c.Param("id")
	var data EmailData
	
	err := db.QueryRow("SELECT id, name, email FROM email_data WHERE id = ?", id).Scan(&data.ID, &data.Name, &data.Email)
	if err != nil {
		log.Printf("Error querying database: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Retrieved data: %+v", data)
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling data to JSON: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Encode jsonData to base64
	base64Data := base64.StdEncoding.EncodeToString(jsonData)

	xmlDataStruct := XMLData{
		JSON:   data,
		Base64: base64Data,
	}
	
	xmlData, err := xml.MarshalIndent(xmlDataStruct, "", "  ")
	if err != nil {
		log.Printf("Error marshalling data to XML: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	currentDate := time.Now().Format("20060102")
	newFilename := currentDate + ".xml"

	// Write XML data to a file
	err = os.WriteFile(newFilename, xmlData, 0644)
	if err != nil {
		log.Printf("Error writing XML data to file: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("XML file has been successfully created and saved locally.")
	c.String(http.StatusOK, "XML file has been successfully created and saved locally.")
}

func moveXMLFile(c *gin.Context) {
	log.Println("moveXMLFile started")
	
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	currentDate := time.Now().Format("20060102")
	newFilename := currentDate + ".xml"
	
	newLocation := dir + "/outputs/" + newFilename

	xmlFile, err := os.ReadFile(newFilename)
	if err != nil {
		log.Printf("Error reading XML file: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = os.WriteFile(newLocation, xmlFile, 0644)
	if err != nil {
		log.Printf("Error writing XML file to new location: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("XML file has been successfully moved.")
	err = os.Remove(newFilename)
	if err != nil {
		log.Printf("Error removing original XML file: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var xmlData XMLData
	err = xml.Unmarshal(xmlFile, &xmlData)
	if err != nil {
		log.Printf("Error unmarshalling XML data: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Exec("UPDATE email_data SET name = ?, email = ? WHERE id = ?", xmlData.JSON.Name, xmlData.JSON.Email, xmlData.JSON.ID)
	if err != nil {
		log.Printf("Error updating database record: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Database record has been successfully updated.")
	c.String(http.StatusOK, "XML file has been successfully moved to the connected archive and the database record has been updated.")
}

