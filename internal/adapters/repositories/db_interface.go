package repositories

// DBClient interface genérica para operações de banco de dados
type DBClient interface {
	Create(value interface{}) error
	First(dest interface{}, conds ...interface{}) error
	Find(dest interface{}, conds ...interface{}) error
	Where(query interface{}, args ...interface{}) DBClient
	Delete(value interface{}, conds ...interface{}) error
	AutoMigrate(dst ...interface{}) error
}