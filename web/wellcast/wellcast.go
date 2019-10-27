package wellcast

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/naveego/api/pipeline/publisher"
	"github.com/naveego/api/types/pipeline"
)

type Publisher struct {
}

// NewPublisher creates a new Wellcast publisher instance
func NewPublisher() publisher.Publisher {
	return &Publisher{}
}

func (p *Publisher) TestConnection(ctx publisher.Context, connSettings map[string]interface{}) (bool, string, error) {
	return true, "", nil
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

	sqlServer := os.Getenv("SRC_SQL_SERVER")
	sqlPort := os.Getenv("SRC_SQL_PORT")
	sqlUser := os.Getenv("SRC_SQL_USER")
	sqlPwd := os.Getenv("SRC_SQL_PWD")

	connString := fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=synergy", sqlServer, sqlPort, sqlUser, sqlPwd)
	conn, err:= sql.Open("mssql", connString)
	if err != nil {
		return err
	}

	dataPoints := make([]pipeline.DataPoint, len(wells))
	count := 0
	for i, well := range wells {
		count++

		wellData, ok := well.(map[string]interface{})

		if !ok {
			continue
		}

		dataPoint := ctx.NewDataPoint("WellAttribute", []string{"Well ID"}, wellData)
		dataPoints[i] = dataPoint

		if count == 250 {

			if err := insertDataPoints(ctx, conn, "WellAttribute", dataPoints, []string{"Well ID"}); err != nil {
				return err
			}

			count = 0
			dataPoints = []pipeline.DataPoint{}
		}
	}

	if len(dataPoints) > 0 {
		if err := insertDataPoints(ctx, conn, "WellAttribute", dataPoints, []string{"Well ID"}); err != nil {
			return err
		}
	}


	ctx.Logger.Infof("Successfully published %d wells to the pipeline", count)
	return nil

}

func getWells(ctx publisher.Context, authToken string) ([]interface{}, error) {
	ctx.Logger.Info("Fetching well info from API")

	apiURL, valid := getStringSetting(ctx.PublisherInstance.Settings, "api_url")
	if !valid {
		return nil, errors.New("Expected setting for 'api_url' but it was not set or not a valid string.")
	}

	layoutName, valid := getStringSetting(ctx.PublisherInstance.Settings, "layout_name")
	if !valid {
		return nil, errors.New("Expected settings for 'layout_name' but was not set or not a valid string.")
	}

	cli := http.Client{}
	resourceURL := fmt.Sprintf("%s/api/v1/Custom/getLayoutAttribute?Layout=%s", apiURL, layoutName)
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
	var scrubbedWells []interface{}
	err = json.NewDecoder(resp.Body).Decode(&wells)
	if err != nil {
		return nil, fmt.Errorf("Could not read API response: %v", err)
	}

	// We cannot send over any properties that contain restricted property
	// name characters.
	for _, wellRaw := range wells {
		well, ok := wellRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Invalid Well record. Unexpected Type")
		}

		for k := range well {
			if strings.Contains(k, ":") {
				delete(well, k)
			}
		}

		scrubbedWells = append(scrubbedWells, well)
	}

	ctx.Logger.Infof("Successfully fetched %d wells from API", len(scrubbedWells))

	return scrubbedWells, nil

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


func insertDataPoints(ctx publisher.Context, conn *sql.DB, entity string, data []pipeline.DataPoint, keyColumns []string) error {

	firstDp := data[0]

	deleteStmt := fmt.Sprintf("DELETE FROM [wellcast].[%s] WHERE ", entity)
	for _, dp := range data {
		deleteStmt += "("

		for _, pk := range keyColumns {
			pkv := dp.Data[pk]
			pkvS := strings.Replace(fmt.Sprintf("%v", pkv), "'", "''", -1)
			deleteStmt += fmt.Sprintf("[%s] = '%v' AND ", pk, pkvS)
		}
		deleteStmt = deleteStmt[0:len(deleteStmt) - 5] + ") OR "
	}
	deleteStmt = deleteStmt[0:len(deleteStmt) - 4]

	_, e := conn.Exec(deleteStmt)
	if e != nil {
		ctx.Logger.Error(e)
		ctx.Logger.Info(deleteStmt)
	}

	var fields []string
	for k := range firstDp.Data {
		fields = append(fields, k)
	}

	sort.Strings(fields)

	insertStmt := fmt.Sprintf("INSERT INTO [wellcast].[%s] (", entity)
	for _, k := range fields {
		insertStmt += "[" + k + "], "
	}

	insertStmt = insertStmt[0:len(insertStmt)-2] + ") VALUES "

	for _, dp := range data {
		insertStmt += "("

		for _, k := range fields {
			v := dp.Data[k]
			if v == nil {
				insertStmt += "NULL, "
			} else {
				vStr := strings.Replace(fmt.Sprintf("%v", v), "'", "''", -1)
				insertStmt += "'" + vStr + "', "
			}
		}

		insertStmt = insertStmt[0:len(insertStmt)-2] + "),"
	}

	insertStmt = insertStmt[0:len(insertStmt)-1]

	_, e = conn.Exec(insertStmt)
	if e != nil {
		ctx.Logger.Error(e)
		ctx.Logger.Info(insertStmt)
	}
	return nil
}
