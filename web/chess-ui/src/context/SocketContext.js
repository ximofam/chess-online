import { createContext, useContext, useEffect, useRef, useState } from "react"
import { accessTokenKey, useAuth } from "./AuthContext"

const SocketContext = createContext(null)

export const SocketProvider = ({ children }) => {
	const { user } = useAuth()
	const socketRef = useRef(null);
	const [connected, setConnected] = useState(false);

	const connect = () => {
		const accessToken = localStorage.getItem(accessTokenKey)
		if (!accessToken) return

		const ws = new WebSocket(`ws://localhost:8080/ws?token=${accessToken}`);

		socketRef.current = ws;

		ws.onopen = () => {
			console.log("WS connected");
			setConnected(true);
		};

		ws.onclose = () => {
			console.log("WS closed");

			setConnected(false);

			setTimeout(() => {
				connect();
			}, 3000);
		};
	};

	useEffect(() => {
		if (!user) return

		connect();

		return () => {
			socketRef.current?.close();
		};

	}, [user]);

	const send = (data) => {
		if (socketRef.current?.readyState === 1) {
			socketRef.current.send(JSON.stringify(data));
		}
	};

	const addEventListener = (type, handler) => {
		socketRef?.current.addEventListener(type, handler)

		return () => {
			socketRef?.current.removeEventListener(type, handler)
		}
	}

	const value = {
		send,
		addEventListener,
		connected
	};

	return (
		<SocketContext.Provider value={value}>
			{children}
		</SocketContext.Provider>
	);
}
export const useSocket = () => useContext(SocketContext)