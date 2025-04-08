package main

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type TokenMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Error loading env file")
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
	}))

	e.POST("/upload", handleUpload)
	e.GET("/health-check", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server Available")
	})
	e.Logger.Fatal(e.Start(":8080"))
}

func handleUpload(c echo.Context) error {
	name := c.FormValue("name")
	description := c.FormValue("description")
	jsonName := c.FormValue("json_name")

	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Arquivo não encontrado")
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	imageHash, err := uploadFileToPinata(file.Filename, src)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Erro ao subir imagem no Pinata")
	}
	imageURL := "https://gateway.pinata.cloud/ipfs/" + imageHash

	metadata := TokenMetadata{
		Name:        name,
		Description: description,
		Image:       imageURL,
	}

	metaHash, err := uploadJSONToPinata(metadata, jsonName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Erro ao subir metadados no Pinata")
	}

	tokenURI := "https://gateway.pinata.cloud/ipfs/" + metaHash

	return c.JSON(http.StatusOK, map[string]string{
		"tokenURI": tokenURI,
	})
}

func uploadFileToPinata(filename string, file io.Reader) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}
	writer.Close()

	req, err := http.NewRequest("POST", "https://api.pinata.cloud/pinning/pinFileToIPFS", body)
	if err != nil {
		return "", err
	}
	token := os.Getenv("API_JWT")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		IpfsHash string `json:"IpfsHash"`
	}
	json.NewDecoder(resp.Body).Decode(&res)

	return res.IpfsHash, nil
}

func uploadJSONToPinata(metadata TokenMetadata, jsonName string) (string, error) {

	payload := map[string]interface{}{
		"pinataMetadata": map[string]string{
			"name": jsonName,
		},
		"pinataContent": metadata,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.pinata.cloud/pinning/pinJSONToIPFS", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}

	token := os.Getenv("API_JWT")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		IpfsHash string `json:"IpfsHash"`
	}
	json.NewDecoder(resp.Body).Decode(&res)

	return res.IpfsHash, nil
}
