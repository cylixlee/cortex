import { useEffect, useRef, useState } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { Send, Loader2 } from "lucide-react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { useChatStore, useAuthStore } from "@/stores"
import { cn } from "@/lib/utils"

export default function ChatPage() {
  const navigate = useNavigate()
  const { id } = useParams()
  const { user } = useAuthStore()
  const {
    currentConversation,
    isLoading,
    isSending,
    loadConversation,
    sendMessage,
  } = useChatStore()

  const [input, setInput] = useState("")
  const [enableRag, setEnableRag] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

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

  return (
    <div className="flex h-full flex-col">
      <ScrollArea className="flex-1">
        <div className="mx-auto max-w-3xl p-4">
          <div className="flex flex-col gap-4">
            {!currentConversation?.messages?.length && !isLoading && (
              <div className="flex min-h-[50vh] items-center justify-center text-muted-foreground">
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
        </div>
      </ScrollArea>

      <div className="p-4">
        <div className="mx-auto max-w-3xl">
          <div className="relative">
            <Textarea
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Type a message..."
              disabled={isSending}
              className="min-h-[80px] resize-none pr-24 pb-14"
            />
            <div className="absolute bottom-2 left-2">
              <Button
                variant={enableRag ? "default" : "outline"}
                size="sm"
                onClick={() => setEnableRag(!enableRag)}
                disabled={isSending}
              >
                Skill RAG
              </Button>
            </div>
            <div className="absolute right-2 bottom-2">
              <Button
                onClick={handleSend}
                disabled={!input.trim() || isSending}
              >
                <Send className="size-4" />
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
