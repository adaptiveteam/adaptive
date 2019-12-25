import $file.Meta
import Meta._

def capitalizeFirst(part: String): String =
    part.charAt(0).toUpper + part.substring(1, part.length)

def capitalizeAll(part: String): String =
    new String(part.toCharArray.map(_.toUpper))

def capitalizePart(part: String): String = part match {
    case "id" => "ID"
    case _ =>
        part.charAt(0).toUpper + part.substring(1, part.length)
}

def snakeCaseName(name: SimpleName): String = name.parts.map(_.map(_.toLower)).mkString("_")

def indentLine(line: String): String = "\t" + line

def indent(lines: List[String]): List[String] = lines.map(indentLine)

def block(lines: List[String]): List[String] = "{" :: lines.map(indentLine) ::: List("}")

def blockNamed(name: String, lines: List[String]): List[String] = name + " {" :: lines.map(indentLine) ::: List("}")

def parensBlockNamed(name: String, lines: List[String]): List[String] = name + " (" :: lines.map(indentLine) ::: List(")")

def bracketBlockNamed(name: String, lines: List[String]): List[String] = name + " [" :: lines.map(indentLine) ::: List("]")
def prependToFirstLine(line0: String, lines: List[String]): List[String] = 
    lines match {
        case head::tail =>     (line0 + head) :: tail
        case Nil => List(line0)
    }

def lines(str: String): List[String] = str.split("\n").toList

def concatLines(lines: List[String]): String = 
    lines.mkString("\n") + "\n"

def goPublicName(name: Name): String = name match {
    case SimpleName(parts) => parts.map(capitalizePart).mkString("")
    case QualifiedName(p, n) => goPrivateName(p) + "." + goPublicName(n)
}

def goPrivateName(name: Name): String = name match {
    case SimpleName(parts) => parts match {
        case head::tail => 
            (head.map(_.toLower) :: 
            tail.map(capitalizePart)
            ).mkString("")
        case Nil => ""
    }
    case QualifiedName(p, n) => goPrivateName(p) + "." + goPrivateName(n)
}

def quote(str: String): String = "\"" + str + "\""
