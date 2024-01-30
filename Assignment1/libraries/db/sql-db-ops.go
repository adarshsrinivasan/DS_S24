package db

import (
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/common"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"github.com/oiime/logrusbun"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var (
	SqlDBClient *bun.DB
)

const (
	PostgresHostEnv     = "POSTGRES_HOST"
	PostgresPortEnv     = "POSTGRES_PORT"
	PostgresUsernameEnv = "POSTGRES_USERNAME"
	PostgresPasswordEnv = "POSTGRES_PASSWORD"
	PostgresDbEnv       = "POSTGRES_DB"

	DefaultPoolSize     = 8   // default connection pool size.
	DefaultIdleTimeouts = -1  // never timeout/close an idle connection.
	DefaultLogSlowQuery = 100 // log db queries slower than 100ms by default

	SchemaName = "marketplace"
)

func NewSQLClient(ctx context.Context, applicationName string) (*bun.DB, error) {
	logrus.Infof("NewSQLClient: Initializing new SQL DB CLient...\n")
	host := common.GetEnv(PostgresHostEnv, "localhost")
	port := common.GetEnv(PostgresPortEnv, "5432")
	username := common.GetEnv(PostgresUsernameEnv, "admin")
	password := common.GetEnv(PostgresPasswordEnv, "admin")
	dbName := common.GetEnv(PostgresDbEnv, "marketplace")

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", host, port)),
		pgdriver.WithUser(username),
		pgdriver.WithPassword(password),
		pgdriver.WithDatabase(dbName),
		pgdriver.WithApplicationName(applicationName),
		pgdriver.WithConnParams(map[string]interface{}{
			"search_path": fmt.Sprintf("%s", SchemaName),
		}),
		pgdriver.WithTLSConfig(&tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: host,
		})))
	sqldb.SetConnMaxIdleTime(DefaultIdleTimeouts)
	sqldb.SetMaxIdleConns(DefaultPoolSize)
	sqldb.SetMaxOpenConns(DefaultPoolSize)
	sqlDB, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName))
	if err != nil {
		err = fmt.Errorf("exception while opening postgres connection: %v", err)
		logrus.Errorf("NewSQLClient: %v\n", err)
		return nil, err
	}
	db := bun.NewDB(sqlDB, pgdialect.New())
	logrusObj := logrus.New()
	logrusObj.SetFormatter(&logrus.TextFormatter{DisableQuote: true})
	db.AddQueryHook(logrusbun.NewQueryHook(logrusbun.QueryHookOptions{
		Logger:          logrusObj,
		LogSlow:         DefaultLogSlowQuery,
		QueryLevel:      logrus.DebugLevel,
		SlowLevel:       logrus.WarnLevel,
		ErrorLevel:      logrus.ErrorLevel,
		MessageTemplate: "{{.Operation}}[{{.Duration}}]: {{.Query}}",
		ErrorTemplate:   "{{.Operation}}[{{.Duration}}]: {{.Query}}: {{.Error}}",
	}))

	if err := db.Ping(); err != nil {
		err := fmt.Errorf("exception while pinging postgres DB. %v", err)
		logrus.Errorf("NewSQLClient: %v\n", err)
		return nil, err
	}

	_, err = db.ExecContext(ctx, "CREATE SCHEMA IF NOT EXISTS ?", bun.Ident(SchemaName))
	if err != nil {
		err := fmt.Errorf("exception while creating %s schema. %v", SchemaName, err)
		logrus.Errorf("NewSQLClient: %v\n", err)
		return nil, err
	}

	_, err = db.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS pg_trgm")
	if err != nil {
		err := fmt.Errorf("exception while creating pg_trgm extention. %v", err)
		logrus.Errorf("NewSQLClient: %v\n", err)
		return nil, err
	}

	_, err = db.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	if err != nil {
		err := fmt.Errorf("exception while creating btree_gin extention. %v", err)
		logrus.Errorf("NewSQLClient: %v\n", err)
		return nil, err
	}

	logrus.Infof("NewSQLClient: New SQL DB CLient Initialized Successfully...\n")
	return db, nil
}

func VerifySQLDatabaseConnection(ctx context.Context, databaseConnection *bun.DB) error {
	logrus.Debugf("VerifySQLDatabaseConnection: Varifying SQL DB CLient...\n")
	if databaseConnection == nil {
		return fmt.Errorf("database connection not initialized")
	}
	return databaseConnection.Ping()
}

func CreateWhereClause(ctx context.Context, whereClause []WhereClauseType) (string, []interface{}, error) {
	var values []interface{}
	var buffer bytes.Buffer

	flag := true
	for i := 0; i < len(whereClause); i++ {
		val := whereClause[i]

		if flag {
			flag = false
		} else {
			buffer.WriteString(" and ")
		}
		var relType string
		if val.RelationType.String() != "" {
			relType = val.RelationType.String()
		} else {
			relType = "="
		}

		// useful column name, so it works in where clauses with joins too
		var fullColumnName string             // placeholder for column name
		fullColumnNameValues := []bun.Ident{} // values of placeholder
		if val.TableAlias == "" {
			fullColumnName = "?TableAlias.?" // ?TableAlias is filled by the bun ORM
		} else {
			fullColumnName = "?.?"
			fullColumnNameValues = append(fullColumnNameValues, bun.Ident(val.TableAlias))
		}
		fullColumnNameValues = append(fullColumnNameValues, bun.Ident(val.ColumnName))

		switch strings.ToLower(relType) {
		case "like":
			buffer.WriteString(fullColumnName)
			buffer.WriteString(val.JsonOperator)
			buffer.WriteString(" ")
			buffer.WriteString(relType)
			buffer.WriteString(" ")
			buffer.WriteString("?")
			colValue, ok := val.ColumnValue.(string)
			if !ok {
				return "", nil, fmt.Errorf("exception while creating where query for tabel %s. Column value not string type", fullColumnName)
			}
			values = append(values, fullColumnNameValues[0])
			if len(fullColumnNameValues) > 1 {
				values = append(values, fullColumnNameValues[1])
			}
			values = append(values, "%"+colValue+"%")
		case "in":
			buffer.WriteString(fullColumnName)
			buffer.WriteString(val.JsonOperator)
			buffer.WriteString(" ")
			buffer.WriteString(relType)
			buffer.WriteString(" ")
			buffer.WriteString("(" + "?" + ")")
			values = append(values, fullColumnNameValues[0])
			if len(fullColumnNameValues) > 1 {
				values = append(values, fullColumnNameValues[1])
			}
			values = append(values, val.ColumnValue)
		case "is":
			buffer.WriteString(fullColumnName)
			buffer.WriteString(val.JsonOperator)
			buffer.WriteString(" ")
			buffer.WriteString(relType)
			buffer.WriteString(" ")
			colValue, ok := val.ColumnValue.(string)
			if !ok {
				return "", nil, fmt.Errorf("exception while creating where query for tabel %s. Column value not string type", fullColumnName)
			}
			if colValue == NullValue || colValue == NotNullValue {
				buffer.WriteString(colValue)
			} else {
				return "", nil, fmt.Errorf("only null and not null values are supported")
			}
			values = append(values, fullColumnNameValues[0])
			if len(fullColumnNameValues) > 1 {
				values = append(values, fullColumnNameValues[1])
			}
		case "any":
			buffer.WriteString("'" + val.ColumnValue.(string) + "'")
			buffer.WriteString(" = ")
			buffer.WriteString(relType)
			buffer.WriteString("(" + fullColumnName + val.JsonOperator + ")")
			values = append(values, fullColumnNameValues[0])
			if len(fullColumnNameValues) > 1 {
				values = append(values, fullColumnNameValues[1])
			}
		default:
			switch val.ColumnValue.(type) {
			case int8, uint8, int16, uint16, int32, uint32, int64, int, uint, uint64, float32, float64:
				buffer.WriteString(fullColumnName)
				buffer.WriteString(val.JsonOperator)
				buffer.WriteString(" ")
				buffer.WriteString(relType)
				buffer.WriteString(" ")

				valPlaceholder := "'" + "?" + "'"
				buffer.WriteString(valPlaceholder)
				values = append(values, fullColumnNameValues[0])
				if len(fullColumnNameValues) > 1 {
					values = append(values, fullColumnNameValues[1])
				}
				values = append(values, val.ColumnValue)

			case string:
				buffer.WriteString(fullColumnName)
				buffer.Write([]byte(val.JsonOperator))
				buffer.WriteString(" ")
				buffer.WriteString(relType)
				buffer.WriteString(" ")
				buffer.WriteString("?")

				colValue, ok := val.ColumnValue.(string)
				if !ok {
					return "", nil, fmt.Errorf("exception while creating where query for tabel %s. Column value not string type", fullColumnName)
				}

				values = append(values, fullColumnNameValues[0])
				if len(fullColumnNameValues) > 1 {
					values = append(values, fullColumnNameValues[1])
				}
				values = append(values, colValue)
			case bool:
				buffer.WriteString(fullColumnName)
				buffer.Write([]byte(val.JsonOperator))
				buffer.WriteString(" ")
				buffer.WriteString(relType)
				buffer.WriteString(" ")
				buffer.WriteString("?")

				values = append(values, fullColumnNameValues[0])
				if len(fullColumnNameValues) > 1 {
					values = append(values, fullColumnNameValues[1])
				}
				values = append(values, val.ColumnValue)
			default:
				buffer.WriteString(fullColumnName)
				buffer.WriteString(val.JsonOperator)
				buffer.WriteString(" ")
				buffer.WriteString(relType)
				buffer.WriteString(" ")
				buffer.WriteString("?")
				values = append(values, fullColumnNameValues[0])
				if len(fullColumnNameValues) > 1 {
					values = append(values, fullColumnNameValues[1])
				}
				values = append(values, val.ColumnValue)
			}
		}
	}
	return buffer.String(), values, nil
}

func PrepareUpdateQuery(ctx context.Context, oldVersion *int, data interface{}, igVersionCheck, colListEmpty bool) *bun.UpdateQuery { // nolint
	q := SqlDBClient.NewUpdate().Model(data).WherePK()
	logrus.Debugf("PrepareUpdateQuery: updateQuery-start: %s", q.String())
	v := reflect.ValueOf(data).Elem()
	for i := 0; i < v.NumField(); i++ {
		valueField := v.Field(i)
		typeField := v.Type().Field(i)
		kindType := valueField.Kind()
		if kindType == reflect.Ptr || typeField.Type.String() == "schema.BaseModel" {
			continue
		}
		pgTag := typeField.Tag.Get("bun")
		customPgTag := typeField.Tag.Get("custom")
		// Expecting first tag as columnName
		columnName := strings.Split(pgTag, ",")[0]
		if valueField.CanSet() {
			if !strings.Contains(customPgTag, "update_invalid") && colListEmpty {
				q = q.Column(strings.TrimSpace(columnName))
			}
			if !igVersionCheck {
				if typeField.Name == "Version" && kindType == reflect.Int {
					oval := valueField.Interface().(int)
					oldVersion = &oval
					newVersionInt64 := int64(*oldVersion + 1)
					valueField.SetInt(newVersionInt64)
					q = q.Where("Version = ?", *oldVersion)
				}
			}
		}
	}
	logrus.Debugf("PrepareUpdateQuery: updateQuery: %s", q.String())
	return q
}

func createColumnList(ctx context.Context, columnList []string) string {
	var buffer bytes.Buffer
	if columnList != nil {
		gblen := len(columnList)
		if gblen >= 1 {
			flag := true
			for i := 0; i < gblen; i++ {
				if flag {
					flag = false
				} else {
					buffer.WriteString(", ")
				}
				buffer.WriteString(strings.ToLower(columnList[i]))
			}
		}
	}
	return buffer.String()
}

func createOrderBy(ctx context.Context, orderByClause []string) (string, error) {
	var buffer bytes.Buffer
	if orderByClause != nil {
		oblen := len(orderByClause)
		if oblen >= 1 {
			flag := true
			for _, value := range orderByClause {
				kv := strings.Split(value, ":")
				if len(kv) != 2 {
					return "", fmt.Errorf("invalid orderBy param %v, it should be of type fieldName:sortingType", value)
				}
				if flag {
					flag = false
				} else {
					buffer.WriteString(", ")
				}
				buffer.WriteString(strings.TrimSpace(kv[0]))
				buffer.WriteString(" ")
				buffer.WriteString(strings.TrimSpace(kv[1]))
			}
		}
	}
	return buffer.String(), nil
}

func ReadUtil(ctx context.Context, tableName string, pagination *Cursor, whereClausefilters []WhereClauseType,
	orderByClause []string, groupByClause, selectedColumns []string, singleRecord bool, result interface{}) (*Cursor, int, error) {

	if err := VerifySQLDatabaseConnection(ctx, SqlDBClient); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. exception while verifying DB connection %v", "Read", tableName, err)
		logrus.Errorf("ReadUtil: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}

	var (
		newOffset     int
		newPagination = new(Cursor)
		readQuery     *bun.SelectQuery
	)

	readQuery = SqlDBClient.NewSelect().Model(result)

	if len(selectedColumns) != 0 {
		colListStr := createColumnList(ctx, selectedColumns)
		if colListStr != "" {
			readQuery = readQuery.ColumnExpr(colListStr)
		}
	}

	if len(whereClausefilters) != 0 {
		queryStr, vals, err := CreateWhereClause(ctx, whereClausefilters)
		if err != nil {
			return nil, http.StatusBadRequest, fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", tableName, err)
		}
		readQuery = readQuery.Where(queryStr, vals...)
	}

	if len(groupByClause) != 0 {
		groupCols := createColumnList(ctx, groupByClause)
		if groupCols != "" {
			readQuery = readQuery.GroupExpr(groupCols)
		}
	}

	if len(orderByClause) != 0 {
		orderByCols, err := createOrderBy(ctx, orderByClause)
		if err != nil {
			return nil, http.StatusBadRequest, fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", tableName, err)
		}
		if orderByCols != "" {
			readQuery = readQuery.OrderExpr(orderByCols)
		}
	}

	if !singleRecord && pagination != nil {
		if pagination.PageSize <= 0 {
			return nil, http.StatusBadRequest, fmt.Errorf("unable to Perform %s Operation on Table: %s. Invalid PazeSize %v", "Read",
				tableName, pagination.PageSize)
		}

		if pagination.PageSize != 0 {
			var offset int
			var cErr error
			if pagination.PageToken == "" {
				offset = 0
				if pagination.PageNum > 0 {
					offset = (pagination.PageNum - 1) * pagination.PageSize
				}
			} else {
				offset, cErr = strconv.Atoi(pagination.PageToken)
				if cErr != nil {
					return nil, http.StatusInternalServerError, fmt.Errorf("unable to Perform %s Operation on Table: %s. Exception while converting pagetoken to offset. %v", "Read",
						tableName, cErr)
				}
			}
			readQuery = readQuery.Limit(pagination.PageSize).Offset(offset)
			newOffset = pagination.PageSize + offset
			newPagination.PageNum = offset/pagination.PageSize + 1
			newPagination.PageSize = pagination.PageSize
			newPagination.PageToken = strconv.Itoa(newOffset)
		}
	}

	if count, err := readQuery.ScanAndCount(ctx); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("unable to Perform %s Operation on Table: %s. Exception while reading data. %v", "Read",
			tableName, err)
	} else if !singleRecord && pagination != nil {
		if newOffset >= count {
			newPagination.PageToken = ""
		}
		newPagination.TotalRecords = uint32(count)
		if newPagination.PageSize > 0 {
			if count%newPagination.PageSize != 0 {
				newPagination.TotalPages = uint32(count/newPagination.PageSize) + 1
			} else {
				newPagination.TotalPages = uint32(count / newPagination.PageSize)
			}
		}
		return newPagination, http.StatusOK, nil
	}
	return nil, http.StatusOK, nil
}

func CreateIndexCreateQuery(ctx context.Context, indexCreateParam IndexParams) string {
	var buffer bytes.Buffer

	if strings.ToLower(indexCreateParam.Type) == "unique" {
		buffer.WriteString("create unique index if not exists ")
	} else {
		buffer.WriteString("create index if not exists ")
	}
	buffer.WriteString(indexCreateParam.Name)
	buffer.WriteString(" on ")
	buffer.WriteString(indexCreateParam.TableName)
	if strings.ToLower(indexCreateParam.Type) != "unique" && indexCreateParam.Type != "" {
		buffer.WriteString(" using ")
		buffer.WriteString(indexCreateParam.Type)
	}
	buffer.WriteString(" ")
	buffer.WriteString("(")
	prepColumnList := strings.Join(indexCreateParam.ColumnNames, ",")
	buffer.WriteString(prepColumnList)
	buffer.WriteString(")")

	return buffer.String()
}
