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

import $file.DaoTemplates
import DaoTemplates._

def daoFunctionsTemplate(dao: ConnectionBasedDao): List[String] = {
// 	lines(
// s"""
// """) ::: 
	{
		val templates = DaoFunctionTemplates(dao.table)
		dao.operations.flatMap(templates.apply)
	}
}

case class DaoFunctionTemplates(table: Table){
	val entity = table.entity
	val idArgs = idArgList(table)
	// val hashKey = table.defaultIndex.hashKey
	// val hashIdArgs = fieldArg(hashKey)
	val idVarNames = entity.primaryKeyFields.map(f => goPrivateName(f.name)).mkString(", ")
	// val hashIdVarName = goPrivateName(hashKey.name)
	val formatIds = entity.primaryKeyFields.map(f => goPrivateName(f.name) + "==%s").mkString(", ")
	val idFieldNames = entity.primaryKeyFields.map(f => goPublicName(f.name))
	// val hashIdFieldName = goPublicName(hashKey.name)
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
			case DaoReadChildren => QueryByHashKeyTemplates(
				table.defaultIndex, 
				isDefaultIndex = true,
			).apply
			case DaoReadOrEmptyRow => List(
				DaoOperationReadOrEmptyTemplate,
				DaoOperationReadOrEmptyUnsafeTemplate,
				DaoOperationReadOrEmptyIncludingInactiveTemplate,
				DaoOperationReadOrEmptyIncludingInactiveUnsafeTemplate
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
			case DaoQueryRowByHashKey(index: Index) => 
				val templates = QueryByHashKeyTemplates(index)
				templates.apply
		}
	}

	def DaoOperationCreateTemplate: String =
		s"""
		|// Create saves the $structName.
		|func Create($structVarName $structName) common.ConnectionProc {
		|	return func (conn common.DynamoDBConnection) (err error) {
		|		emptyFields, ok := ${structVarName}.CollectEmptyFields()
		|		if ok {
		|			${
						if(entity.supports(CreatedModifiedTimesTrait)) {
							s"""$structVarName.ModifiedAt = core.CurrentRFCTimestamp()
							|	$structVarName.CreatedAt = $structVarName.ModifiedAt
							|	""".stripMargin
						} else ""
					}
		|			err = conn.Dynamo.PutTableEntry($structVarName, TableName(conn.ClientID))
		|		} else {
		|			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		|		}
		|		return
		|	}
		|}
		|""".stripMargin

	def DaoOperationCreateUnsafeTemplate: String =
		s"""
		|// CreateUnsafe saves the $structName.
		|func CreateUnsafe($entityArgValue) func (conn common.DynamoDBConnection) {
		|	return func (conn common.DynamoDBConnection) {
		|		err2 := Create($structVarName)(conn)
		|		core.ErrorHandler(err2, "daos/$structName", fmt.Sprintf("Could not create $formatIds in %s\\n", ${idFieldNames.map{f => structVarName + "." + f}.mkString(", ")}, TableName(conn.ClientID)))
		|	}
		|}
		|""".stripMargin

	def DaoOperationReadTemplate: String =
		s"""
		|// Read reads $structName
		|func Read($idArgs) func (conn common.DynamoDBConnection) (out $structName, err error) {
		|	return func (conn common.DynamoDBConnection) (out $structName, err error) {
		|		var outs []$structName
		|		outs, err = ReadOrEmpty($idVarNames)(conn)
		|		if err == nil && len(outs) == 0 {
		|			err = fmt.Errorf("Not found $formatIds in %s\\n", $idVarNames, TableName(conn.ClientID))
		|		}
		|		if len(outs) > 0 {
		|			out = outs[0]
		|		}
		|		return
		|	}
		|}
		|""".stripMargin

	def DaoOperationReadUnsafeTemplate: String =
		s"""
		|// ReadUnsafe reads the $structName. Panics in case of any errors
		|func ReadUnsafe($idArgs) func (conn common.DynamoDBConnection) $structName {
		|	return func (conn common.DynamoDBConnection) $structName {
		|		out, err2 := Read($idVarNames)(conn)
		|		core.ErrorHandler(err2, "daos/$structName", fmt.Sprintf("Error reading $formatIds in %s\\n", $idVarNames, TableName(conn.ClientID)))
		|		return out
		|	}
		|}
		|""".stripMargin

	def DaoOperationReadOrEmptyTemplate: String =
		s"""
		|// ReadOrEmpty reads $structName
		|func ReadOrEmpty($idArgs) func (conn common.DynamoDBConnection) (out []$structName, err error) {
		|	return func (conn common.DynamoDBConnection) (out []$structName, err error) {
		|       out, err = ReadOrEmptyIncludingInactive($idVarNames)(conn)
		|       ${
					if(entity.supports(DeactivationTrait)){ 
						s"out = ${structName}FilterActive(out)"
					} else ""
				}
		|       
		|		return
		|	}
		|}
		|""".stripMargin

	def DaoOperationReadOrEmptyUnsafeTemplate: String =
		s"""
		|// ReadOrEmptyUnsafe reads the $structName. Panics in case of any errors
		|func ReadOrEmptyUnsafe($idArgs) func (conn common.DynamoDBConnection) []$structName {
		|	return func (conn common.DynamoDBConnection) []$structName {
		|		out, err2 := ReadOrEmpty($idVarNames)(conn)
		|		core.ErrorHandler(err2, "daos/$structName", fmt.Sprintf("Error while reading $formatIds in %s\\n", $idVarNames, TableName(conn.ClientID)))
		|		return out
		|	}
		|}
		|""".stripMargin

	def DaoOperationReadOrEmptyIncludingInactiveTemplate: String =
		s"""
		|// ReadOrEmptyIncludingInactive reads $structName
		|func ReadOrEmptyIncludingInactive($idArgs) func (conn common.DynamoDBConnection) (out []$structName, err error) {
		|	return func (conn common.DynamoDBConnection) (out []$structName, err error) {
		|		var outOrEmpty $structName
		|		ids := idParams($idVarNames)
		|		var found bool
		|		found, err = conn.Dynamo.GetItemOrEmptyFromTable(TableName(conn.ClientID), ids, &outOrEmpty)
		|		if found {
		|			if ${entity.primaryKeyFields.map(fld => "outOrEmpty." + goPublicName(fld.name) + " == " + goPrivateName(fld.name)).mkString(" && ")} {
		|				out = append(out, outOrEmpty)
		|			} else {
		|				err = fmt.Errorf("Requested ids: $formatIds are different from the found ones: $formatIds", $idVarNames, ${idFieldNames.map("outOrEmpty." + _).mkString(", ")}) // unexpected error: found ids != ids
		|			}
		|		}
		|		err = errors.Wrapf(err, "$structName DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		|		return
		|	}
		|}
		|""".stripMargin

	def DaoOperationReadOrEmptyIncludingInactiveUnsafeTemplate: String =
		s"""
		|// ReadOrEmptyIncludingInactiveUnsafe reads the $structName. Panics in case of any errors
		|func ReadOrEmptyIncludingInactiveUnsafeIncludingInactive($idArgs) func (conn common.DynamoDBConnection) []$structName {
		|	return func (conn common.DynamoDBConnection) []$structName {
		|		out, err2 := ReadOrEmptyIncludingInactive($idVarNames)(conn)
		|		core.ErrorHandler(err2, "daos/$structName", fmt.Sprintf("Error while reading $formatIds in %s\\n", $idVarNames, TableName(conn.ClientID)))
		|		return out
		|	}
		|}
		|""".stripMargin

	def DaoOperationCreateOrUpdateTemplate: String =
		s"""
		|// CreateOrUpdate saves the $structName regardless of if it exists.
		|func CreateOrUpdate($structVarName $structName) common.ConnectionProc {
		|	return func (conn common.DynamoDBConnection) (err error) {
		|		${
					if(entity.supports(CreatedModifiedTimesTrait)) {
						s"""$structVarName.ModifiedAt = core.CurrentRFCTimestamp()
					|	if $structVarName.CreatedAt == "" { $structVarName.CreatedAt = $structVarName.ModifiedAt }
					|	""".stripMargin
					} else ""
				}
		|		var olds []$structName
		|		olds, err = ReadOrEmpty(${idFieldNames.map(structVarName + "." + _).mkString(", ")})(conn)
		|		err = errors.Wrapf(err, "$structName DAO.CreateOrUpdate(id = $formatIds) couldn't ReadOrEmpty", ${idFieldNames.map(structVarName + "." + _).mkString(", ")})
		|		if err == nil {
		|			if len(olds) == 0 {
		|				err = Create($structVarName)(conn)
		|				err = errors.Wrapf(err, "$structName DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
		|			} else {
		|				emptyFields, ok := ${structVarName}.CollectEmptyFields()
		|				if ok {
		|					old := olds[0]
		|					${Option.when(entity.supports(CreatedModifiedTimesTrait))(s"$structVarName.CreatedAt  = old.CreatedAt"             ).getOrElse("")}
		|					${Option.when(entity.supports(CreatedModifiedTimesTrait))(s"$structVarName.ModifiedAt = core.CurrentRFCTimestamp()").getOrElse("")}
		|					key := idParams(${idFieldNames.map("old." + _).mkString(", ")})
		|					expr, exprAttributes, names := updateExpression($structVarName, old)
		|					input := dynamodb.UpdateItemInput{
		|						ExpressionAttributeValues: exprAttributes,
		|						TableName:                 aws.String(TableName(conn.ClientID)),
		|						Key:                       key,
		|						ReturnValues:              aws.String("UPDATED_NEW"),
		|						UpdateExpression:          aws.String(expr),
		|					}
		|					if names != nil { input.ExpressionAttributeNames = *names } // workaround for a pointer to an empty slice
		|					if  len(exprAttributes) > 0 { // if there some changes
		|						err = conn.Dynamo.UpdateItemInternal(input)
		|					} else {
		|						// WARN: no changes.
		|					}
		|					err = errors.Wrapf(err, "$structName DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
		|				} else {
		|					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
		|				}
		|			}
		|		}
		|		return 
		|	}
		|}
		|""".stripMargin

	def DaoOperationCreateOrUpdateUnsafeTemplate: String =
		s"""
		|// CreateOrUpdateUnsafe saves the $structName regardless of if it exists.
		|func CreateOrUpdateUnsafe($structVarName $structName) func (conn common.DynamoDBConnection) {
		|	return func (conn common.DynamoDBConnection) {
		|		err2 := CreateOrUpdate($structVarName)(conn)
		|		core.ErrorHandler(err2, "daos/$structName", fmt.Sprintf("could not create or update %v in %s\\n", $structVarName, TableName(conn.ClientID)))
		|	}
		|}
		|""".stripMargin

	def DaoOperationDeleteTemplate: String =
		s"""
		|// Delete removes $structName from db
		|func Delete($idArgs) func (conn common.DynamoDBConnection) error {
		|	return func (conn common.DynamoDBConnection) error {
		|		return conn.Dynamo.DeleteEntry(TableName(conn.ClientID), idParams($idVarNames))
		|	}
		|}
		|""".stripMargin

	def DaoOperationDeleteUnsafeTemplate: String =
		s"""
		|// DeleteUnsafe deletes $structName and panics in case of errors.
		|func DeleteUnsafe($idArgs) func (conn common.DynamoDBConnection) {
		|	return func (conn common.DynamoDBConnection) {
		|		err2 := Delete($idVarNames)(conn)
		|		core.ErrorHandler(err2, "daos/$structName", fmt.Sprintf("Could not delete $formatIds in %s\\n", $idVarNames, TableName(conn.ClientID)))
		|	}
		|}
		|""".stripMargin

	def DaoOperationDeactivationTemplate: String =
		s"""
		|// Deactivate "removes" $structName. 
		|// The mechanism is adding timestamp to `DeactivatedOn` field. 
		|// Then, if this field is not empty, the instance is considered to be "active"
		|func Deactivate($idArgs) func (conn common.DynamoDBConnection) error {
		|	return func (conn common.DynamoDBConnection) error {
		|		instance, err2 := Read($idVarNames)(conn)
		|		if err2 == nil {
		|			instance.${goPublicName(deactivatedAtField.name)} = core.CurrentRFCTimestamp()
    	|			err2 = CreateOrUpdate(instance)(conn)
		|		}
		|		return err2
		|	}
		|}
		|""".stripMargin
	def DaoOperationDeactivationUnsafeTemplate: String =
		s"""
		|// DeactivateUnsafe "deletes" $structName and panics in case of errors.
		|func DeactivateUnsafe($idArgs) func (conn common.DynamoDBConnection) {
		|	return func (conn common.DynamoDBConnection) {
		|		err2 := Deactivate($idVarNames)(conn)
		|		core.ErrorHandler(err2, "daos/$structName", fmt.Sprintf("Could not deactivate $formatIds in %s\\n", $idVarNames, TableName(conn.ClientID)))
		|	}
		|}
		|""".stripMargin


	case class QueryTemplates(index: Index, isDefaultIndex: Boolean = false) {
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
			|func ReadBy${indexShortName}($args) func (conn common.DynamoDBConnection) (out []$structName, err error) {
			|	return func (conn common.DynamoDBConnection) (out []$structName, err error) {
			|		var instances []$structName
			|		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			|			IndexName: "$indexFullName",
			|			Condition: "${index.fields.zipWithIndex.map{case (fld, i) => dbExprParam(fld) + " = :a" + i.toString }.mkString(" and ")}",
			|			Attributes: map[string]interface{}{
			|	${index.fields.zipWithIndex.map{case (fld, i) =>  "\t\t\t\":a" + i.toString + "\": " + varName(fld) + ","}.mkString("\n")}
			|			},
			|		}, map[string]string{${
						index.fields
							.filter(f => isDynamoReserved(f.dbName))
							.map{ case fld => 
								"\"" + dbExprParam(fld) + "\": \"" + dbName(fld) + "\""
							}.mkString(", ")}}, true, -1, &instances)
			|		out = ${if(supportsDeactivation) structName+"FilterActive(instances)" else "instances" }
			|		return
			|	}
			|}
			|""".stripMargin
		}

		def ReadByIndexUnsafeTemplate: String = {
			s"""
			|func ReadBy${indexShortName}Unsafe($args) func (conn common.DynamoDBConnection) (out []$structName) {
			|	return func (conn common.DynamoDBConnection) (out []$structName) {
			|		out, err2 := ReadBy${indexShortName}(${index.fields.map(varName).mkString(", ")})(conn)
			|		core.ErrorHandler(err2, "daos/$structName", fmt.Sprintf("Could not query $indexFullName on %s table\\n", TableName(conn.ClientID)))
			|		return
			|	}
			|}
			|""".stripMargin
		}
	}
	case class QueryByHashKeyTemplates(index: Index, isDefaultIndex: Boolean = false) {
		val indexHashKeyName = goPublicName(index.hashKey.name)
		val fields = List(index.hashKey)
		val indexFullName = indexName(index)
		val args = fieldArg(index.hashKey)
		def apply: List[String] = 
			List(
				ReadByHashKeyIndexTemplate,
				ReadByHashKeyIndexUnsafeTemplate
			)
		def ReadByHashKeyIndexTemplate: String = {
			s"""
			|func ReadByHashKey${indexHashKeyName}($args) func (conn common.DynamoDBConnection) (out []$structName, err error) {
			|	return func (conn common.DynamoDBConnection) (out []$structName, err error) {
			|		var instances []$structName
			|		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			|			${if(isDefaultIndex) "" else "IndexName: string(" + indexFullName + "),"}
			|			Condition: "${dbExprParam(index.hashKey) + " = :a" }",
			|			Attributes: map[string]interface{}{
			|				":a" : ${varName(index.hashKey)},
			|			},
			|		}, map[string]string{${
						fields
							.filter(f => isDynamoReserved(f.dbName))
							.map{ case fld => 
								"\"" + dbExprParam(fld) + "\": \"" + dbName(fld) + "\""
							}.mkString(", ")}}, true, -1, &instances)
			|		out = ${if(supportsDeactivation) structName+"FilterActive(instances)" else "instances" }
			|		return
			|	}
			|}
			|""".stripMargin
		}

		def ReadByHashKeyIndexUnsafeTemplate: String = {
			s"""
			|func ReadByHashKey${indexHashKeyName}Unsafe($args) func (conn common.DynamoDBConnection) (out []$structName) {
			|	return func (conn common.DynamoDBConnection) (out []$structName) {
			|		out, err2 := ReadByHashKey${indexHashKeyName}(${varName(index.hashKey)})(conn)
			|		core.ErrorHandler(err2, "daos/$structName", fmt.Sprintf("Could not query $indexFullName on %s table\\n", TableName(conn.ClientID)))
			|		return
			|	}
			|}
			|""".stripMargin
		}
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

}

 