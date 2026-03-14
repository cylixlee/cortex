import { useEffect, useRef, useState } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { Send } from "lucide-react"
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
  const [displayedTitle, setDisplayedTitle] = useState("")
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const prevTitleRef = useRef("")

  useEffect(() => {
    const newTitle = currentConversation?.title || ""
    if (newTitle && newTitle !== prevTitleRef.current) {
      prevTitleRef.current = newTitle
      let index = 0
      const animate = () => {
        if (index < newTitle.length) {
          setDisplayedTitle(newTitle.slice(0, index + 1))
          index++
          setTimeout(animate, 30)
        }
      }
      animate()
    } else if (!newTitle) {
      prevTitleRef.current = ""
      setDisplayedTitle("")
    }
  }, [currentConversation?.title])

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
    <div className="relative flex flex-col h-full min-h-0">
      {displayedTitle && (
        <div className="border-b px-4 py-3">
          <h2 className="text-lg font-semibold">{displayedTitle}</h2>
        </div>
      )}

      <ScrollArea className="flex-1 min-h-0">
        <div className="min-h-full p-4">
          <div className="flex flex-col gap-4">
            {!currentConversation?.messages?.length && !isLoading && (
              <div className="absolute inset-0 flex items-center justify-center">
                <p className="text-muted-foreground">Start a conversation...</p>
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
                {msg.role === "user" ? (
                  <>
                    <div
                      className={cn(
                        "max-w-[70%] rounded-lg px-4 py-2",
                        "bg-primary text-primary-foreground"
                      )}
                    >
                      <p className="whitespace-pre-wrap">{msg.content}</p>
                    </div>
                    <Avatar className="size-8">
                      <AvatarFallback>
                        {user?.email?.[0]?.toUpperCase() || "U"}
                      </AvatarFallback>
                    </Avatar>
                  </>
                ) : (
                  <div className="w-full">
                    <p className="whitespace-pre-wrap break-all">{msg.content}</p>
                  </div>
                )}
              </div>
            ))}
            {isSending && (
              <div className="w-full">
                <span className="text-sm text-muted-foreground">
                  Thinking...
                </span>
              </div>
            )}
            <div ref={messagesEndRef} />
          </div>
        </div>
      </ScrollArea>

      <div className="p-4">
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
            <Button onClick={handleSend} disabled={!input.trim() || isSending}>
              <Send className="size-4" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
