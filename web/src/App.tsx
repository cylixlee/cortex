import { useEffect } from "react"
import { RouterProvider } from "react-router-dom"
import { Toaster } from "sonner"
import router from "@/router"
import { useAuthStore } from "@/stores"

export function App() {
  const checkAuth = useAuthStore((s) => s.checkAuth)

  useEffect(() => {
    checkAuth()
  }, [checkAuth])

  return (
    <>
      <RouterProvider router={router} />
      <Toaster position="top-right" />
    </>
  )
}

export default App
