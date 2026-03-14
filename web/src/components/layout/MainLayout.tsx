import { useEffect, useState } from "react"
import { Outlet, Navigate, useNavigate, useParams } from "react-router-dom"
import { Plus, Trash2, Database, LogOut } from "lucide-react"
import { useAuthStore, useChatStore } from "@/stores"
import { Button } from "@/components/ui/button"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { cn } from "@/lib/utils"

export default function MainLayout() {
  const navigate = useNavigate()
  const { id } = useParams()
  const { user, isAuthenticated, logout } = useAuthStore()
  const {
    conversations,
    currentConversation,
    loadConversations,
    deleteConversation,
    clearCurrentConversation,
  } = useChatStore()

  const [deleteId, setDeleteId] = useState<string | null>(null)

  useEffect(() => {
    loadConversations()
  }, [loadConversations])

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  const handleNewChat = () => {
    clearCurrentConversation()
    navigate("/chat")
  }

  const handleDelete = (convId: string) => {
    setDeleteId(convId)
  }

  const confirmDelete = async () => {
    if (!deleteId) return
    await deleteConversation(deleteId)
    if (deleteId === currentConversation?.id) {
      navigate("/chat")
    }
    setDeleteId(null)
  }

  return (
    <div className="flex h-screen">
      <aside className="flex w-64 flex-col border-r bg-card">
        <div className="flex flex-col gap-2 p-4">
          <Button
            variant="ghost"
            size="sm"
            className="w-full justify-start gap-2"
            onClick={handleNewChat}
          >
            <Plus className="size-4" />
            New Chat
          </Button>
          <Button
            variant="default"
            size="sm"
            className="w-full justify-start gap-2"
            onClick={() => navigate("/skills")}
          >
            <Database className="size-4" />
            My Skills
          </Button>
        </div>
        <ScrollArea className="w-full flex-1">
          <div className="flex w-full flex-col gap-1 p-2">
            {conversations.map((conv) => (
              <div
                key={conv.id}
                onClick={() => navigate(`/chat/${conv.id}`)}
                className={cn(
                  "flex cursor-pointer items-center gap-2 rounded-md px-3 py-2 text-sm hover:bg-accent",
                  id === conv.id && "bg-accent"
                )}
              >
                <span className="truncate-sm">{conv.title}</span>
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-6 opacity-50 hover:opacity-100"
                  onClick={(e) => {
                    e.stopPropagation()
                    handleDelete(conv.id)
                  }}
                >
                  <Trash2 className="size-3" />
                </Button>
              </div>
            ))}
          </div>
        </ScrollArea>
        <div className="p-4">
          <div className="flex items-center gap-3">
            <Avatar className="size-8">
              <AvatarFallback>
                {user?.email?.[0]?.toUpperCase() || "U"}
              </AvatarFallback>
            </Avatar>
            <div className="min-w-0 flex-1">
              <p className="truncate text-sm font-medium">{user?.email}</p>
            </div>
            <Button variant="ghost" size="icon" onClick={logout}>
              <LogOut className="size-4" />
            </Button>
          </div>
        </div>
      </aside>
      <main className="min-h-0 flex-1">
        <div className="relative mx-auto h-full max-w-5xl p-6">
          <Outlet />
        </div>
      </main>

      <AlertDialog open={!!deleteId} onOpenChange={() => setDeleteId(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete this conversation?</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={confirmDelete}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
