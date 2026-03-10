import { useEffect, useState } from "react"
import ChessBoard from "../components/ChessBoard"
import PlayerInfo from "../components/PlayerInfo"
import ChatBox from "../components/ChatBox"
import { useLocation, useNavigate } from "react-router-dom"
import { useSocket } from "../context/SocketContext"
import NoficateGameResult from "../components/NoficateGameResult"

function Chess() {
	const location = useLocation()
	const navigate = useNavigate()
	const { send, addEventListener, connected } = useSocket()
	const [fen, setFen] = useState("")
	const [msgs, setMsgs] = useState([])
	const [white, setWhite] = useState(null)
	const [whiteTime, setWhiteTime] = useState(0)
	const [black, setBlack] = useState(null)
	const [blackTime, setBlackTime] = useState(0)
	const [spectators, setSpectators] = useState([])
	const [turn, setTurn] = useState('')
	const [isPlaying, setIsPlaying] = useState(false)
	const [validMoves, setValidMoves] = useState([])
	const [gameResult, setGameResult] = useState(null)
	const playerRole = location.state?.role


	const handlers = {
		INFO_LIST: (data) => {
			setWhite(data.white)
			setBlack(data.black)
			setSpectators(data.spectators)
			setMsgs(data.chat_messages)
			setIsPlaying(data.is_playing)
			const gameState = data.game_state
			setFen(gameState.fen)
			setWhiteTime(gameState.white_time)
			setBlackTime(gameState.black_time)
			setTurn(gameState.turn)
		},
		ROOM_DELETE: (data) => { navigate('/lobby') },
		PLAYER_JOIN: (data) => {
			switch (data.role) {
				case 'White':
					setWhite({ id: data.id, username: data.username, time: 600 })
					break
				case 'Black':
					setBlack({ id: data.id, username: data.username, time: 600 })
					break
				case 'Spectator':
					setSpectators(prev => [...prev, data])
			}
		},
		PLAYER_LEAVE: (data) => {
			switch (data.role) {
				case 'Black':
					setBlack(null)
				case 'Spectator':
					setSpectators(prev => prev.filter(x => x.id !== data.id))
			}
		},
		CHAT: (data) => setMsgs(prev => [...prev, data]),
		PLAY: (data) => {
			setIsPlaying(true)
		},
		GAME_STATE: (data) => {
			setFen(data.fen)
			setWhiteTime(data.white_time)
			setBlackTime(data.black_time)
			setTurn(data.turn)
			setValidMoves(data.valid_moves)
		},
		GAMEOVER: (data) => {
			setIsPlaying(false)
			setGameResult(data)
		}
	}

	useEffect(() => {
		if (!connected) return

		const unSub = addEventListener("message", (event) => {
			const msg = JSON.parse(event.data)

			handlers[msg.type]?.(msg.data)
		})

		send({ event: "INFO_LIST" })

		return () => {
			send({ event: 'LEAVE' })
			unSub()
		}
	}, [connected])

	useEffect(() => {
		if (!isPlaying) return

		const interval = setInterval(() => {
			if (turn === "White") {
				setWhiteTime(t => Math.max(t - 1e9, 0))
			} else {
				setBlackTime(t => Math.max(t - 1e9, 0))
			}
		}, 1000)

		return () => clearInterval(interval)
	}, [isPlaying, turn])

	const handleChat = (msg) => {
		send({
			event: "CHAT",
			payload: {
				text: msg
			}
		})
	}

	const handleLeaveRoom = () => {
		send({ event: "LEAVE" })
		navigate('/lobby')
	}


	return (
		<div className="container mt-4">
			{(playerRole !== 'White' || playerRole !== 'Black' || !isPlaying) &&
				<button class="btn btn-danger position-absolute top-0 end-0 m-3" onClick={handleLeaveRoom}>
					Leave
				</button>}

			{gameResult && <NoficateGameResult gameResult={gameResult} handleLeaveRoom={handleLeaveRoom} />}

			<div className="row">

				{/* LEFT: BOARD */}
				<div className="col-md-8 text-center">

					<PlayerInfo player={black} color="Black" time={blackTime} />

					<ChessBoard fen={fen}
						isPlaying={isPlaying}
						playerRole={playerRole}
						turn={turn}
						validMoves={validMoves}
						handleMove={(move) =>
							send({
								event: 'MOVE',
								payload: { uci: move }
							})} />

					<PlayerInfo player={white} color="White" time={whiteTime} />

					{!isPlaying && playerRole === 'White' &&
						<button onClick={() => { send({ event: 'PLAY' }) }}>Play</button>}

				</div>

				{/* RIGHT PANEL */}
				<div className="col-md-4 d-flex flex-column" style={{ height: "600px" }}>
					{/* Spectators */}
					<div className="card p-2 mb-2 d-flex flex-column" style={{ height: "50%" }}>
						<h6>👀 Spectators ({spectators.length})</h6>

						<div style={{ overflowY: "auto", flex: 1 }}>
							{spectators.map((s) => (
								<div key={s.id}>
									{s.username}
								</div>
							))}
						</div>
					</div>

					{/* Chat */}
					<div style={{ height: "50%" }}>
						<ChatBox
							messages={msgs}
							onSend={handleChat}
						/>
					</div>
				</div>

			</div>

		</div>
	)
}

export default Chess