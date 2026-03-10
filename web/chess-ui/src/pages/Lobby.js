import { useEffect, useState } from "react"
import { useSocket } from "../context/SocketContext"
import OnlineUsersBox from "../components/OnlineUsersBox"
import { useNavigate } from "react-router-dom"


const Lobby = () => {
	const navigate = useNavigate()
	const { send, addEventListener, connected } = useSocket()
	const [showCreateRoom, setShowCreateRoom] = useState(false)
	const [createRoomForm, setCreateRoomForm] = useState({
		name: "",
		allow_spectate: false,
		max_spectators: 0
	})
	const [rooms, setRooms] = useState([])
	const [onlineUsers, setOnlineUsers] = useState([])

	const handlers = {
		INFO_LIST: (data) => {
			setOnlineUsers(data.online_users)
			setRooms(data.rooms)
		},
		USER_ONLINE: (data) => {
			setOnlineUsers(prev => [...prev, data])
		},
		USER_LEAVE: (data) => {
			setOnlineUsers(prev => prev.filter(u => u.id !== data.id))
		},
		ROOM_CREATE: (data) => setRooms(prev => [...prev, data]),
		ROOM_DELETE: (data) => setRooms(prev => prev.filter(x => x.id !== data.id)),
		ROOM_UPDATE: (data) => setRooms(prev => prev.map(x => x.id === data.id ? data : x)),
		ROOM_JOIN_SUCCESS: (data) => navigate('/chess', { state: { role: data } })
	}

	useEffect(() => {
		if (!connected) return

		const unSub = addEventListener("message", (event) => {
			const msg = JSON.parse(event.data)

			handlers[msg.type]?.(msg.data)
		})

		send({ event: "INFO_LIST" })

		return () => {
			unSub()
		}
	}, [connected])


	const handleFormChange = (e) => {
		const { name, value, type, checked } = e.target

		setCreateRoomForm(prev => ({
			...prev,
			[name]:
				type === "checkbox"
					? checked
					: type === "number"
						? Number(value)
						: value
		}))
	}

	const handleCreateRoom = () => {
		send({
			event: "ROOM_CREATE",
			payload: createRoomForm
		})
	}

	const handleJoin = (roomId) => {
		send({
			event: "ROOM_JOIN",
			payload: {
				room_id: roomId,
				role: "Black"
			}
		})
	}

	const handleWatch = (roomId) => {
		send({
			event: "ROOM_JOIN",
			payload: {
				room_id: roomId,
				role: "Spectator"
			}
		})
	}

	return (
		<div className="container-fluid mt-4">
			<div className="row">

				<div className="col-md-8 d-flex flex-column" style={{ height: "calc(100vh - 60px)" }}>

					<div className="mb-3 d-flex justify-content-between align-items-center">
						<h3 className="mb-0">♟ Active Rooms</h3>

						<button
							className="btn btn-success"
							onClick={() => setShowCreateRoom(true)}
						>
							Create Room
						</button>

						{showCreateRoom && (
							<div className="modal fade show d-block">
								<div className="modal-dialog modal-dialog-centered">
									<div className="modal-content">

										<div className="modal-header">
											<h5 className="modal-title">Create Room</h5>
											<button
												className="btn-close"
												onClick={() => setShowCreateRoom(false)}
											/>
										</div>

										<div className="modal-body">

											<div className="mb-3">
												<label className="form-label">Room Name</label>
												<input
													className="form-control"
													name="name"
													value={createRoomForm.name}
													onChange={handleFormChange}
												/>
											</div>

											<div className="form-check mb-3">
												<input
													className="form-check-input"
													type="checkbox"
													name="allow_spectate"
													checked={createRoomForm.allow_spectate}
													onChange={handleFormChange}
												/>
												<label className="form-check-label">
													Allow Spectators
												</label>
											</div>

											<div className="mb-3">
												<label className="form-label">Max Spectators</label>
												<input
													type="number"
													className="form-control"
													name="max_spectators"
													value={createRoomForm.max_spectators}
													onChange={handleFormChange}
													disabled={!createRoomForm.allow_spectate}
												/>
											</div>

										</div>

										<div className="modal-footer">
											<button
												className="btn btn-secondary"
												onClick={() => setShowCreateRoom(false)}
											>
												Cancel
											</button>

											<button
												className="btn btn-success"
												onClick={handleCreateRoom}
											>
												Create
											</button>
										</div>

									</div>
								</div>
							</div>
						)}
					</div>

					{!connected && (
						<p className="text-danger">Disconnected...</p>
					)}

					{/* Scroll container */}
					<div
						style={{
							flex: 1,
							overflowY: "auto",
							paddingRight: "8px"
						}}
					>
						{rooms.map((room) => {
							const isUnlimited = room.max_spectators === 0
							const isFull =
								!isUnlimited &&
								room.curr_observer >= room.max_spectators

							return (
								<div key={room.id} className="card shadow-sm mb-3">
									<div className="card-body d-flex justify-content-between align-items-center">

										{/* Room Info */}
										<div>
											<h5 className="mb-1">{room.name}</h5>

											<div className="text-muted small">
												♔ {room.white} vs ♚ {room.black || "Waiting..."}
											</div>

											<div className="small">
												Spectators: {room.curr_observer}
												{isUnlimited
													? " / ∞"
													: ` / ${room.max_spectators}`}
											</div>
										</div>

										{/* Buttons */}
										<div className="d-flex gap-2">
											<button
												className="btn btn-success btn-sm"
												disabled={!room.join_able}
												onClick={() => handleJoin(room.id)}
											>
												{room.join_able ? "Join" : "Closed"}
											</button>

											<button
												className="btn btn-primary btn-sm"
												disabled={!room.allow_spectate || isFull}
												onClick={() => handleWatch(room.id)}
											>
												{!room.allow_spectate ? "Spectator Off" : isFull ? "Full" : "Watch"}
											</button>
										</div>

									</div>
								</div>
							)
						})}
					</div>
				</div>

				{/* RIGHT: ONLINE USERS */}
				<OnlineUsersBox title="Online users" onlineUsers={onlineUsers} />

			</div>
		</div>
	)
}

export default Lobby