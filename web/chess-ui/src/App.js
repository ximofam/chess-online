import { BrowserRouter, Route, Routes } from "react-router-dom"
import { AuthProvider } from "./context/AuthContext"
import { SocketProvider } from "./context/SocketContext"
import ProtectedRoute from "./components/ProtectedRoute"
import Home from "./pages/Home"
import Lobby from "./pages/Lobby"
import Login from "./pages/Login"
import Register from "./pages/Register"
import Chess from "./pages/Chess"
import Main from "./components/layouts/Main"
import HistoryMatch from "./pages/HistoryMatch"

const App = () => {
  return (
    <AuthProvider>
      <SocketProvider>
        <BrowserRouter>
          <Routes>
            <Route element={<Main />}>
              <Route path="/" element={<Home />} />
              <Route path="/lobby" element={<ProtectedRoute><Lobby /></ProtectedRoute>} />
              <Route path="/history" element={<ProtectedRoute><HistoryMatch /></ProtectedRoute>} />
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
            </Route>

            <Route path="/chess" element={<ProtectedRoute><Chess /></ProtectedRoute>} />

          </Routes>
        </BrowserRouter>
      </SocketProvider>
    </AuthProvider>
  )
}

export default App