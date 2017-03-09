package mssql

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/naveego/api/pipeline/publisher"
	"github.com/naveego/api/types/pipeline"
	"github.com/naveego/api/utils"
)

type Publisher struct{}

// NewPublisher creates a new MSSQL publisher instance
func NewPublisher() publisher.Publisher {
	return &Publisher{}
}

func (p *Publisher) TestConnection(ctx publisher.Context, connSettings map[string]interface{}) (bool, string, error) {
	connString, err := buildConnectionString(connSettings, 10)
	if err != nil {
		return false, "could not connect to server", err
	}
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		return false, "could not connect to server", err
	}
	defer conn.Close()

	stmt, err := conn.Prepare("select 1")
	if err != nil {
		return false, "could not connect to server", err
	}
	defer stmt.Close()

	return true, "successfully connected to server", nil
}

func (p *Publisher) Shapes(ctx publisher.Context) (pipeline.ShapeDefinitions, error) {
	defs := pipeline.ShapeDefinitions{}

	connString, err := buildConnectionString(ctx.PublisherInstance.Settings, 30)
	if err != nil {
		return defs, err
	}
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		return defs, err
	}
	defer conn.Close()

	rows, err := conn.Query(`select s.Name, o.Name, c.Name, ty.name from
sys.objects o
INNER JOIN sys.schemas s ON (o.schema_id = s.schema_id)
INNER JOIN sys.columns c ON (o.object_id = c.object_id)
INNER JOIN sys.types ty ON (c.user_type_id = ty.user_type_id)
where type IN ('U', 'V')
ORDER BY s.Name, o.Name, c.column_id`)

	if err != nil {
		return defs, err
	}
	defer rows.Close()

	var schemaName string
	var tableName string
	var columnName string
	var columnType string

	s := map[string]*pipeline.ShapeDefinition{}

	for rows.Next() {
		err = rows.Scan(&schemaName, &tableName, &columnName, &columnType)
		if err != nil {
			continue
		}

		shapeName := tableName
		if schemaName != "dbo" {
			shapeName = fmt.Sprintf("%s_%s", schemaName, tableName)
		}

		shapeDef, ok := s[shapeName]
		if !ok {
			shapeDef = &pipeline.ShapeDefinition{
				Name: shapeName,
			}
			s[shapeName] = shapeDef
		}

		shapeDef.Properties = append(shapeDef.Properties, pipeline.PropertyDefinition{
			Name: columnName,
			Type: convertSQLType(columnType),
		})
	}

	for _, sd := range s {
		defs = append(defs, *sd)
	}

	// Sort the shapes by Name
	sort.Sort(pipeline.SortShapesByName(defs))

	return defs, nil
}

func (p *Publisher) Publish(ctx publisher.Context, dataTransport publisher.DataTransport) {
}

func buildConnectionString(settings map[string]interface{}, timeout int) (string, error) {
	mr := utils.NewMapReader(settings)
	server, ok := mr.ReadString("server")
	if !ok {
		return "", errors.New("server cannot be null or empty")
	}
	db, ok := mr.ReadString("database")
	if !ok {
		return "", errors.New("database cannot be null or empty")
	}
	auth, ok := mr.ReadString("auth")
	if !ok {
		return "", errors.New("auth type must be provided")
	}

	connStr := []string{
		"server=" + server,
		"database=" + db,
		"connection timeout=10",
	}

	if auth == "sql" {
		username, _ := mr.ReadString("username")
		pwd, _ := mr.ReadString("password")
		connStr = append(connStr, "user id="+username)
		connStr = append(connStr, "password="+pwd)
	} else {
		connStr = append(connStr, "")
	}

	return strings.Join(connStr, ";"), nil
}

func convertSQLType(t string) string {
	switch t {
	case "datetime":
	case "date":
	case "time":
	case "smalldatetime":
		return "date"
	case "bigint":
	case "int":
	case "smallint":
	case "tinyint":
		return "integer"
	case "decimal":
	case "float":
	case "money":
	case "smallmoney":
		return "float"
	case "bit":
		return "bool"
	}

	return "string"
}
