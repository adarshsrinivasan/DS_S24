package sql

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/common"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/db"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/oiime/logrusbun"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type clientObj struct {
	tx        *sql.Tx
	bunClient *bun.DB
}

const (
	PostgresHostEnv     = "POSTGRES_HOST"
	PostgresPortEnv     = "POSTGRES_PORT"
	PostgresUsernameEnv = "POSTGRES_USERNAME"
	PostgresPasswordEnv = "POSTGRES_PASSWORD"
	PostgresDbEnv       = "POSTGRES_DB"
	PostgresMaxConnEnv  = "POSTGRES_MAX_CONN"

	DefaultIdleTimeouts = -1  // never timeout/close an idle connection.
	DefaultLogSlowQuery = 100 // log db queries slower than 100ms by default

	DialTimeOut  = 5 * time.Second // time to wait for connection before exit
	ReadTimeOut  = 5 * time.Second // time to wait for read to complete before exit
	WriteTimeOut = 5 * time.Second // time to wait for write to complete before exit

)

var (
	poolObj *connPool
)

func getClient(ctx context.Context, applicationName, schemaName string) (*clientObj, error) {
	maxConn, _ := strconv.Atoi(common.GetEnv(PostgresMaxConnEnv, "500"))
	if poolObj == nil {
		poolObj = &connPool{}
		if err := poolObj.initialize(ctx, applicationName, schemaName, maxConn); err != nil {
			err = fmt.Errorf("exception while initializing SQL connection pool. %v", err)
			logrus.Errorf("NewClient: %v\n", err)
			return nil, err
		}
	}
	return poolObj.getClient(ctx), nil
}

func NewClient(ctx context.Context, applicationName, schemaName string) (*clientObj, error) {
	host := common.GetEnv(PostgresHostEnv, "localhost")
	port := common.GetEnv(PostgresPortEnv, "5432")
	username := common.GetEnv(PostgresUsernameEnv, "admin")
	password := common.GetEnv(PostgresPasswordEnv, "admin")
	dbName := common.GetEnv(PostgresDbEnv, "marketplace")
	maxConn, _ := strconv.Atoi(common.GetEnv(PostgresMaxConnEnv, "300"))

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(fmt.Sprintf("postgres://%s:@%s:%s/%s?sslmode=disable", username, host, port, dbName)),
		pgdriver.WithPassword(password),
		pgdriver.WithApplicationName(applicationName),
		pgdriver.WithDialTimeout(DialTimeOut),
		pgdriver.WithReadTimeout(ReadTimeOut),
		pgdriver.WithWriteTimeout(WriteTimeOut),
		pgdriver.WithConnParams(map[string]interface{}{
			"search_path": fmt.Sprintf("%s", schemaName),
		})))
	sqldb.SetConnMaxIdleTime(DefaultIdleTimeouts)
	sqldb.SetMaxIdleConns(maxConn)
	sqldb.SetMaxOpenConns(maxConn)

	_, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName))
	if err != nil {
		err = fmt.Errorf("exception while opening postgres connection: %v", err)
		logrus.Errorf("NewClient: %v\n", err)
		return nil, err
	}
	bunDBObj := bun.NewDB(sqldb, pgdialect.New())
	logrusObj := logrus.New()
	logrusObj.SetFormatter(&logrus.TextFormatter{DisableQuote: true})
	bunDBObj.AddQueryHook(logrusbun.NewQueryHook(logrusbun.QueryHookOptions{
		Logger:          logrusObj,
		LogSlow:         DefaultLogSlowQuery,
		QueryLevel:      logrus.DebugLevel,
		SlowLevel:       logrus.WarnLevel,
		ErrorLevel:      logrus.ErrorLevel,
		MessageTemplate: "{{.Operation}}[{{.Duration}}]: {{.Query}}",
		ErrorTemplate:   "{{.Operation}}[{{.Duration}}]: {{.Query}}: {{.Error}}",
	}))
	client := &clientObj{
		bunClient: bunDBObj,
	}

	return client, nil
}

func (client *clientObj) Initialize(ctx context.Context, schemaName string) error {
	if err := client.VerifyConnection(ctx); err != nil {
		err := fmt.Errorf("exception while verifying SQL DB connection. %v", err)
		logrus.Errorf("Initialize: %v\n", err)
		return err
	}

	if _, err := client.bunClient.ExecContext(ctx, "CREATE SCHEMA IF NOT EXISTS ?", bun.Ident(schemaName)); err != nil {
		err := fmt.Errorf("exception while creating %s schema. %v", schemaName, err)
		logrus.Errorf("Initialize: %v\n", err)
		return err
	}

	if _, err := client.bunClient.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS pg_trgm"); err != nil {
		err := fmt.Errorf("exception while creating pg_trgm extention. %v", err)
		logrus.Errorf("Initialize: %v\n", err)
		return err
	}

	if _, err := client.bunClient.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin"); err != nil {
		err := fmt.Errorf("exception while creating btree_gin extention. %v", err)
		logrus.Errorf("Initialize: %v\n", err)
		return err
	}

	logrus.Infof("Initialize: SQL DB CLient Initialized Successfully...\n")
	return nil
}

func (client *clientObj) VerifyConnection(ctx context.Context) error {
	logrus.Debugf("VerifyConnection: Varifying SQL DB CLient...\n")
	if client.bunClient == nil {
		return fmt.Errorf("database connection not initialized")
	}
	return client.bunClient.Ping()
}

func (client *clientObj) CreateTable(ctx context.Context, model interface{}, tableName string, foreignKeys []db.ForeignKey) error {
	createTableQuery := client.bunClient.NewCreateTable().
		Model(model).
		IfNotExists()

	for _, fk := range foreignKeys {
		createTableQuery.ForeignKey(client.prepareForeignKeyQuery(fk))
	}

	_, err := createTableQuery.Exec(ctx)
	if err != nil {
		err := fmt.Errorf("exception while creaiting event table %s. %v", err, tableName)
		logrus.Errorf("CreateTable: %v\n", err)
		return err
	}
	return nil
}

func (client *clientObj) Insert(ctx context.Context, model interface{}, tableName string) error {
	if _, err := client.bunClient.NewInsert().Model(model).Exec(ctx); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", tableName, err)
		logrus.Errorf("InsertOne: %v\n", err)
		return err
	}
	return nil
}

func (client *clientObj) Read(ctx context.Context, tableName string, pagination *db.Cursor, whereClauseFilters []db.WhereClauseType,
	orderByClause []string, groupByClause, selectedColumns []string, singleRecord bool, result interface{}) (*db.Cursor, error) {

	var (
		newOffset     int
		newPagination = new(db.Cursor)
		readQuery     *bun.SelectQuery
	)

	readQuery = client.bunClient.NewSelect().Model(result)

	if len(selectedColumns) != 0 {
		colListStr := client.createColumnList(ctx, selectedColumns)
		if colListStr != "" {
			readQuery = readQuery.ColumnExpr(colListStr)
		}
	}

	if len(whereClauseFilters) != 0 {
		queryStr, vals, err := client.createWhereClause(ctx, whereClauseFilters)
		if err != nil {
			return nil, fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", tableName, err)
		}
		readQuery = readQuery.Where(queryStr, vals...)
	}

	if len(groupByClause) != 0 {
		groupCols := client.createColumnList(ctx, groupByClause)
		if groupCols != "" {
			readQuery = readQuery.GroupExpr(groupCols)
		}
	}

	if len(orderByClause) != 0 {
		orderByCols, err := client.createOrderBy(ctx, orderByClause)
		if err != nil {
			return nil, fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", tableName, err)
		}
		if orderByCols != "" {
			readQuery = readQuery.OrderExpr(orderByCols)
		}
	}

	if !singleRecord && pagination != nil {
		if pagination.PageSize <= 0 {
			return nil, fmt.Errorf("unable to Perform %s Operation on Table: %s. Invalid PazeSize %v", "Read",
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
					return nil, fmt.Errorf("unable to Perform %s Operation on Table: %s. Exception while converting pagetoken to offset. %v", "Read",
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
		return nil, fmt.Errorf("unable to Perform %s Operation on Table: %s. Exception while reading data. %v", "Read",
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
		return newPagination, nil
	}
	return nil, nil
}

func (client *clientObj) Update(ctx context.Context, model interface{}, tableName string, igVersionCheck bool) error {
	updateQuery := client.prepareUpdateQuery(ctx, model, false)
	_, err := updateQuery.Exec(ctx)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", tableName, err)
		logrus.Errorf("Update: %v\n", err)
		return err
	}
	return nil
}

func (client *clientObj) Delete(ctx context.Context, model interface{}, tableName string, whereClauseFilters []db.WhereClauseType) error {
	deleteQuery := client.bunClient.NewDelete().
		Model(model)

	// prepare whereClause.
	queryStr, vals, err := client.createWhereClause(ctx, whereClauseFilters)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", tableName, err)
		logrus.Errorf("Delete: %v\n", err)
		return err
	}
	deleteQuery = deleteQuery.Where(queryStr, vals...)

	if _, err := deleteQuery.Exec(ctx); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", tableName, err)
		logrus.Errorf("Delete: %v\n", err)
		return err
	}
	return nil
}

func (client *clientObj) Close(ctx context.Context) error {
	//poolObj.close(ctx, client)
	return client.bunClient.Close()
}

func (client *clientObj) prepareForeignKeyQuery(foreignKeyObj db.ForeignKey) string {
	query := fmt.Sprintf("(\"%s\") REFERENCES \"%s\" (\"%s\")", foreignKeyObj.ColumnName, foreignKeyObj.SrcTableName, foreignKeyObj.SrcColumnName)

	if foreignKeyObj.CascadeDelete {
		query += " ON DELETE CASCADE"
	}
	return query
}

func (client *clientObj) prepareUpdateQuery(ctx context.Context, data interface{}, igVersionCheck bool) *bun.UpdateQuery { // nolint
	q := client.bunClient.NewUpdate().Model(data).WherePK()
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
			if !strings.Contains(customPgTag, "update_invalid") {
				q = q.Column(strings.TrimSpace(columnName))
			}
			if !igVersionCheck {
				if typeField.Name == "Version" && kindType == reflect.Int {
					oldVersion := valueField.Interface().(int)
					newVersionInt64 := int64(oldVersion + 1)
					valueField.SetInt(newVersionInt64)
					q = q.Where("Version = ?", oldVersion)
				}
			}
		}
	}
	logrus.Debugf("PrepareUpdateQuery: updateQuery: %s", q.String())
	return q
}

func (client *clientObj) createWhereClause(ctx context.Context, whereClause []db.WhereClauseType) (string, []interface{}, error) {
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
			if colValue == db.NullValue || colValue == db.NotNullValue {
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

func (client *clientObj) createColumnList(ctx context.Context, columnList []string) string {
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

func (client *clientObj) createOrderBy(ctx context.Context, orderByClause []string) (string, error) {
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

func (client *clientObj) createIndexCreateQuery(ctx context.Context, indexCreateParam db.IndexParams) string {
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
