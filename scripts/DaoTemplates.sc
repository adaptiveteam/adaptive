import $file.Meta
import Meta._

import $file.Templates
import Templates._

import $file.GoTemplates
import GoTemplates._

import $file.Dsl
import Dsl._

import $file.DynamoTemplates
import DynamoTemplates._

import $file.VirtualFields
import VirtualFields._

def idArgList(table: Table): String =
	table.entity.primaryKeyFields.map(fieldArg).mkString(", ")

def indexArgList(index: Index): String = fieldArg(index.hashKey) + (
    index.rangeKey.map(fieldArg).map(", " + _).getOrElse("")
)

def entityArg(entity: Entity): String = arg(entity.parentField)

def entitySliceArg(entity: Entity): String = argSlice(entity.parentField)

def interfaceTemplate(dao: Dao): List[String] = {
	val table = dao.table
	val idArgs = idArgList(table)
	val entityArgValue = entityArg(table.entity)
	val entityArgSliceValue = entitySliceArg(table.entity)
	val supportsDeactivation = dao.table.entity.supports(DeactivationTrait)

	def interfaceOperationTemplate(operation: DaoOperation): List[String] = {
		operation match {
			case DaoCreateRow => List(
				s"Create($entityArgValue) error",
				s"CreateUnsafe($entityArgValue)",
			)
			case DaoReadRow => List(
				s"Read($idArgs) ($entityArgValue, err error)",
				s"ReadUnsafe($idArgs) ($entityArgValue)",
			)
			case DaoReadOrEmptyRow => List(
				s"ReadOrEmpty($idArgs) ($entityArgSliceValue, err error)",
				s"ReadOrEmptyUnsafe($idArgs) ($entityArgSliceValue)",
			)
			case DaoUpdateRow => List(
				s"CreateOrUpdate($entityArgValue) error",
				s"CreateOrUpdateUnsafe($entityArgValue)",
			)
			case DaoDeleteRow if !supportsDeactivation => List(
				s"Delete($idArgs) error",
				s"DeleteUnsafe($idArgs)",
			)
			case DaoDeleteRow if supportsDeactivation => List(
				s"Deactivate($idArgs) error",
				s"DeactivateUnsafe($idArgs)",
			)
			case DaoQueryRow(index: Index) => 
				val indexName = goPublicName(index.name.init)
				val args = indexArgList(index)
				List(
					s"ReadBy${indexName}($args) ($entityArgSliceValue, err error)",
					s"ReadBy${indexName}Unsafe($args) ($entityArgSliceValue)",
				)
		}
	}

    blockNamed(
        "type DAO interface",
        dao.operations.flatMap(interfaceOperationTemplate)
    )
}

def daoTemplate(dao: Dao): List[String] = {
	interfaceTemplate(dao) ::: 
	lines(
s"""
// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
	Namespace string                  `json:"namespace"`
	Name      string                  `json:"name"`
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, clientID string) DAO {
	if clientID == "" { panic("Cannot create DAO without clientID") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: clientID + "_${dynamoName(dao.table.entity.name)}",
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic("Cannot create DAO without tableName") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}
""") ::: {
		val templates = OperationImplementationTemplates(dao.table)
		dao.operations.flatMap(templates.apply) ::: 
			templates.utilitiesTemplate ::: 
			templates.allParamsTemplate :::
			templates.updateExpressionTemplate
	}
}

case class OperationImplementationTemplates(table: Table){
	val entity = table.entity
	val idArgs = idArgList(table)
	val idVarNames = entity.primaryKeyFields.map(f => goPrivateName(f.name)).mkString(", ")
	val formatIds = entity.primaryKeyFields.map(f => goPrivateName(f.name) + "==%s").mkString(", ")
	val idFieldNames = entity.primaryKeyFields.map(f => goPublicName(f.name))
	val idDbNames = entity.primaryKeyFields.map(f => dynamoName(f.dbName))
	val entityArgValue = entityArg(table.entity)
	val structVarName = goPrivateName(table.entity.name)
	val structName = goPublicName(table.entity.name)
	val supportsDeactivation = table.entity.supports(DeactivationTrait)

	def apply(operation: DaoOperation): List[String] = {
		operation match {
			case DaoCreateRow => List(
				DaoOperationCreateTemplate,
				DaoOperationCreateUnsafeTemplate
			)
			case DaoReadRow => List(
				DaoOperationReadTemplate,
				DaoOperationReadUnsafeTemplate
			)
			case DaoReadOrEmptyRow => List(
				DaoOperationReadOrEmptyTemplate,
				DaoOperationReadOrEmptyUnsafeTemplate
			)
			case DaoUpdateRow => List(
				DaoOperationCreateOrUpdateTemplate,
				DaoOperationCreateOrUpdateUnsafeTemplate
			)
			case DaoDeleteRow if !supportsDeactivation => List(
				DaoOperationDeleteTemplate,
				DaoOperationDeleteUnsafeTemplate
			)
			case DaoDeleteRow if supportsDeactivation => List(
				DaoOperationDeactivationTemplate,
				DaoOperationDeactivationUnsafeTemplate
			)
			case DaoQueryRow(index: Index) => 
				val templates = QueryTemplates(index)
				templates.apply
		}
	}

	def DaoOperationCreateTemplate: String =
s"""
// Create saves the $structName.
func (d DAOImpl) Create($structVarName $structName) error {
	emptyFields, ok := ${structVarName}.CollectEmptyFields()
	if ok {
		${
			if(entity.supports(CreatedModifiedTimesTrait)) {
				s"""$structVarName.ModifiedAt = core.TimestampLayout.Format(time.Now())
			|	$structVarName.CreatedAt = $structVarName.ModifiedAt
			|	""".stripMargin
			} else ""
		}err = d.Dynamo.PutTableEntry($structVarName, d.Name)
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}
"""

	def DaoOperationCreateUnsafeTemplate: String =
s"""
// CreateUnsafe saves the $structName.
func (d DAOImpl) CreateUnsafe($entityArgValue) {
	err := d.Create($structVarName)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create $formatIds in %s\\n", ${idFieldNames.map{f => structVarName + "." + f}.mkString(", ")}, d.Name))
}
"""
	def DaoOperationReadTemplate: String =
s"""
// Read reads $structName
func (d DAOImpl) Read($idArgs) (out $structName, err error) {
	var outs []$structName
	outs, err = d.ReadOrEmpty($idVarNames)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found $formatIds in %s\\n", $idVarNames, d.Name)
	}
	return
}
"""
	def DaoOperationReadUnsafeTemplate: String =
s"""
// ReadUnsafe reads the $structName. Panics in case of any errors
func (d DAOImpl) ReadUnsafe($idArgs) $structName {
	out, err := d.Read($idVarNames)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error reading $formatIds in %s\\n", $idVarNames, d.Name))
	return out
}
"""

	def DaoOperationReadOrEmptyTemplate: String =
s"""
// ReadOrEmpty reads $structName
func (d DAOImpl) ReadOrEmpty($idArgs) (out []$structName, err error) {
	var outOrEmpty $structName
	ids := idParams($idVarNames)
	err = d.Dynamo.QueryTable(d.Name, ids, &outOrEmpty)
	if ${entity.primaryKeyFields.map(fld => "outOrEmpty." + goPublicName(fld.name) + " == " + goPrivateName(fld.name)).mkString(" && ")} {
		out = append(out, outOrEmpty)
	} else if err != nil && strings.HasPrefix(err.Error(), "In table ") {
		err = nil // expected not-found error	
	}
	err = errors.Wrapf(err, "$structName DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}
"""
	def DaoOperationReadOrEmptyUnsafeTemplate: String =
s"""
// ReadOrEmptyUnsafe reads the $structName. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe($idArgs) []$structName {
	out, err := d.ReadOrEmpty($idVarNames)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error while reading $formatIds in %s\\n", $idVarNames, d.Name))
	return out
}
"""

	def DaoOperationCreateOrUpdateTemplate: String =
s"""
// CreateOrUpdate saves the $structName regardless of if it exists.
func (d DAOImpl) CreateOrUpdate($structVarName $structName) (err error) {
	${
		if(entity.supports(CreatedModifiedTimesTrait)) {
			s"""$structVarName.ModifiedAt = core.TimestampLayout.Format(time.Now())
           |	if $structVarName.CreatedAt == "" { $structVarName.CreatedAt = $structVarName.ModifiedAt }
           |	""".stripMargin
		} else ""
	}
	var olds []$structName
	olds, err = d.ReadOrEmpty(${idFieldNames.map(structVarName + "." + _).mkString(", ")})
	err = errors.Wrapf(err, "$structName DAO.CreateOrUpdate(id = %v) couldn't ReadOrEmpty", key)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create($structVarName)
			err = errors.Wrapf(err, "$structName DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			emptyFields, ok := ${structVarName}.CollectEmptyFields()
			if ok {
				old := olds[0]
				${
					if(entity.supports(CreatedModifiedTimesTrait)) {
						s"$structVarName.ModifiedAt = core.TimestampLayout.Format(time.Now())" + "\n"
					} else ""
				}
				key := idParams(${idFieldNames.map("old." + _).mkString(", ")})
				expr, exprAttributes, names := updateExpression($structVarName, old)
				input := dynamodb.UpdateItemInput{
					ExpressionAttributeValues: exprAttributes,
					TableName:                 aws.String(d.Name),
					Key:                       key,
					ReturnValues:              aws.String("UPDATED_NEW"),
					UpdateExpression:          aws.String(expr),
				}
				if names != nil { input.ExpressionAttributeNames = *names } // workaround for a pointer to an empty slice
				if err == nil {
					err = d.Dynamo.UpdateItemInternal(input)
				}
				err = errors.Wrapf(err, "$structName DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, d.Name, expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}
"""
	def DaoOperationCreateOrUpdateUnsafeTemplate: String =
s"""
// CreateOrUpdateUnsafe saves the $structName regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe($structVarName $structName) {
	err := d.CreateOrUpdate($structVarName)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("could not create or update %v in %s\\n", $structVarName, d.Name))
}
"""

	def DaoOperationDeleteTemplate: String =
s"""
// Delete removes $structName from db
func (d DAOImpl)Delete($idArgs) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams($idVarNames))
}
"""
	def DaoOperationDeleteUnsafeTemplate: String =
s"""
// DeleteUnsafe deletes $structName and panics in case of errors.
func (d DAOImpl)DeleteUnsafe($idArgs) {
	err := d.Delete($idVarNames)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not delete $formatIds in %s\\n", $idVarNames, d.Name))
}
"""

	def DaoOperationDeactivationTemplate: String =
		s"""
		|// Deactivate "removes" $structName. 
		|// The mechanism is adding timestamp to `DeactivatedOn` field. 
		|// Then, if this field is not empty, the instance is considered to be "active"
		|func (d DAOImpl)Deactivate($idArgs) error {
		|	instance, err := d.Read($idVarNames)
		|	if err == nil {
		|		instance.DeactivatedOn = core.ISODateLayout.Format(time.Now())
		|		err = d.CreateOrUpdate(instance)
		|	}
		|	return err
		|}
		|""".stripMargin
	def DaoOperationDeactivationUnsafeTemplate: String =
		s"""
		|// DeactivateUnsafe "deletes" $structName and panics in case of errors.
		|func (d DAOImpl)DeactivateUnsafe($idArgs) {
		|	err := d.Deactivate($idVarNames)
		|	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not deactivate $formatIds in %s\\n", $idVarNames, d.Name))
		|}
		|""".stripMargin


	case class QueryTemplates(index: Index) {
		val indexShortName = goPublicName(index.name.init)
		val indexFullName = indexName(index)
		val args = indexArgList(index)
		def apply: List[String] = 
			List(
				ReadByIndexTemplate,
				ReadByIndexUnsafeTemplate
			)
		def ReadByIndexTemplate: String = {
s"""
func (d DAOImpl)ReadBy${indexShortName}($args) (out []$structName, err error) {
	var instances []$structName
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "$indexFullName",
		Condition: "${index.fields.zipWithIndex.map{case (fld, i) => dbExprParam(fld) + " = :a" + i.toString }.mkString(" and ")}",
		Attributes: map[string]interface{}{
${index.fields.zipWithIndex.map{case (fld, i) =>  "\t\t\t\":a" + i.toString + "\": " + varName(fld) + ","}.mkString("\n")}
		},
	}, map[string]string{${
		index.fields
			.filter(f => isDynamoReserved(f.dbName))
			.map{ case fld => 
				"\"" + dbExprParam(fld) + "\": \"" + dbName(fld) + "\""
			}.mkString(", ")}}, true, -1, &instances)
	out = ${if(supportsDeactivation) structName+"FilterActive(instances)" else "instances" }
	return
}
"""
		}
		def ReadByIndexUnsafeTemplate: String = {
s"""
func (d DAOImpl)ReadBy${indexShortName}Unsafe($args) (out []$structName) {
	out, err := d.ReadBy${indexShortName}(${index.fields.map(varName).mkString(", ")})
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query $indexFullName on %s table\\n", d.Name))
	return
}
"""
		}
				}

	def utilitiesTemplate: List[String] = {
		blockNamed(s"func idParams($idArgs) map[string]*dynamodb.AttributeValue", 
			blockNamed("params := map[string]*dynamodb.AttributeValue",
				entity.primaryKeyFields.map(fld => 
					"\"" + dynamoName(fld.dbName) + "\": common.Dyn" + dynamoType(fld.tpe) + "(" +
					(fld.tpe match {
						case t: TypeAliasTypeInfo => s"string(${goPrivateName(fld.name)})" 
						case SimpleTypeInfo(_, _: QualifiedName, _, _, _, _) => s"string(${goPrivateName(fld.name)})" 
						case _ => goPrivateName(fld.name)
					}) + "),"
				)
			) :::
			lines("return params")
		)
	}
	def isFieldChangedTemplate(fld: Field, newName: String, oldName: String): String = {
		val accessField = newName + "." + fieldName(fld)
		val accessOldField = oldName + "." + fieldName(fld)
		fld.tpe match {
			case s:SimpleTypeInfo if s.dynamoType == "SS" =>
			  s"!common.StringArraysEqual($accessField, $accessOldField)"
			case _ => s"""$accessField != $accessOldField"""
		}
	}
	def dynFieldValueExpr(structVarName: String, fld: Field): String = {
		val accessField = structVarName + "." + fieldName(fld)
		// val accessOldField = "old." + fieldName(fld)
		val fldValueAsString = fld.tpe match {
			case SimpleTypeInfo(_, QualifiedName(_,_), _, _, _, _) => s"string($accessField)"
			case TypeAliasTypeInfo(_) => s"string($accessField)"
			case _ => accessField
		}
		"common.Dyn" + dynamoType(fld.tpe) + "(" + fldValueAsString + ")"
	}
	def allParamsTemplate: List[String] = {
		
		val body = (
			if(entity.fields.exists(_.tpe.isInstanceOf[StructTypeInfo])) 
				List("panic(\"struct fields are not supported in " + structName + ".CreateOrUpdate/allParams\")", "return")
			else
				lines(s"""params = map[string]*dynamodb.AttributeValue{}""") ++
				(entity.fields ::: entity.virtualFields).filterNot(_.tpe.isInstanceOf[StructTypeInfo]).zipWithIndex.map{
					case (fld, i) => 
						val changed = isFieldChangedTemplate(fld, structVarName, "old")
						val dynAttrValue = dynFieldValueExpr(structVarName, fld)
						s"""if $changed { params[":a$i"] = """ + dynAttrValue + " }"
				} ++
				lines("return")
		)
		blockNamed(s"func allParams($structVarName $structName, old $structName) (params map[string]*dynamodb.AttributeValue)",
			body
		)
	}

	def updateExpressionTemplate: List[String] = {
		val body = 			
		(entity.fields ::: entity.virtualFields).zipWithIndex.map{
			case (fld@Field(_, _ : StructTypeInfo, _, _), _) => 
				"panic(\"struct fields are not supported in " + structName + s".CreateOrUpdate/updateExpression " + fieldName(fld) + "\")"
			case (fld, i) => 
				val changed = isFieldChangedTemplate(fld, structVarName, "old")
				val dbExpr = dbExprParam(fld)
				val dynAttrValue = dynFieldValueExpr(structVarName, fld)
				val paramsStmt = s"""params[":a$i"] = """ + dynAttrValue
				val updateStmt = s"updateParts = append(updateParts, "+"\"" + dbExpr + " = :a" + i.toString + "\")"
				val nameStmt = (
					if(dbExpr.startsWith("#")) 
						"fldName := \"" + dbName(fld) + "\"; names[\""+dbExpr+"\"] = &fldName" 
					else 
						""
					)
				s"if $changed { "+updateStmt+"; "+ paramsStmt + "; " + nameStmt + " }"
		}		
		blockNamed(
			s"func updateExpression($structVarName $structName, old $structName) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string)",
			"var updateParts []string" :: 
			"params = map[string]*dynamodb.AttributeValue{}" :: 
			"names := map[string]*string{}" :: 
			body ::: 
			List("""expr = "set " + strings.Join(updateParts, ", ")""",
				"if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty",
				"return"
			)
		)
	}

}
