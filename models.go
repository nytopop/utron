package utron

import (
	"errors"
	"reflect"

	mgo "gopkg.in/mgo.v2"

	// support mysql, sqlite3 and postgresql
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)

// Model facilitate database interactions, supports postgres, mysql and foundation
type Model struct {
	models map[string]reflect.Value
	isOpen bool
	DB     *mgo.Database
}

// NewModelWithConfig creates a new model, and opens database connection based on cfg settings
func NewModelWithConfig(cfg *Config) (*Model, error) {
	m := NewModel()
	if err := m.OpenWithConfig(cfg); err != nil {
		return nil, err
	}
	return m, nil
}

// NewModel returns a new Model without opening database connection
func NewModel() *Model {
	return &Model{
		models: make(map[string]reflect.Value),
	}
}

// IsOpen returns true if the Model has already established connection
// to the database
func (m *Model) IsOpen() bool {
	return m.isOpen
}

// OpenWithConfig opens database connection with the settings found in cfg
func (m *Model) OpenWithConfig(cfg *Config) error {
	session, err := mgo.Dial(cfg.DatabaseConn)
	if err != nil {
		return err
	}
	db := session.DB(cfg.Database)
	m.DB = db
	m.isOpen = true
	return nil
}

// Register adds the values to the models registry
func (m *Model) Register(values ...interface{}) error {

	// do not work on them.models first, this is like an insurance policy
	// whenever we encounter any error in the values nothing goes into the registry
	models := make(map[string]reflect.Value)
	if len(values) > 0 {
		for _, val := range values {
			rVal := reflect.ValueOf(val)
			if rVal.Kind() == reflect.Ptr {
				rVal = rVal.Elem()
			}
			switch rVal.Kind() {
			case reflect.Struct:
				models[getTypName(rVal.Type())] = reflect.New(rVal.Type())
			default:
				return errors.New("utron: models must be structs")
			}
		}
	}
	for k, v := range models {
		m.models[k] = v
	}
	return nil
}
