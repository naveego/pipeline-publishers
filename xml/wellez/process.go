package wellez

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/naveego/api/pipeline/publisher"
	"github.com/naveego/api/types/pipeline"

)

func processFile(ctx publisher.Context, transport publisher.DataTransport, tmpDir string, file fileInfo) error {

	fileDir := filepath.Join(tmpDir, file.LocalDirName)
	var xmlFiles []os.FileInfo

	filepath.Walk(fileDir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if strings.HasSuffix(f.Name(), ".xml") {
				xmlFiles = append(xmlFiles, f)
			}
		}

		return nil
	})

	if len(xmlFiles) != 1 {
		return fmt.Errorf("Expected 1 file per archive, but found %d", len(xmlFiles))
	}

	xmlFile := xmlFiles[0].Name()

	ctx.Logger.Infof("Processing data from file '%s", xmlFile)

	f, err := os.Open(filepath.Join(tmpDir, file.LocalDirName, xmlFile))
	if err != nil {
		return err
	}
	defer f.Close()

	data := Data{}
	err = xml.NewDecoder(f).Decode(&data)
	if err != nil {
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

	if err := sendDataPoints(ctx, conn, data.CompletionCost, "CompletionCost", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, conn, data.CompletionCostItem, "CompletionCostItem", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, conn, data.CostAllocation, "CostAllocation", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, conn, data.CostAllocationItem, "CostAllocationItem", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, conn, data.DailyOps, "DailyOps", []string{"well_id", "report_Date", "job_number"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, conn, data.DrillingCost, "DrillingCost", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, conn, data.DrillingCostItem, "DrillingCostItem", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, conn, data.FacilitiesCost, "FacilitiesCost", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, conn, data.FacilitiesCostItem, "FacilitiesCostItem", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, conn, data.JobDetails, "JobDetails", []string{"well_id", "job_number"}); err != nil {
		return err
	}

	if err := sendDataPoints(ctx, conn, data.LocationInfo, "LocationInfo", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, conn, data.WellInfo, "WellInfo", []string{"well_id"}); err != nil {
		return err
	}

	return nil
}

func toJSONMap(s interface{}) map[string]interface{} {
	m := map[string]interface{}{}
	b, _ := json.Marshal(s)
	d := json.NewDecoder(strings.NewReader(string(b)))
	d.UseNumber()
	d.Decode(&m)
	return m
}

func sendDataPoints(ctx publisher.Context, conn *sql.DB, data interface{}, entity string, keyNames []string) error {

	dataPoints := []pipeline.DataPoint{}
	count := 0
	total := 0



	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			datum := s.Index(i).Interface()
			dp := ctx.NewDataPoint(entity, keyNames, toJSONMap(datum))
			dataPoints = append(dataPoints, dp)
			count++
			total++

			if count == 250 {

				if err := insertDataPoints(ctx, conn, entity, dataPoints, keyNames); err != nil {
					return err
				}

				count = 0
				dataPoints = []pipeline.DataPoint{}
			}

			if (total % 10000) == 0 {
				ctx.Logger.Info("Loaded %v data points for %s", total, entity)
			}
		}
	}

	if len(dataPoints) > 0 {
		if err := insertDataPoints(ctx, conn, entity, dataPoints, keyNames); err != nil {
			return err
		}
	}

	ctx.Logger.Info(fmt.Sprintf("Successfully processed %v data points for %s", total, entity))

	return nil
}

func insertDataPoints(ctx publisher.Context, conn *sql.DB, entity string, data []pipeline.DataPoint, keyColumns []string) error {

	firstDp := data[0]

	deleteStmt := fmt.Sprintf("DELETE FROM [wellez].[%s] WHERE ", entity)
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

	insertStmt := fmt.Sprintf("INSERT INTO [wellez].[%s] (", entity)
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


