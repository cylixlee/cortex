import { create } from "zustand"
import type { Conversation, Message, ConversationWithMessages } from "@/types"
import * as conversationApi from "@/api/conversation"

interface ChatState {
  conversations: Conversation[]
  currentConversation: ConversationWithMessages | null
  isLoading: boolean
  isSending: boolean
  loadConversations: () => Promise<void>
  loadConversation: (id: string) => Promise<void>
  createConversation: (title: string) => Promise<Conversation>
  deleteConversation: (id: string) => Promise<void>
  sendMessage: (
    message: string,
    enableRag?: boolean
  ) => Promise<string | undefined>
  clearCurrentConversation: () => void
}

export const useChatStore = create<ChatState>((set, get) => ({
  conversations: [],
  currentConversation: null,
  isLoading: false,
  isSending: false,

  loadConversations: async () => {
    set({ isLoading: true })
    try {
      const conversations = await conversationApi.listConversations()
      set({ conversations })
    } catch (error) {
      console.error("Failed to load conversations:", error)
    } finally {
      set({ isLoading: false })
    }
  },

  loadConversation: async (id: string) => {
    set({ isLoading: true })
    try {
      const conversation = await conversationApi.getConversation(id)
      set({ currentConversation: conversation })
    } catch (error) {
      console.error("Failed to load conversation:", error)
    } finally {
      set({ isLoading: false })
    }
  },

  createConversation: async (title: string) => {
    const conversation = await conversationApi.createConversation(title)
    set((state) => ({
      conversations: [conversation, ...state.conversations],
    }))
    return conversation
  },

  deleteConversation: async (id: string) => {
    await conversationApi.deleteConversation(id)
    set((state) => ({
      conversations: state.conversations.filter((c) => c.id !== id),
      currentConversation:
        state.currentConversation?.id === id ? null : state.currentConversation,
    }))
  },

  sendMessage: async (message: string, enableRag?: boolean) => {
    const { currentConversation } = get()
    set({ isSending: true })

    const userMessage: Message = {
      id: Date.now().toString(),
      role: "user",
      content: message,
      created_at: new Date().toISOString(),
    }

    const assistantMessage: Message = {
      id: (Date.now() + 1).toString(),
      role: "assistant",
      content: "",
      created_at: new Date().toISOString(),
    }

    if (currentConversation) {
      set({
        currentConversation: {
          ...currentConversation,
          messages: [
            ...currentConversation.messages,
            userMessage,
            assistantMessage,
          ],
        },
      })
    } else {
      const newConv: ConversationWithMessages = {
        id: "",
        title: message.slice(0, 50),
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        messages: [userMessage, assistantMessage],
      }
      set({ currentConversation: newConv })
    }

    try {
      const conversationId = await conversationApi.sendMessage(
        message,
        currentConversation?.id,
        (chunk) => {
          set((state) => {
            if (!state.currentConversation) return state
            const messages = [...state.currentConversation.messages]
            const lastMsg = messages[messages.length - 1]
            if (lastMsg?.role === "assistant") {
              lastMsg.content += chunk
            }
            return {
              currentConversation: { ...state.currentConversation, messages },
            }
          })
        },
        enableRag
      )
      return conversationId
    } catch (error) {
      console.error("Failed to send message:", error)
      throw error
    } finally {
      set({ isSending: false })
    }
  },

  clearCurrentConversation: () => {
    set({ currentConversation: null })
  },
}))
