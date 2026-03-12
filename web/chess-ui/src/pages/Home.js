import { useAuth } from "../context/AuthContext"

const Home = () => {
	return (
		<div className="container text-center mt-5">
			<h1 className="display-4 fw-bold">Welcome to Chess Online</h1>
			<p className="lead mt-3">
				This project uses Go for the backend, WebSocket for real-time communication,
				and ReactJS for the frontend.
			</p>
		</div>
	)
}

export default Home