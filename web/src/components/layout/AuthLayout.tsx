import { Outlet, Navigate } from "react-router-dom"
import { useAuthStore } from "@/stores"

export default function AuthLayout() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)

  if (isAuthenticated) {
    return <Navigate to="/chat" replace />
  }

  return <Outlet />
}
