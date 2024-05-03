# go-rest-api

Go application that serves as a RESTful API for managing email data. It includes functionalities for storing email data in a MySQL database, retrieving email data, converting it to XML format with associated binary data encoded in base64, saving the XML file locally, moving the XML file to another location, and updating the corresponding database record.

## Prerequisites
- Go installed on your system
- MySQL database accessible
- Required Go packages:
   - github.com/gin-gonic/gin for HTTP web framework
   - github.com/go-sql-driver/mysql for MySQL driver
 
## Installation
- Clone or download the repository to your local machine.
- Navigate to the project directory in your terminal.
- Run the following command to install dependencies:
    - go mod tidy

## Usage
- go run main.go

## Testing
- You can test the endpoints using tools like cURL, Postman, or any HTTP client of your choice.
