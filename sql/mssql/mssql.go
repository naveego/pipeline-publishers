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

type Publisher struct {
	conn *sql.DB
}

// NewPublisher creates a new MSSQL publisher instance
func NewPublisher() publisher.Publisher {
	return &Publisher{}
}

func (p *Publisher) TestConnection(ctx publisher.Context) (bool, string, error) {
	connString, err := buildConnectionString(ctx.Settings, 10)
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

func (p *Publisher) Init(ctx publisher.Context) error {
	connString, err := buildConnectionString(ctx.Settings, 30)
	if err != nil {
		ctx.Logger.Warn("Could not connect to publisher", err)
		return err
	}
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		ctx.Logger.Warn("Could not open connection to publisher", err)
		return err
	}
	p.conn = conn
	return nil
}

func (p *Publisher) Dispose(ctx publisher.Context) error {
	if p.conn != nil {
		p.conn.Close()
		p.conn = nil
	}

	return nil
}

func (p *Publisher) Shapes(ctx publisher.Context) (pipeline.ShapeDefinitions, error) {
	defs := pipeline.ShapeDefinitions{}

	connString, err := buildConnectionString(ctx.Settings, 30)
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
			shapeName = fmt.Sprintf("%s__%s", schemaName, tableName)
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

func (p *Publisher) Publish(ctx publisher.Context, shape pipeline.ShapeDefinition, dataTransport publisher.DataTransport) {

	columns := []string{}

	for _, v := range shape.Properties {
		columns = append(columns, v.Name)
	}

	schemaName := "dbo"
	tableName := shape.Name

	if strings.Contains(shape.Name, "__") {
		idx := strings.Index(shape.Name, "__")
		schemaName = tableName[:idx]
		tableName = tableName[idx+2:]
	}

	colStr := "[" + strings.Join(columns, "],[") + "]"
	query := fmt.Sprintf("SELECT %s FROM [%s].[%s]", colStr, schemaName, tableName)

	ctx.Logger.Debugf("Query: %s", query)
	rows, err := p.conn.Query(query)
	if err != nil {
		ctx.Logger.Warn("Could not query publisher", err)
		return
	}
	defer rows.Close()

	vals := make([]interface{}, len(columns))
	for i := 0; i < len(columns); i++ {
		vals[i] = new(interface{})
	}

	for rows.Next() {
		err := rows.Scan(vals...)
		if err != nil {
			ctx.Logger.Warn("Error reading row", err)
			continue
		}

		d := map[string]interface{}{}

		for i := 0; i < len(vals); i++ {
			colName := columns[i]
			d[colName] = vals[i]
		}

		dp := pipeline.DataPoint{
			Repository: "matador",
			Data:       d,
		}

		dps := []pipeline.DataPoint{dp}
		dataTransport.Send(dps)

	}

	dataTransport.Done()

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
