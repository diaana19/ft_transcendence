import { BrowserRouter, Routes, Route } from "react-router-dom"
import Layout from "../components/layout/Layout"
import Feed from "../features/posts/Feed"
import Profile from "../features/user/Profile"

export default function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>

        {/* Rutes of Layout */}
        <Route element={<Layout />}>
          <Route path="/" element={<Feed />} />
          <Route path="/profile" element={<Profile />} />
        </Route>

      </Routes>
    </BrowserRouter>
  )
}