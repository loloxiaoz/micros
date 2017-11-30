package orm

import (
	"fmt"
	"strings"
)

// Define callbacks for updating
func init() {
	DefaultCallback.Update().Register("micros:update", updateCallback)
}

// updateCallback the callback used to update data to database
func updateCallback(scope *Scope) {
	if !scope.HasError() {
		var sqls []string

		if updateAttrs, ok := scope.InstanceGet("micros:update_attrs"); ok {
			for column, value := range updateAttrs.(map[string]interface{}) {
				sqls = append(sqls, fmt.Sprintf("%v = %v", scope.Quote(column), scope.AddToVars(value)))
			}
		} else {
			for _, field := range scope.Fields() {
				if scope.changeableField(field) {
					if !field.IsPrimaryKey && field.IsNormal {
						sqls = append(sqls, fmt.Sprintf("%v = %v", scope.Quote(field.DBName), scope.AddToVars(field.Field.Interface())))
					}
				}
			}
		}

		var extraOption string
		if str, ok := scope.Get("micros:update_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		if len(sqls) > 0 {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v%v%v",
				scope.QuotedTableName(),
				strings.Join(sqls, ", "),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}
