import $file.Meta
import Meta._

import $file.Schema
import Schema._

import $file.ProjectTemplates
import ProjectTemplates._

import $file.FileEffects
import FileEffects._

@main def main(args: String*) = {
    val rootDir = ".."
    val fcs = renderProjects(workspace).map(rebase(rootDir))
    if(args.contains("--dry-run"))
        fcs.foreach(dryRun)
    else
        fcs.foreach(saveFile)
}
  