import { Outlet, Navigate } from "react-router-dom"
import { MessageSquare, Database, LogOut } from "lucide-react"
import { useAuthStore } from "@/stores"
import { Button } from "@/components/ui/button"

export default function MainLayout() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const logout = useAuthStore((s) => s.logout)

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  return (
    <div className="flex h-screen flex-col">
      <header className="flex h-14 items-center justify-between border-b bg-card px-4">
        <div className="flex items-center gap-4">
          <span className="font-bold">Cortex</span>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => (window.location.href = "/chat")}
          >
            <MessageSquare className="mr-1 size-4" />
            Chat
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => (window.location.href = "/skills")}
          >
            <Database className="mr-1 size-4" />
            Skills
          </Button>
          <Button variant="ghost" size="icon" onClick={logout}>
            <LogOut className="size-4" />
          </Button>
        </div>
      </header>
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  )
}
