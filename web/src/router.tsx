import { createBrowserRouter, Navigate } from "react-router-dom"
import { AuthLayout, MainLayout } from "@/components/layout"
import {
  LoginPage,
  RegisterPage,
  ChatPage,
  SkillListPage,
  SkillUploadPage,
  SkillDetailPage,
} from "@/pages"

const router = createBrowserRouter([
  {
    element: <AuthLayout />,
    children: [
      {
        path: "/login",
        element: <LoginPage />,
      },
      {
        path: "/register",
        element: <RegisterPage />,
      },
    ],
  },
  {
    element: <MainLayout />,
    children: [
      {
        path: "/",
        element: <Navigate to="/chat" replace />,
      },
      {
        path: "/chat",
        element: <ChatPage />,
      },
      {
        path: "/chat/:id",
        element: <ChatPage />,
      },
      {
        path: "/skills",
        element: <SkillListPage />,
      },
      {
        path: "/skills/upload",
        element: <SkillUploadPage />,
      },
      {
        path: "/skills/:id",
        element: <SkillDetailPage />,
      },
    ],
  },
])

export default router
