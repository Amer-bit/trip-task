package tripDetails

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"

	"trip/types"
)

func buildUrl(dto types.Location) string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Access environment variable
	googleApiDomain := os.Getenv("GOOGLE_API_DOMAIN")
	matrixApiKey := os.Getenv("MATRIX_API_KEY")
	reponseFormat := `json`
	matrixPath := `/maps/api/distancematrix/` + reponseFormat
	queryParams := `?destinations=` + url.QueryEscape(dto.Destination) + `&origins=` + url.QueryEscape(dto.Origin) + `&units=metric&key=` + matrixApiKey
	// Specify the API endpoint URL
	apiURL := googleApiDomain + matrixPath + queryParams
	return apiURL
}

func GetTripInfo(dto types.Location) (types.MatrixElement, error) {
	// Make a GET request
	apiURL := buildUrl(dto)
	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return types.MatrixElement{}, err
	}
	// defer response.Body.Close()
	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading the response:", err)
		return types.MatrixElement{}, err
	}
	var data types.MatrixResponse
	parsingErr := json.Unmarshal(body, &data)
	if parsingErr != nil {
		fmt.Println("Error reading the response:", err)
		return types.MatrixElement{}, err
	}
	firstElement := data.Rows[0].Elements[0]
	if firstElement.Status != `OK` {
		return types.MatrixElement{}, errors.New(`Location ` + firstElement.Status)
	}

	return firstElement, nil
}
