import { Link } from "react-router-dom"
import { useAuth } from "../../context/AuthContext"

const Header = () => {
	const { user, logout } = useAuth()

	return (
		<nav className="navbar navbar-expand-lg navbar-dark bg-dark px-4">
			<Link className="navbar-brand fw-bold logo" to="/">
				♟ Chess Online
			</Link>

			{user &&
				<div className="navbar-nav ms-4">
					<Link className="nav-link custom-link" to="/lobby">
						Lobby
					</Link>

					<Link className="nav-link custom-link" to="/history">
						Match History
					</Link>
				</div>}

			<div className="ms-auto d-flex align-items-center gap-3">
				{user ? (
					<>
						<div className="d-flex align-items-center gap-2 text-white">
							<img
								src="https://www.gravatar.com/avatar/?d=mp"
								alt="avatar"
								width="35"
								height="35"
								className="rounded-circle border"
							/>
							<span className="fw-semibold">
								{user.username}
							</span>
						</div>

						<button
							className="btn btn-outline-light btn-sm"
							onClick={logout}
						>
							Logout
						</button>
					</>
				) : (
					<Link className="btn btn-outline-light btn-sm" to="/login">
						Login
					</Link>
				)}
			</div>
		</nav>
	)
}

export default Header