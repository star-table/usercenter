package mysql

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/opentracing/opentracing-go"
	"github.com/qustavo/sqlhooks/v2"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/threadlocal"
	"github.com/star-table/usercenter/pkg/tracing"

	"upper.io/db.v3/lib/sqlbuilder"
	upper "upper.io/db.v3/mysql"

	"strconv"
)

var mysqlMutex sync.Mutex
var sess sqlbuilder.Database

type Config struct {
	Host         string
	Port         int
	User         string
	Pwd          string
	Database     string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  int
}

type Client struct {
	config *Config
}

func NewClient(config Config) *Client {
	return &Client{
		config: &config,
	}
}

// Hooks satisfies the sqlhook.Hooks interface
type Hooks struct{}

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	if tracing.EnableTracing() {
		if v, ok := threadlocal.Mgr.GetValue(consts.JaegerContextSpanKey); ok {
			if parentSpan, ok := v.(opentracing.Span); ok {
				spanCtx := parentSpan.Context()
				span := tracing.StartSpan("mysql opt", opentracing.ChildOf(spanCtx))
				span.SetTag("sql", query)
				span.SetTag("args", args)
				span.SetTag("operation", "mysql opt")
				return context.WithValue(ctx, "traceSpan", span), nil
			}
		}
	}
	return ctx, nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	if tracing.EnableTracing() {
		if v := ctx.Value("traceSpan"); v != nil {
			if span, ok := v.(opentracing.Span); ok {
				span.Finish()
			}
		}
	}
	return ctx, nil
}

func init() {
	sql.Register("mysql-hooks", sqlhooks.Wrap(&mysql.MySQLDriver{}, &Hooks{}))
}

func (c Client) GetConnect() (sqlbuilder.Database, error) {
	if c.config == nil {
		panic(errors.New("Mysql Datasource Configuration is missing!"))
	}
	if sess == nil {
		mysqlMutex.Lock()
		defer mysqlMutex.Unlock()
		if sess == nil {
			var err error
			sess, err = c.initSess()
			if err != nil {
				return nil, err
			}
		}
	}
	if err := sess.Ping(); err != nil {
		sess, err = c.initSess()
		if err != nil {
			return nil, err
		}
	}
	return sess, nil
}

func (c Client) initSess() (sqlbuilder.Database, error) {

	conf := c.config
	settings := &upper.ConnectionURL{
		User:     conf.User,
		Password: conf.Pwd,
		Database: conf.Database,
		Host:     conf.Host + ":" + strconv.Itoa(conf.Port),
		Socket:   "",
		Options: map[string]string{
			"parseTime": "true",
			"loc":       "Local",
			"charset":   "utf8mb4",
			"collation": "utf8mb4_unicode_ci",
		},
	}

	sess, err := upper.Open(settings)
	if err != nil {
		return nil, err
	}
	maxOpenConns := 50
	maxIdleConns := 10
	maxLifetime := 300
	if c.config.MaxOpenConns > 0 {
		maxOpenConns = c.config.MaxOpenConns
	}
	if c.config.MaxIdleConns > 0 {
		maxIdleConns = c.config.MaxIdleConns
	}
	if c.config.MaxLifetime > 0 {
		maxLifetime = c.config.MaxLifetime
	}
	sess.SetMaxOpenConns(maxOpenConns)
	sess.SetMaxIdleConns(maxIdleConns)
	sess.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)
	return sess, nil
}

type Domain interface {
	TableName() string
}

func CloseTx(conn sqlbuilder.Database, tx sqlbuilder.Tx) {
	if tx != nil {
		if err := tx.Close(); err != nil {
			Info(err)
		}
	}
}

func Close(conn sqlbuilder.Database) {
	if conn != nil {
		if err := conn.Close(); err != nil {
			Info(err)
		}
	}
}

func Rollback(tx sqlbuilder.Tx) {
	err := tx.Rollback()
	if err != nil {
		Info(err)
	}
}

type Upd map[string]interface{}
