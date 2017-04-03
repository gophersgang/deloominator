package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"github.com/lucapette/deloominator/pkg/app"
)

type cell struct {
	Value string `json:"value"`
}

type row struct {
	Cells []cell `json:"cells"`
}

type column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type rawResults struct {
	Total   int      `json:"total"`
	Columns []column `json:"columns"`
	Rows    []row    `json:"rows"`
}

type table struct {
	Name string `json:"name"`
}

type dataSource struct {
	Name   string   `json:"name"`
	Tables []*table `json:"tables"` //better if it's a rawResults.
}

var schema graphql.Schema

func GraphQLHandler(app *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/json")

		query, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		res := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: string(query),
			Context:       context.WithValue(context.Background(), "app", app),
		})

		if res.HasErrors() {
			w.WriteHeader(http.StatusBadRequest)
		}

		rJSON, err := json.Marshal(res)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(rJSON)
	}
}

func ResolveDataSources(p graphql.ResolveParams) (interface{}, error) {
	var dataSources []*dataSource
	app := p.Context.Value("app").(*app.App)
	for _, ds := range app.GetDataSources() {
		name := ds.DSN().DBName
		log.WithField("schema_name", name).Info("query metadata")

		start := time.Now()

		qr, err := ds.Tables()
		if err != nil {
			return dataSources, err
		}

		ts := make([]*table, len(qr.Rows))
		for i, t := range qr.Rows {
			ts[i] = &table{Name: t[0].Value}
		}

		log.WithFields(log.Fields{
			"schema_name": name,
			"n_tables":    len(qr.Rows),
			"spent":       time.Now().Sub(start),
		}).Info("tables loaded")

		dataSources = append(dataSources, &dataSource{Name: name, Tables: ts})
	}

	return dataSources, nil
}

func ResolveQuery(p graphql.ResolveParams) (interface{}, error) {
	source := p.Args["source"].(string)
	input := p.Args["input"].(string)
	app := p.Context.Value("app").(*app.App)

	qr, err := app.GetDataSources()[source].Query(input)

	columns := make([]column, len(qr.Columns))

	for i, col := range qr.Columns {
		columns[i].Name = col.Name
	}

	rows := make([]row, len(qr.Rows))

	for i, r := range qr.Rows {
		rows[i].Cells = make([]cell, len(qr.Columns))

		for j, c := range r {
			rows[i].Cells[j].Value = c.Value
		}
	}

	return rawResults{
		Total:   len(qr.Rows),
		Columns: columns,
		Rows:    rows,
	}, err
}

func init() {
	cellType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Cell",
		Description: "A cell represents a single piece of returnted data",
		Fields: graphql.Fields{
			"value": &graphql.Field{
				Description: "Value of the cell",
				Type:        graphql.NewNonNull(graphql.String),
			},
		},
	})
	rowType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Row",
		Description: "A row holds the representation of a set of cells of the raw data returned by the data source",
		Fields: graphql.Fields{
			"cells": &graphql.Field{
				Description: "Name of the column",
				Type:        graphql.NewList(cellType),
			},
		},
	})

	columnType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Column",
		Description: "A column holds the representation of one columnd of the raw data returned by a data source",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Description: "Name of the column",
				Type:        graphql.String,
			},
			"type": &graphql.Field{
				Description: "Type of the column",
				Type:        graphql.String,
			},
		},
	})

	rawResultsType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "RawResults",
		Description: "RawResults represents a collection of raw data returned by a data source",
		Fields: graphql.Fields{
			"total": &graphql.Field{
				Description: "Total count of returned results",
				Type:        graphql.Int,
			},
			"columns": &graphql.Field{
				Description: "Columns of the returned results",
				Type:        graphql.NewList(columnType),
			},
			"rows": &graphql.Field{
				Description: "Rows of the returned results",
				Type:        graphql.NewList(rowType),
			},
		},
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			return true
		},
	})

	queryResultType := graphql.NewUnion(graphql.UnionConfig{
		Name:        "QueryResult",
		Description: "QueryResult represents all the possible outcomes of a Query",
		Types:       []*graphql.Object{rawResultsType},
	})

	tableType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Table",
		Description: fmt.Sprintf("A table of a data source"),
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	dataSourceType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "DataSource",
		Description: fmt.Sprintf("A DataSource represents a single source of data to analyze"),
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"tables": &graphql.Field{
				Type: graphql.NewList(tableType),
			},
		},
	})

	fields := graphql.Fields{
		"DataSources": &graphql.Field{
			Type:    graphql.NewList(dataSourceType),
			Resolve: ResolveDataSources,
		},
		"Query": &graphql.Field{
			Type: queryResultType,
			Args: graphql.FieldConfigArgument{
				"source": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: ResolveQuery,
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "Query", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}

	var err error
	schema, err = graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
}
