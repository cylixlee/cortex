import { useEffect } from "react"
import { useParams, useNavigate } from "react-router-dom"
import { ArrowLeft, Download, Loader2 } from "lucide-react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { ScrollArea } from "@/components/ui/scroll-area"
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

export default function SkillDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const {
    currentSkill,
    isLoading,
    loadSkill,
    downloadSkill,
    clearCurrentSkill,
  } = useSkillStore()

  useEffect(() => {
    if (id) {
      loadSkill(id)
    }
    return () => clearCurrentSkill()
  }, [id, loadSkill, clearCurrentSkill])

  const handleDownload = async () => {
    if (!currentSkill) return
    try {
      await downloadSkill(currentSkill.id, currentSkill.name)
      toast.success("Download started")
    } catch {
      toast.error("Failed to download skill")
    }
  }

  if (isLoading) {
    return (
      <div>
        <Skeleton className="mb-6 h-8 w-32" />
        <Skeleton className="h-64 w-full" />
      </div>
    )
  }

  if (!currentSkill) {
    return (
      <div>
        <p>Skill not found</p>
      </div>
    )
  }

  return (
    <div>
      <div className="mb-6 flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={() => navigate("/skills")}>
          <ArrowLeft className="size-4" />
        </Button>
        <div className="flex-1">
          <h1 className="text-2xl font-bold">{currentSkill.name}</h1>
          <p className="text-muted-foreground">
            {currentSkill.description || "No description"}
          </p>
        </div>
        {currentSkill.status === "completed" && (
          <Button onClick={handleDownload}>
            <Download className="mr-2 size-4" />
            Download
          </Button>
        )}
      </div>

      <div className="mb-6 flex items-center gap-2">
        <Badge
          variant={
            currentSkill.status === "completed"
              ? "secondary"
              : currentSkill.status === "failed"
                ? "destructive"
                : "outline"
          }
        >
          {stageNames[currentSkill.stage] || currentSkill.status}
        </Badge>
        <span className="text-sm text-muted-foreground">
          Created {new Date(currentSkill.created_at).toLocaleDateString()}
        </span>
      </div>

      {currentSkill.status !== "completed" ? (
        <Card>
          <CardContent className="flex items-center justify-center py-12">
            <div className="flex items-center gap-2 text-muted-foreground">
              <Loader2 className="size-4 animate-spin" />
              Processing...
            </div>
          </CardContent>
        </Card>
      ) : (
        <Tabs defaultValue="overview" className="w-full">
          <TabsList>
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="references">
              References ({currentSkill.skill?.references.length || 0})
            </TabsTrigger>
          </TabsList>

          <TabsContent value="overview">
            <Card>
              <CardHeader>
                <CardTitle>Skill Overview</CardTitle>
              </CardHeader>
              <CardContent>
                <ScrollArea className="h-[500px]">
                  <pre className="text-sm whitespace-pre-wrap">
                    {currentSkill.skill?.overview || "No overview available"}
                  </pre>
                </ScrollArea>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="references">
            <div className="flex flex-col gap-4">
              {currentSkill.skill?.references.map((ref, idx) => (
                <Card key={idx}>
                  <CardHeader>
                    <CardTitle className="text-base">{ref.filename}</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <ScrollArea className="h-[400px]">
                      <pre className="font-mono text-sm whitespace-pre-wrap">
                        {ref.content}
                      </pre>
                    </ScrollArea>
                  </CardContent>
                </Card>
              ))}
            </div>
          </TabsContent>
        </Tabs>
      )}
    </div>
  )
}
