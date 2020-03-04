import $file.Meta
import Meta._

import $file.Templates
import Templates._

import $file.DynamoReservedWords
import DynamoReservedWords._

def dynamoType(tpe: TypeInfo): String = tpe match {
    case SimpleTypeInfo(_, _, _, dbType, _, _) => dbType
    case TypeAliasTypeInfo(TypeAlias(_, tpe)) => dynamoType(tpe)
    case _ => "<unsupported type>"
}

def dynamoName(name: SimpleName): String = 
  snakeCaseName(name)

def dbName(field: Field): String = dynamoName(field.dbName)

def isDynamoReserved(name: SimpleName): Boolean =
    dynamoReservedWordsLower.contains(dynamoName(name))

def dbExprParam(field: Field): String = {
    val name = dbName(field)
    if(dynamoReservedWordsLower.contains(name)) 
        "#" + name 
    else
        name
}

def indexName(index: Index): String = 
    goPublicName(index.name)
