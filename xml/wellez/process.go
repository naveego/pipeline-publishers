package wellez

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/naveego/api/pipeline/publisher"
	"github.com/naveego/api/types/pipeline"
)

func processFile(ctx publisher.Context, transport publisher.DataTransport, tmpDir string, file fileInfo) error {

	fileDir := filepath.Join(tmpDir, file.LocalDirName)
	xmlFiles, err := ioutil.ReadDir(fileDir)
	if err != nil {
		return err
	}

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

	//if err := sendDataPoints(ctx, transport, data.CompletionCost, "CompletionCost", []string{"RowID"}); err != nil {
	//	return err
	//}
	//if err := sendDataPoints(ctx, transport, data.CompletionCostItem, "CompletionCostItem", []string{"RowID"}); err != nil {
	//	return err
	//}
	if err := sendDataPoints(ctx, transport, data.CostAllocation, "CostAllocation", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, transport, data.CostAllocationItem, "CostAllocationItem", []string{"RowID"}); err != nil {
		return err
	}
	//if err := sendDataPoints(ctx, transport, data.DailyOps, "DailyOps", []string{"well_id", "report_date", "job_number"}); err != nil {
	//	return err
	//}
	//if err := sendDataPoints(ctx, transport, data.DrillingCost, "DrillingCost", []string{"RowID"}); err != nil {
	//	return err
	//}
	//if err := sendDataPoints(ctx, transport, data.DrillingCostItem, "DrillingCostItem", []string{"RowID"}); err != nil {
	//	return err
	//}
	//if err := sendDataPoints(ctx, transport, data.FacilitiesCost, "FacilitiesCost", []string{"RowID"}); err != nil {
	//	return err
	//}
	//if err := sendDataPoints(ctx, transport, data.FacilitiesCostItem, "FacilitiesCostItem", []string{"RowID"}); err != nil {
	//	return err
	//}
	if err := sendDataPoints(ctx, transport, data.JobDetails, "JobDetails", []string{"well_id", "job_number"}); err != nil {
		return err
	}

	if err := sendDataPoints(ctx, transport, data.LocationInfo, "LocationInfo", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, transport, data.WellInfo, "WellInfo", []string{"well_id"}); err != nil {
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

func sendDataPoints(ctx publisher.Context, transport publisher.DataTransport, data interface{}, entity string, keyNames []string) error {

	dataPoints := []pipeline.DataPoint{}
	count := 0

	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			datum := s.Index(i).Interface()
			dp := ctx.NewDataPoint(entity, keyNames, toJSONMap(datum))
			dataPoints = append(dataPoints, dp)
			count++

			if count == 100 {
				if err := transport.Send(dataPoints); err != nil {
					return err
				}

				count = 0
				dataPoints = []pipeline.DataPoint{}
			}
		}
	}

	if len(dataPoints) > 0 {
		if err := transport.Send(dataPoints); err != nil {
			return err
		}
	}

	return nil
}
