import { createContext, useContext, useEffect, useState } from "react"
import Apis from "../configs/Apis"

const AuthContext = createContext(null)

export const accessTokenKey = "accessToken"

export const AuthProvider = ({ children }) => {
	const [user, setUser] = useState(null)


	useEffect(() => {
		const accessToken = localStorage.getItem(accessTokenKey)
		if (!accessToken) return

		Apis.get("/auth/me", {
			headers: {
				Authorization: `Bearer ${accessToken}`
			}
		}).then(res => setUser(res.data))
			.catch(err => console.log(err))
	}, [])

	const login = (userData, accessTokenData) => {
		setUser(userData)
		localStorage.setItem(accessTokenKey, accessTokenData)
	}

	const logout = () => {
		setUser(null)
		localStorage.removeItem(accessTokenKey)
	}

	return (
		<AuthContext.Provider value={{ user, login, logout }}>
			{children}
		</AuthContext.Provider>
	)

}

export const useAuth = () => useContext(AuthContext)