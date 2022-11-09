package gen

import (
	"fmt"
	"gorm.io/gen/pkg/generate"
	model2 "gorm.io/gen/pkg/model"
	"reflect"
	"regexp"
	"strings"

	"gorm.io/gen/field"
	"gorm.io/gorm/schema"
)

// ModelOpt field option
type ModelOpt = model2.Option

var ns = schema.NamingStrategy{}

var (
	// FieldNew add new field (any type your want)
	FieldNew = func(fieldName, fieldType, fieldTag string) model2.CreateFieldOpt {
		return func(*model2.Field) *model2.Field {
			return &model2.Field{
				Name:         fieldName,
				Type:         fieldType,
				OverwriteTag: fieldTag,
			}
		}
	}
	// FieldIgnore ignore some columns by name
	FieldIgnore = func(columnNames ...string) model2.FilterFieldOpt {
		return func(m *model2.Field) *model2.Field {
			for _, name := range columnNames {
				if m.ColumnName == name {
					return nil
				}
			}
			return m
		}
	}
	// FieldIgnoreReg ignore some columns by RegExp
	FieldIgnoreReg = func(columnNameRegs ...string) model2.FilterFieldOpt {
		regs := make([]regexp.Regexp, len(columnNameRegs))
		for i, reg := range columnNameRegs {
			regs[i] = *regexp.MustCompile(reg)
		}
		return func(m *model2.Field) *model2.Field {
			for _, reg := range regs {
				if reg.MatchString(m.ColumnName) {
					return nil
				}
			}
			return m
		}
	}
	// FieldRename specify field name in generated struct
	FieldRename = func(columnName string, newName string) model2.ModifyFieldOpt {
		return func(m *model2.Field) *model2.Field {
			if m.ColumnName == columnName {
				m.Name = newName
			}
			return m
		}
	}
	// FieldComment specify field comment in generated struct
	FieldComment = func(columnName string, comment string) model2.ModifyFieldOpt {
		return func(m *model2.Field) *model2.Field {
			if m.ColumnName == columnName {
				m.ColumnComment = comment
				m.MultilineComment = strings.Contains(comment, "\n")
			}
			return m
		}
	}
	// FieldType specify field type in generated struct
	FieldType = func(columnName string, newType string) model2.ModifyFieldOpt {
		return func(m *model2.Field) *model2.Field {
			if m.ColumnName == columnName {
				m.Type = newType
			}
			return m
		}
	}
	// FieldTypeReg specify field type in generated struct by RegExp
	FieldTypeReg = func(columnNameReg string, newType string) model2.ModifyFieldOpt {
		reg := regexp.MustCompile(columnNameReg)
		return func(m *model2.Field) *model2.Field {
			if reg.MatchString(m.ColumnName) {
				m.Type = newType
			}
			return m
		}
	}
	// FieldGenType specify field gen type in generated dao
	FieldGenType = func(columnName string, newType string) model2.ModifyFieldOpt {
		return func(m *model2.Field) *model2.Field {
			if m.ColumnName == columnName {
				m.CustomGenType = newType
			}
			return m
		}
	}
	// FieldGenTypeReg specify field gen type in generated dao  by RegExp
	FieldGenTypeReg = func(columnNameReg string, newType string) model2.ModifyFieldOpt {
		reg := regexp.MustCompile(columnNameReg)
		return func(m *model2.Field) *model2.Field {
			if reg.MatchString(m.ColumnName) {
				m.CustomGenType = newType
			}
			return m
		}
	}
	// FieldTag specify GORM tag and JSON tag
	FieldTag = func(columnName string, gormTag, jsonTag string) model2.ModifyFieldOpt {
		return func(m *model2.Field) *model2.Field {
			if m.ColumnName == columnName {
				m.GORMTag, m.JSONTag = gormTag, jsonTag
			}
			return m
		}
	}
	// FieldJSONTag specify JSON tag
	FieldJSONTag = func(columnName string, jsonTag string) model2.ModifyFieldOpt {
		return func(m *model2.Field) *model2.Field {
			if m.ColumnName == columnName {
				m.JSONTag = jsonTag
			}
			return m
		}
	}
	// FieldJSONTagWithNS specify JSON tag with name strategy
	FieldJSONTagWithNS = func(schemaName func(columnName string) (tagContent string)) model2.ModifyFieldOpt {
		return func(m *model2.Field) *model2.Field {
			if schemaName != nil {
				m.JSONTag = schemaName(m.ColumnName)
			}
			return m
		}
	}
	// FieldWithNS specify JSON tag with name strategy
	FieldWithNS = func(fieldFunc func(param *model2.Field) (ret *model2.Field)) model2.ModifyFieldOpt {
		return func(m *model2.Field) *model2.Field {
			if fieldFunc != nil {
				m = fieldFunc(m)
			}
			return m
		}
	}
	// FieldGORMTag specify GORM tag
	FieldGORMTag = func(columnName string, gormTag string) model2.ModifyFieldOpt {
		return func(m *model2.Field) *model2.Field {
			if m.ColumnName == columnName {
				m.GORMTag = gormTag
			}
			return m
		}
	}
	// FieldNewTag add new tag
	FieldNewTag = func(columnName string, newTag string) model2.ModifyFieldOpt {
		return func(m *model2.Field) *model2.Field {
			if m.ColumnName == columnName {
				m.NewTag += " " + newTag
			}
			return m
		}
	}
	// FieldNewTagWithNS add new tag with name strategy
	FieldNewTagWithNS = func(tagName string, schemaName func(columnName string) string) model2.ModifyFieldOpt {
		return func(m *model2.Field) *model2.Field {
			if schemaName == nil {
				schemaName = func(name string) string { return name }
			}
			m.NewTag = fmt.Sprintf(`%s %s:"%s"`, m.NewTag, tagName, schemaName(m.ColumnName))
			return m
		}
	}
	// FieldTrimPrefix trim column name's prefix
	FieldTrimPrefix = func(prefix string) model2.ModifyFieldOpt {
		return func(m *model2.Field) *model2.Field {
			m.Name = strings.TrimPrefix(m.Name, prefix)
			return m
		}
	}
	// FieldTrimSuffix trim column name's suffix
	FieldTrimSuffix = func(suffix string) model2.ModifyFieldOpt {
		return func(m *model2.Field) *model2.Field {
			m.Name = strings.TrimSuffix(m.Name, suffix)
			return m
		}
	}
	// FieldAddPrefix add prefix to struct's memeber name
	FieldAddPrefix = func(prefix string) model2.ModifyFieldOpt {
		return func(m *model2.Field) *model2.Field {
			m.Name = prefix + m.Name
			return m
		}
	}
	// FieldAddSuffix add suffix to struct's memeber name
	FieldAddSuffix = func(suffix string) model2.ModifyFieldOpt {
		return func(m *model2.Field) *model2.Field {
			m.Name += suffix
			return m
		}
	}
	// FieldRelate relate to table in database
	FieldRelate = func(relationship field.RelationshipType, fieldName string, table *generate.QueryStructMeta, config *field.RelateConfig) model2.CreateFieldOpt {
		if config == nil {
			config = &field.RelateConfig{}
		}
		if config.JSONTag == "" {
			config.JSONTag = ns.ColumnName("", fieldName)
		}
		return func(*model2.Field) *model2.Field {
			return &model2.Field{
				Name:         fieldName,
				Type:         config.RelateFieldPrefix(relationship) + table.StructInfo.Type,
				JSONTag:      config.JSONTag,
				GORMTag:      config.GORMTag,
				NewTag:       config.NewTag,
				OverwriteTag: config.OverwriteTag,

				Relation: field.NewRelationWithType(
					relationship, fieldName, table.StructInfo.Package+"."+table.StructInfo.Type,
					table.Relations()...),
			}
		}
	}
	// FieldRelateModel relate to exist table model
	FieldRelateModel = func(relationship field.RelationshipType, fieldName string, relModel interface{}, config *field.RelateConfig) model2.CreateFieldOpt {
		st := reflect.TypeOf(relModel)
		if st.Kind() == reflect.Ptr {
			st = st.Elem()
		}
		fieldType := st.String()

		if config == nil {
			config = &field.RelateConfig{}
		}
		if config.JSONTag == "" {
			config.JSONTag = ns.ColumnName("", fieldName)
		}

		return func(*model2.Field) *model2.Field {
			return &model2.Field{
				Name:         fieldName,
				Type:         config.RelateFieldPrefix(relationship) + fieldType,
				JSONTag:      config.JSONTag,
				GORMTag:      config.GORMTag,
				NewTag:       config.NewTag,
				OverwriteTag: config.OverwriteTag,

				Relation: field.NewRelationWithModel(relationship, fieldName, fieldType, relModel),
			}
		}
	}

	// WithMethod add custom method for table model
	WithMethod = func(methods ...interface{}) model2.AddMethodOpt {
		return func() []interface{} { return methods }
	}
)

var (
	DefaultMethodTableWithNamer = (&defaultModel{}).TableName
)

type defaultModel struct {
}

func (*defaultModel) TableName(namer schema.Namer) string {
	if namer == nil {
		return "@@table"
	}
	return namer.TableName("@@table")
}
