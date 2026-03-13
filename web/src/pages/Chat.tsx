import { useEffect, useRef, useState } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { Send, Trash2, Plus, Loader2 } from "lucide-react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { useChatStore, useAuthStore } from "@/stores"
import { cn } from "@/lib/utils"

export default function ChatPage() {
  const navigate = useNavigate()
  const { id } = useParams()
  const { user } = useAuthStore()
  const {
    conversations,
    currentConversation,
    isLoading,
    isSending,
    loadConversations,
    loadConversation,
    createConversation,
    deleteConversation,
    sendMessage,
  } = useChatStore()

  const [input, setInput] = useState("")
  const [enableRag, setEnableRag] = useState(true)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    loadConversations()
  }, [loadConversations])

  useEffect(() => {
    if (id) {
      loadConversation(id)
    }
  }, [id, loadConversation])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }, [currentConversation?.messages])

  const handleSend = async () => {
    if (!input.trim() || isSending) return

    const message = input.trim()
    setInput("")

    try {
      const conversationId = await sendMessage(message, enableRag)
      if (conversationId && !id) {
        navigate(`/chat/${conversationId}`)
        loadConversations()
      }
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Failed to send message"
      )
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const handleNewChat = async () => {
    const conv = await createConversation("New Chat")
    navigate(`/chat/${conv.id}`)
  }

  const handleDelete = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation()
    await deleteConversation(id)
    if (id === currentConversation?.id) {
      navigate("/chat")
    }
  }

  return (
    <div className="flex h-screen bg-background">
      <aside className="w-64 border-r bg-card">
        <div className="flex h-14 items-center border-b px-4">
          <Button
            variant="ghost"
            size="sm"
            className="w-full justify-start gap-2"
            onClick={handleNewChat}
          >
            <Plus className="size-4" />
            New Chat
          </Button>
        </div>
        <ScrollArea className="h-[calc(100vh-3.5rem)]">
          <div className="flex flex-col gap-1 p-2">
            {conversations.map((conv) => (
              <div
                key={conv.id}
                onClick={() => navigate(`/chat/${conv.id}`)}
                className={cn(
                  "flex cursor-pointer items-center gap-2 rounded-md px-3 py-2 text-sm hover:bg-accent",
                  id === conv.id && "bg-accent"
                )}
              >
                <span className="flex-1 truncate">{conv.title}</span>
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-6 opacity-50 hover:opacity-100"
                  onClick={(e) => handleDelete(conv.id, e)}
                >
                  <Trash2 className="size-3" />
                </Button>
              </div>
            ))}
          </div>
        </ScrollArea>
      </aside>

      <main className="flex flex-1 flex-col">
        <header className="flex h-14 items-center justify-between border-b px-4">
          <h1 className="text-lg font-medium">
            {currentConversation?.title || "New Chat"}
          </h1>
          <label className="flex items-center gap-2 text-sm">
            <input
              type="checkbox"
              checked={enableRag}
              onChange={(e) => setEnableRag(e.target.checked)}
              className="rounded border-input"
            />
            Enable RAG
          </label>
        </header>

        <ScrollArea className="flex-1 p-4">
          <div className="flex flex-col gap-4">
            {!currentConversation?.messages?.length && !isLoading && (
              <div className="flex h-full items-center justify-center text-muted-foreground">
                <p>Start a conversation...</p>
              </div>
            )}
            {currentConversation?.messages?.map((msg) => (
              <div
                key={msg.id}
                className={cn(
                  "flex gap-3",
                  msg.role === "user" ? "justify-end" : "justify-start"
                )}
              >
                {msg.role === "assistant" && (
                  <Avatar className="size-8">
                    <AvatarFallback>AI</AvatarFallback>
                  </Avatar>
                )}
                <div
                  className={cn(
                    "max-w-[70%] rounded-lg px-4 py-2",
                    msg.role === "user"
                      ? "bg-primary text-primary-foreground"
                      : "bg-muted"
                  )}
                >
                  <p className="whitespace-pre-wrap">{msg.content}</p>
                </div>
                {msg.role === "user" && (
                  <Avatar className="size-8">
                    <AvatarFallback>
                      {user?.email?.[0]?.toUpperCase() || "U"}
                    </AvatarFallback>
                  </Avatar>
                )}
              </div>
            ))}
            {isSending && (
              <div className="flex gap-3">
                <Avatar className="size-8">
                  <AvatarFallback>AI</AvatarFallback>
                </Avatar>
                <div className="flex items-center gap-2 rounded-lg bg-muted px-4 py-2">
                  <Loader2 className="size-4 animate-spin" />
                  <span className="text-sm text-muted-foreground">
                    Thinking...
                  </span>
                </div>
              </div>
            )}
            <div ref={messagesEndRef} />
          </div>
        </ScrollArea>

        <div className="border-t p-4">
          <div className="flex gap-2">
            <Input
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Type a message..."
              disabled={isSending}
            />
            <Button onClick={handleSend} disabled={!input.trim() || isSending}>
              <Send className="size-4" />
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
