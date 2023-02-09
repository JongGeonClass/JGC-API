package dbmodel

// JGC에서 사용할 테이블들입니다.
var tables []interface{} = nil
var tableNames []string = nil

// 새로운 테이블을 등록합니다.
func AddTable(tableName string, table interface{}) {
	tables = append(tables, table)
	tableNames = append(tableNames, tableName)
}

// 등록된 모든 테이블을 가져옵니다.
func GetTables() ([]interface{}, []string) {
	return tables, tableNames
}
