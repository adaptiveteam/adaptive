import $file.Meta
import Meta._

def readFile(path: java.io.File): String =
    scala.io.Source.fromFile(path).getLines.mkString("\n")+"\n"

def saveFile(fc: FileWithContent): Unit = {
    import java.nio.file.{Paths, Files}
    import java.nio.charset.StandardCharsets
        
    val path = Paths.get(fc.sf.path)
    path.toFile.getParentFile.mkdirs
    if(path.toFile.exists && readFile(path.toFile) == fc.content) {
        println("Skipping unchanged " + fc.sf.path)
    } else {
        Files.write(
            path, 
            fc.content.getBytes(java.nio.charset.StandardCharsets.UTF_8)
        )
        println(fc.sf.path)
    }
}

def dryRun(fc: FileWithContent): Unit = fc match { case FileWithContent(SourceFile(path), content) => 
    println(path)
    println("".padTo(path.length, '='))
    println(content)
    println("".padTo(path.length, '='))
    println
}

