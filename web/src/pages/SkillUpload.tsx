import { useState, useEffect } from "react"
import { useNavigate } from "react-router-dom"
import { Loader2, Upload } from "lucide-react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
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

export default function SkillUploadPage() {
  const navigate = useNavigate()
  const { uploadSkill, subscribeStatus, skillStatus, loadSkills } =
    useSkillStore()

  const [name, setName] = useState("")
  const [file, setFile] = useState<File | null>(null)
  const [uploading, setUploading] = useState(false)
  const [skillId, setSkillId] = useState<string | null>(null)
  const [currentStatus, setCurrentStatus] = useState<{
    stage: number
    name: string
  } | null>(null)

  useEffect(() => {
    if (!skillId) return
    let unsub: (() => void) | undefined
    subscribeStatus(skillId).then((fn) => {
      unsub = fn
    })
    return () => unsub?.()
  }, [skillId, subscribeStatus])

  useEffect(() => {
    setCurrentStatus(skillStatus)
  }, [skillStatus])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!name.trim() || !file) return

    setUploading(true)
    try {
      const uploadId = await uploadSkill(name.trim(), file)
      setSkillId(uploadId)
      toast.success("Upload started")
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Upload failed")
      setUploading(false)
    }
  }

  useEffect(() => {
    if (!uploading || !skillId) return
    const checkInterval = setInterval(async () => {
      const { useSkillStore: store } = await import("@/stores")
      const status = store.getState().skillStatus
      if (status) {
        setCurrentStatus(status)
        if (status.name === "completed" || status.name === "failed") {
          setUploading(false)
          clearInterval(checkInterval)
          if (status.name === "completed") {
            toast.success("Skill processed successfully")
            loadSkills()
            navigate("/skills")
          } else {
            toast.error("Skill processing failed")
          }
        }
      }
    }, 2000)
    return () => clearInterval(checkInterval)
  }, [uploading, skillId, loadSkills, navigate])

  const isProcessing = uploading && skillId

  return (
    <div className="h-full w-full">
      <div className="flex h-full items-center justify-center">
        <Card className="w-full max-w-md">
          <CardContent className="pt-6">
            <form onSubmit={handleSubmit} className="flex flex-col gap-4">
              <div className="flex flex-col gap-2">
                <Label htmlFor="name">Skill Name</Label>
                <Input
                  id="name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="my-awesome-skill"
                  disabled={isProcessing}
                />
              </div>
              <div className="flex flex-col gap-2">
                <Input
                  id="file"
                  type="file"
                  accept=".zip"
                  onChange={(e) => setFile(e.target.files?.[0] || null)}
                  disabled={isProcessing}
                />
              </div>

              {isProcessing && (
                <div className="flex flex-col gap-2 rounded-md bg-muted p-3">
                  <div className="flex items-center justify-between text-sm">
                    <span>Processing...</span>
                    <span className="text-muted-foreground">
                      {currentStatus
                        ? stageNames[currentStatus.stage]
                        : "Initializing"}
                    </span>
                  </div>
                  <Progress
                    value={currentStatus?.name === "completed" ? 100 : 50}
                  />
                </div>
              )}

              <Button
                type="submit"
                disabled={!name.trim() || !file || isProcessing}
              >
                {isProcessing ? (
                  <>
                    <Loader2 className="mr-2 size-4 animate-spin" />
                    Processing...
                  </>
                ) : (
                  <>
                    <Upload className="mr-2 size-4" />
                    Upload
                  </>
                )}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
