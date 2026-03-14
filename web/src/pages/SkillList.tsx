import { useEffect } from "react"
import { useNavigate } from "react-router-dom"
import { Plus, Download, Trash2 } from "lucide-react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import { useSkillStore } from "@/stores"

const stageNames = [
  "Pending",
  "Uploading",
  "Processing",
  "Overview",
  "API Retrieval",
  "Completed",
  "Failed",
]

export default function SkillListPage() {
  const navigate = useNavigate()
  const { skills, isLoading, loadSkills, deleteSkill, downloadSkill } =
    useSkillStore()

  useEffect(() => {
    loadSkills()
  }, [loadSkills])

  const handleDelete = async (id: string) => {
    try {
      await deleteSkill(id)
      toast.success("Skill deleted")
    } catch {
      toast.error("Failed to delete skill")
    }
  }

  const handleDownload = async (id: string, name: string) => {
    try {
      await downloadSkill(id, name)
      toast.success("Download started")
    } catch {
      toast.error("Failed to download skill")
    }
  }

  return (
    <div className="h-full w-full">
      <div className="flex h-full flex-col">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold">My Skills</h1>
          <Button onClick={() => navigate("/skills/upload")}>
            <Plus className="mr-2 size-4" />
            Upload Skill
          </Button>
        </div>

        <div className="flex flex-1">
          {isLoading ? (
            <div className="w-full space-y-4">
              {[1, 2, 3].map((i) => (
                <Skeleton key={i} className="h-32 w-full" />
              ))}
            </div>
          ) : skills.length === 0 ? (
            <div className="flex h-full w-full flex-col items-center justify-center">
              <p className="mb-4 text-muted-foreground">No skills yet</p>
              <Button onClick={() => navigate("/skills/upload")}>
                Upload your first skill
              </Button>
            </div>
          ) : (
            <div className="w-full space-y-4 pt-4">
              {skills.map((skill) => (
                <Card
                  key={skill.id}
                  className="cursor-pointer transition-colors hover:bg-accent/50"
                  onClick={() => navigate(`/skills/${skill.id}`)}
                >
                  <CardHeader className="flex flex-row items-center justify-between pb-2">
                    <div>
                      <CardTitle>{skill.name}</CardTitle>
                      <CardDescription>
                        {skill.description || "No description"}
                      </CardDescription>
                    </div>
                    <div className="flex items-center gap-2">
                      {skill.status === "completed" && (
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={(e) => {
                            e.stopPropagation()
                            handleDownload(skill.id, skill.name)
                          }}
                        >
                          <Download className="mr-1 size-4" />
                          Download
                        </Button>
                      )}
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={(e) => {
                          e.stopPropagation()
                          handleDelete(skill.id)
                        }}
                      >
                        <Trash2 className="size-4" />
                      </Button>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="flex items-center gap-2">
                      <Badge
                        variant={
                          skill.status === "completed"
                            ? "secondary"
                            : skill.status === "failed"
                              ? "destructive"
                              : "outline"
                        }
                      >
                        {stageNames[skill.stage] || skill.status}
                      </Badge>
                      <span className="text-xs text-muted-foreground">
                        Created{" "}
                        {new Date(skill.created_at).toLocaleDateString()}
                      </span>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
