import { getToken, API_BASE } from "./client"
import type { Skill, SkillDetail, SkillStatus } from "@/types"

export async function listSkills(): Promise<Skill[]> {
  const token = getToken()
  const response = await fetch(`${API_BASE}/skills`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  const data = await response.json()
  return data.skills || []
}

export async function uploadSkill(
  name: string,
  file: File,
  onProgress?: (progress: number) => void
): Promise<{ skill_id: string; status: string }> {
  return new Promise((resolve, reject) => {
    const token = getToken()
    const xhr = new XMLHttpRequest()
    const formData = new FormData()
    formData.append("name", name)
    formData.append("file", file)

    xhr.upload.addEventListener("progress", (e) => {
      if (e.lengthComputable && onProgress) {
        const progress = Math.round((e.loaded / e.total) * 100)
        onProgress(progress)
      }
    })

    xhr.addEventListener("load", () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        resolve(JSON.parse(xhr.responseText))
      } else {
        reject(new Error(`HTTP error! status: ${xhr.status}`))
      }
    })

    xhr.addEventListener("error", () => {
      reject(new Error("Upload failed"))
    })

    xhr.open("POST", `${API_BASE}/skills/upload`)
    xhr.setRequestHeader("Authorization", `Bearer ${token}`)
    xhr.send(formData)
  })
}

export async function getSkill(id: string): Promise<SkillDetail> {
  const token = getToken()
  const response = await fetch(`${API_BASE}/skills/${id}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  return response.json()
}

export async function deleteSkill(id: string): Promise<void> {
  const token = getToken()
  const response = await fetch(`${API_BASE}/skills/${id}`, {
    method: "DELETE",
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }
}

export async function getSkillStatus(id: string): Promise<SkillStatus> {
  const token = getToken()
  const response = await fetch(`${API_BASE}/skills/${id}/status`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  return response.json()
}

export async function subscribeSkillStatus(
  id: string,
  onStage: (stage: number, name: string) => void
): Promise<() => void> {
  const token = getToken()
  const response = await fetch(`${API_BASE}/skills/${id}/status/sse`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })

  if (!response.ok || !response.body) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  const reader = response.body.getReader()
  const decoder = new TextDecoder()

  const read = async () => {
    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      const text = decoder.decode(value)
      const lines = text.split("\n")

      for (const line of lines) {
        if (line.startsWith("data: ")) {
          const data = JSON.parse(line.slice(6))
          onStage(data.stage, data.name)

          if (data.name === "completed" || data.name === "failed") {
            return
          }
        }
      }
    }
  }

  read()

  return () => {
    reader.cancel()
  }
}

export async function downloadSkill(
  id: string,
  filename: string
): Promise<void> {
  const token = getToken()
  const response = await fetch(`${API_BASE}/skills/${id}/download`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  const blob = await response.blob()
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement("a")
  a.href = url
  a.download = filename + ".zip"
  document.body.appendChild(a)
  a.click()
  a.remove()
  window.URL.revokeObjectURL(url)
}
