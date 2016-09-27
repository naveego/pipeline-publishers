package wellcast

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/naveego/api/pipeline/publisher"
	"github.com/naveego/api/types/pipeline"
)

type Publisher struct {
}

// NewPublisher creates a new Wellcast publisher instance
func NewPublisher() publisher.Publisher {
	return &Publisher{}
}

func (p *Publisher) Shapes(ctx publisher.Context) (map[string]pipeline.Shape, error) {
	return nil, nil
}

func (p *Publisher) Publish(ctx publisher.Context, dataTransport publisher.DataTransport) {
	writeCommonLogs(ctx, "publish")
	publishWells(ctx, dataTransport)
}

func publishWells(ctx publisher.Context, dataTransport publisher.DataTransport) error {
	authToken, err := getAuthToken(ctx)
	if err != nil {
		ctx.Logger.Error("Error authenticating to API: ", err)
		return err
	}

	wells, err := getWells(ctx, authToken)
	if err != nil {
		ctx.Logger.Error("Error fetching wells from API: ", err)
		return err
	}

	dataPoints := make([]pipeline.DataPoint, len(wells))

	for i, well := range wells {

		wellData, ok := well.(map[string]interface{})

		if !ok {
			continue
		}

		dataPoint := ctx.NewDataPoint("WellAttribute", []string{"$id"}, wellData)
		dataPoints[i] = dataPoint
	}

	err = dataTransport.Send(dataPoints)
	if err != nil {
		ctx.Logger.Error("Error publishing data to pipeline: ", err)
		return err
	}

	ctx.Logger.Infof("Successfully published %d wells to the pipeline", len(dataPoints))
	return nil

}

func getWells(ctx publisher.Context, authToken string) ([]interface{}, error) {
	ctx.Logger.Info("Fetching well info from API")

	apiURL, valid := getStringSetting(ctx.PublisherInstance.Settings, "api_url")
	if !valid {
		return nil, errors.New("Expected setting for 'api_url' but it was not set or not a valid string.")
	}

	cli := http.Client{}
	resourceURL := fmt.Sprintf("%s/api/v1/Custom/getWellAttribute", apiURL)
	req, err := http.NewRequest("GET", resourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Authorization", authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cli.Do(req)
	if resp == nil && err != nil {
		return nil, fmt.Errorf("Could not execute request: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API Error. Status Code %d", resp.StatusCode)
	}

	var wells []interface{}
	err = json.NewDecoder(resp.Body).Decode(&wells)
	if err != nil {
		return nil, fmt.Errorf("Could not read API response: %v", err)
	}

	ctx.Logger.Infof("Successfully fetched %d wells from API", len(wells))

	return wells, nil

}

func getAuthToken(ctx publisher.Context) (string, error) {
	ctx.Logger.Info("Authenticating to Wellcast Api")

	apiURL, ok := getStringSetting(ctx.PublisherInstance.Settings, "api_url")
	if !ok {
		return "", fmt.Errorf("Expected setting for 'api_url' but it was not set or not a valid string.")
	}

	user, ok := getStringSetting(ctx.PublisherInstance.Settings, "user")
	if !ok {
		return "", fmt.Errorf("Expected setting for 'user' but it was not set or not a valid string")
	}

	password, ok := getStringSetting(ctx.PublisherInstance.Settings, "password")
	if !ok {
		return "", fmt.Errorf("Expected setting for 'password' but it was not set or not a valid string")
	}

	cli := http.Client{}
	authURL := fmt.Sprintf("%s/api/v2/auth/token?userName=%s&password=%s", apiURL, user, password)

	ctx.Logger.Debugf("Calling Authentication Endpoint with URL: %s", authURL)

	req, err := http.NewRequest("POST", authURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := cli.Do(req)
	if resp == nil && err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return "", fmt.Errorf("The API returned a status code of %d", resp.StatusCode)
	}

	var respJSON map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respJSON)
	if err != nil {
		return "", fmt.Errorf("Error decoding response: %v", err)
	}

	rawAuthToken, ok := respJSON["AuthToken"]
	if !ok {
		return "", errors.New("The response did not contain an AuthToken property")
	}

	authToken, ok := rawAuthToken.(string)
	if !ok {
		return "", errors.New("The response contained an AuthToken property that was not a valid string")
	}

	return authToken, nil

}

func getStringSetting(settings map[string]interface{}, name string) (string, bool) {

	rawValue, ok := settings[name]
	if !ok {
		return "", false
	}

	value, ok := rawValue.(string)
	if !ok {
		return "", false
	}

	return value, true
}

func writeCommonLogs(ctx publisher.Context, action string) {
	ctx.Logger.Infof("Starting action %s", action)
	apiURL, _ := getStringSetting(ctx.PublisherInstance.Settings, "api_url")
	ctx.Logger.Infof("Using API Endpoint: %s", apiURL)
}
