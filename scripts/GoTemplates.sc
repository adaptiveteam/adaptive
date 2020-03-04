import $file.Meta
import Meta._

import $file.Templates
import Templates._

import $file.VirtualFields
import VirtualFields._

def jsonName(name: SimpleName): String = name.parts.map(_.map(_.toLower)).mkString("_")

def jsonOptional(tpe: TypeInfo): String = tpe match {
    case SimpleTypeInfo(_, _, _, _, true, _) => ",omitempty"
    case _ => ""
}

def jsonHint(fld: Field): String = " `json:\"" + jsonName(fld.dbName) + jsonOptional(fld.tpe) + "\"`"

def goType(tpe: TypeInfo): String = tpe match {
    case SimpleTypeInfo(_, golangTypeName: SimpleName, _, _, _, _) => goPrivateName(golangTypeName)
    case SimpleTypeInfo(_, golangTypeName: QualifiedName, _, _, _, _) => goPublicName(golangTypeName)
    case StructTypeInfo(name) => goPublicName(name)
    case TypeAliasTypeInfo(TypeAlias(name, _)) => goPublicName(name)
}


def defaultTypeValueLiteralTemplate(tpe: TypeInfo): String = tpe match {
    case SimpleTypeInfo(_,_, defaultTypeValueLiteral, _, _, _) => defaultTypeValueLiteral
    case StructTypeInfo(name) => goPublicName(name) + "{}"
}
  
def fieldArg(field: Field): String = arg(field)

def arg(field: Field): String = 
    varName(field) + " " + goType(field.tpe)

def argSlice(field: Field): String = 
    varName(field) + " []" + goType(field.tpe)

def varName(field: Field): String = 
    goPrivateName(field.name)

def fieldName(field: Field): String = 
    goPublicName(field.name)

def commentsTemplate(comments: List[String]): List[String] =
    comments.map("// " + _)

def structField(fld: Field): List[String] = fld match {
    case Field(name, tpe, comments, _) =>
        val n = goPublicName(fld.name)
        commentsTemplate(comments) :+ 
        (
            (if(n.isEmpty)"" else n + " ") + goType(tpe) + 
                jsonHint(fld)
        )
}

def structTemplate(entity: Entity): List[String] = entity match {
    case Entity(name, _, _, comments, _) =>
        commentsTemplate(comments) :::
        blockNamed(
            "type " + goPublicName(name) + " struct ", 
            (entity.fields ::: entity.virtualFields).flatMap(structField)
        ) ::: 
        deactivationFilter(entity) :::
        isValidTemplate(entity) :::
        toJsonTemplate(entity)
}

def importClauseTemplate(i: ImportClause): String = i match {
    case ImportClause(Some(p), url) => p + " \"" + url + "\"" 
    case ImportClause(None, url) => "\"" + url + "\"" 
}

def importsTemplate(i: Imports): List[String] = 
  i.importClauses.headOption.toList.flatMap(_ => 
    parensBlockNamed(
        "import",
        i.importClauses.map(importClauseTemplate)
    )
  )

def deactivationFilter(entity: Entity): List[String] =
    entity.traits.filter(_ == DeactivationTrait).flatMap(_ => 
        lines(
        s"""
            |// ${goPublicName(entity.name)}FilterActive removes deactivated values
            |func ${goPublicName(entity.name)}FilterActive(in []${goPublicName(entity.name)}) (res []${goPublicName(entity.name)}) {
            |	for _, i := range in {
            |		if i.${goPublicName(deactivatedAtField.name)} == "" {
            |			res = append(res, i)
            |		}
            |	}
            |	return
            |}
            |""".stripMargin
        )
    )

def makefileGoTestTemplate(packages: List[Package]): String = {
    val names = packages.map(_.name).map(goPrivateName)
    """
    |SHELL := /bin/bash
    |
    |""".stripMargin +
    names.map(n => 
            s"""
                |$n/test.log: $n/*.go
                |\tpushd ${n}; go test -v; date > test.log; popd
                |""".stripMargin)
        .mkString("\n") +
    s"""
    |test-all: ${names.map(_ + "/test.log").mkString(" ")}
    |""".stripMargin 
    // +
    // names
    //     .map(n => s"\tpushd ${n}; go test -v; popd")
    //     .mkString(";\\\n")
}

/** This adds CollectEmptyFields method to the entity.
  * It's used for validation in DAO.Create.
  */
def isValidTemplate(entity: Entity): List[String] =
    "" ::
    "// CollectEmptyFields returns entity field names that are empty." ::
    "// It also returns the boolean ok-flag if the list is empty." ::
    blockNamed(
        s"func (${goPrivateName(entity.name)} ${goPublicName(entity.name)})CollectEmptyFields() (emptyFields []string, ok bool)",
        entity.fields.collect{ 
                case f@Field(_, SimpleTypeInfo(_, _, _, _,false, false),_,_) => f
            }
            .map{ f => s"if ${goPrivateName(entity.name)}.${goPublicName(f.name)} == " + defaultTypeValueLiteralTemplate(f.tpe) + 
                " { emptyFields = append(emptyFields, \""+goPublicName(f.name)+"\")}"
            } :::
        "ok = len(emptyFields) == 0" :: 
        "return" :: Nil
    )

def toJsonTemplate(entity: Entity): List[String] =
    "// ToJSON returns json string" ::
    s"func (${goPrivateName(entity.name)} ${goPublicName(entity.name)}) ToJSON() (string, error) {" ::
    s"	b, err := json.Marshal(${goPrivateName(entity.name)})" ::
    "	return string(b), err" ::
    "}" :: Nil
