package gormspanner

import (
	"context"
	"database/sql"
	"fmt"
	spannergorm "github.com/googleapis/go-gorm-spanner"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/errorhandler"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"os"
	"time"
)

func NewCamelCaseNamer() *CamelCaseNamer {
	return &CamelCaseNamer{}
}

type CamelCaseNamer struct{}

func (n *CamelCaseNamer) TableName(table string) string {
	return table
}

func (n *CamelCaseNamer) SchemaName(table string) string {
	return table
}

func (n *CamelCaseNamer) ColumnName(table, column string) string {
	return column
}

func (n *CamelCaseNamer) JoinTableName(joinTable string) string {
	return joinTable
}

func (n *CamelCaseNamer) RelationshipFKName(rel schema.Relationship) string {
	return rel.Name
}

func (n *CamelCaseNamer) CheckerName(table, column string) string {
	return table + "_" + column
}

func (n *CamelCaseNamer) IndexName(table, column string) string {
	return table + column + "Idx"
}

func (n *CamelCaseNamer) UniqueName(table, column string) string {
	return table + column + "UIdx"
}

type ISession interface {
	Association(column string) *gorm.Association
	Model(value interface{}) (tx *gorm.DB)
	Clauses(conds ...clause.Expression) (tx *gorm.DB)
	Table(name string, args ...interface{}) (tx *gorm.DB)
	Distinct(args ...interface{}) (tx *gorm.DB)
	Select(query interface{}, args ...interface{}) (tx *gorm.DB)
	Omit(columns ...string) (tx *gorm.DB)
	Where(query interface{}, args ...interface{}) (tx *gorm.DB)
	Not(query interface{}, args ...interface{}) (tx *gorm.DB)
	Or(query interface{}, args ...interface{}) (tx *gorm.DB)
	Joins(query string, args ...interface{}) (tx *gorm.DB)
	InnerJoins(query string, args ...interface{}) (tx *gorm.DB)
	Group(name string) (tx *gorm.DB)
	Having(query interface{}, args ...interface{}) (tx *gorm.DB)
	Order(value interface{}) (tx *gorm.DB)
	Limit(limit int) (tx *gorm.DB)
	Offset(offset int) (tx *gorm.DB)
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) (tx *gorm.DB)
	Preload(query string, args ...interface{}) (tx *gorm.DB)
	Attrs(attrs ...interface{}) (tx *gorm.DB)
	Assign(attrs ...interface{}) (tx *gorm.DB)
	Unscoped() (tx *gorm.DB)
	Raw(sql string, values ...interface{}) (tx *gorm.DB)
	Create(value interface{}) (tx *gorm.DB)
	CreateInBatches(value interface{}, batchSize int) (tx *gorm.DB)
	Save(value interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Take(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Last(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) *gorm.DB
	FirstOrInit(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	FirstOrCreate(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Update(column string, value interface{}) (tx *gorm.DB)
	Updates(values interface{}) (tx *gorm.DB)
	UpdateColumn(column string, value interface{}) (tx *gorm.DB)
	UpdateColumns(values interface{}) (tx *gorm.DB)
	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
	Count(count *int64) (tx *gorm.DB)
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	Scan(dest interface{}) (tx *gorm.DB)
	Pluck(column string, dest interface{}) (tx *gorm.DB)
	ScanRows(rows *sql.Rows, dest interface{}) error
	Connection(fc func(tx *gorm.DB) error) (err error)
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
	Begin(opts ...*sql.TxOptions) *gorm.DB
	Commit() *gorm.DB
	Rollback() *gorm.DB
	SavePoint(name string) *gorm.DB
	RollbackTo(name string) *gorm.DB
	Exec(sql string, values ...interface{}) (tx *gorm.DB)
	Session(config *gorm.Session) *gorm.DB
	WithContext(ctx context.Context) *gorm.DB
	Debug() (tx *gorm.DB)
	Set(key string, value interface{}) *gorm.DB
	Get(key string) (interface{}, bool)
	InstanceSet(key string, value interface{}) *gorm.DB
	InstanceGet(key string) (interface{}, bool)
	AddError(err error) error
	DB() (*sql.DB, error)
	SetupJoinTable(model interface{}, field string, joinTable interface{}) error
	Use(plugin gorm.Plugin) error
	ToSQL(queryFn func(tx *gorm.DB) *gorm.DB) string
	Migrator() gorm.Migrator
	AutoMigrate(dst ...interface{}) error
}

// NewSession method
var NewSession = func(ctx context.Context, opt config.Spanner) (*gorm.DB, error) {
	if opt.WithoutAuthentication && opt.GRPCInsecure {
		if err := os.Setenv("SPANNER_EMULATOR_HOST", opt.EndPoint); err != nil {
			return nil, err
		}
	}
	return gorm.Open(spannergorm.New(spannergorm.Config{
		DriverName: "spanner",
		DSN:        fmt.Sprintf("projects/%s/instances/%s/databases/%s", opt.ProjectID, opt.Instance, opt.Database),
	}), &gorm.Config{
		NamingStrategy:                           NewCamelCaseNamer(),
		NowFunc:                                  func() time.Time { return time.Now().UTC() },
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: false,
	})
}

func NewExtendGormSpannerDatabase(logger *zap.Logger, opt config.Spanner) (*gorm.DB, func(), error) {
	session, err := NewSession(context.Background(), opt)
	if err != nil {
		return nil, nil, err
	}
	logger.Info("Gorm Spanner init complete", zap.String("system", "Database"))

	return session,
		func() {
			defer errorhandler.PanicErrorHandler("Core", "")
			sqlDB, err := session.DB()
			if err != nil {
				panic(errorhandler.NewErrDBDisconnection(err))
			}
			if err := sqlDB.Close(); err != nil {
				panic(errorhandler.NewErrDBDisconnection(err))
			}
		}, nil
}
