import { getToken } from './auth'

const API_BASE = 'http://localhost:8080/api/v1'

export interface Conversation {
  id: string
  title: string
  created_at: string
  updated_at: string
}

export interface Message {
  id: string
  role: 'user' | 'assistant'
  content: string
  created_at: string
}

export interface ConversationWithMessages extends Conversation {
  messages: Message[]
}

async function fetchWithAuth(url: string, options: RequestInit = {}): Promise<Response> {
  const token = getToken()
  if (!token) throw new Error('Not authenticated')

  const response = await fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      Authorization: `Bearer ${token}`,
    },
  })
  return response
}

export async function listConversations(): Promise<Conversation[]> {
  const response = await fetchWithAuth(`${API_BASE}/conversations`)
  if (!response.ok) throw new Error('Failed to list conversations')
  return response.json()
}

export async function createConversation(title: string): Promise<Conversation> {
  const response = await fetchWithAuth(`${API_BASE}/conversations`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title }),
  })
  if (!response.ok) throw new Error('Failed to create conversation')
  return response.json()
}

export async function getConversation(id: string): Promise<ConversationWithMessages> {
  const response = await fetchWithAuth(`${API_BASE}/conversations/${id}`)
  if (!response.ok) throw new Error('Failed to get conversation')
  return response.json()
}

export async function deleteConversation(id: string): Promise<void> {
  const response = await fetchWithAuth(`${API_BASE}/conversations/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) throw new Error('Failed to delete conversation')
}

export async function sendMessage(
  message: string,
  conversationId?: string,
  onChunk?: (content: string) => void
): Promise<string | undefined> {
  const token = getToken()
  if (!token) throw new Error('Not authenticated')

  const body: { message: string; conversation_id?: string } = { message }
  if (conversationId) {
    body.conversation_id = conversationId
  }

  const response = await fetch(`${API_BASE}/chat`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(body),
  })

  if (!response.ok) {
    const err = await response.json()
    throw new Error(err.error || 'Chat failed')
  }

  const reader = response.body?.getReader()
  if (!reader) throw new Error('No response body')

  const decoder = new TextDecoder()
  let newConversationId: string | undefined

  while (true) {
    const { done, value } = await reader.read()
    if (done) break

    const text = decoder.decode(value)
    const lines = text.split('\n')

    for (const line of lines) {
      if (line.startsWith('data: ')) {
        const data = line.slice(6)
        
        if (data.startsWith('{') && data.includes('conversation_id')) {
          try {
            const parsed = JSON.parse(data)
            newConversationId = parsed.conversation_id
          } catch (e) {}
          continue
        }
        
        if (data === '[DONE]') {
          return newConversationId
        }
        if (data.startsWith('[ERROR]')) {
          throw new Error(data)
        }
        onChunk?.(data)
      }
    }
  }
  
  return newConversationId
}
