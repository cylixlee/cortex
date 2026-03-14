import { getToken, API_BASE } from "./client"
import type { Conversation, ConversationWithMessages } from "@/types"

export async function listConversations(): Promise<Conversation[]> {
  const token = getToken()
  const response = await fetch(`${API_BASE}/conversations`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!response.ok) throw new Error("Failed to list conversations")
  return response.json()
}

export async function createConversation(title: string): Promise<Conversation> {
  const token = getToken()
  const response = await fetch(`${API_BASE}/conversations`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ title }),
  })
  if (!response.ok) throw new Error("Failed to create conversation")
  return response.json()
}

export async function getConversation(
  id: string
): Promise<ConversationWithMessages> {
  const token = getToken()
  const response = await fetch(`${API_BASE}/conversations/${id}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!response.ok) throw new Error("Failed to get conversation")
  return response.json()
}

export async function deleteConversation(id: string): Promise<void> {
  const token = getToken()
  const response = await fetch(`${API_BASE}/conversations/${id}`, {
    method: "DELETE",
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!response.ok) throw new Error("Failed to delete conversation")
}

export async function sendMessage(
  message: string,
  conversationId?: string,
  onChunk?: (content: string) => void,
  enableRag?: boolean,
  onTitle?: (conversationId: string, title: string) => void
): Promise<string | undefined> {
  const token = getToken()
  if (!token) throw new Error("Not authenticated")

  const body: {
    message: string
    conversation_id?: string
    enable_rag?: boolean
  } = { message }
  if (conversationId) {
    body.conversation_id = conversationId
  }
  if (enableRag) {
    body.enable_rag = enableRag
  }

  const response = await fetch(`${API_BASE}/chat`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(body),
  })

  if (!response.ok) {
    const err = await response.json()
    throw new Error(err.error || "Chat failed")
  }

  const reader = response.body?.getReader()
  if (!reader) throw new Error("No response body")

  const decoder = new TextDecoder()
  let newConversationId: string | undefined

  while (true) {
    const { done, value } = await reader.read()
    if (done) break

    const text = decoder.decode(value)
    const lines = text.split("\n")

    for (const line of lines) {
      if (line.startsWith("data: ")) {
        const data = line.slice(6).trim()

        if (data === "[DONE]") {
          return newConversationId
        }

        if (data.startsWith("{")) {
          try {
            const parsed = JSON.parse(data)
            if (parsed.conversation_id) {
              newConversationId = parsed.conversation_id
            }
            if (parsed.title && onTitle && parsed.conversation_id) {
              onTitle(parsed.conversation_id, parsed.title)
            } else if (parsed.error) {
              throw new Error(parsed.error)
            }
          } catch {
            if (data) {
              onChunk?.(data)
            }
          }
          continue
        }

        if (data) {
          onChunk?.(data)
        }
      }
    }
  }

  return newConversationId
}
