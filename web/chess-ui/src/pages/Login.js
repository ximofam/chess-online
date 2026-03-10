import { useState } from "react"
import { useNavigate, Link } from "react-router-dom"
import Apis, { endpoint } from "../configs/Apis"
import { useAuth } from "../context/AuthContext"

export default function Login() {
	const navigate = useNavigate()
	const { login } = useAuth()
	const [form, setForm] = useState({
		username: "",
		password: ""
	})
	const [error, setError] = useState("")
	const [loading, setLoading] = useState(false)

	const handleChange = (e) => {
		setForm({
			...form,
			[e.target.name]: e.target.value
		})
	}

	const handleSubmit = async (e) => {
		e.preventDefault()
		setError("")

		try {
			const res = await Apis.post("/auth/login", {
				username: form.username,
				password: form.password
			})

			const accessToken = res.data.access_token

			const userRes = await Apis.get("/auth/me", {
				headers: {
					Authorization: `Bearer ${accessToken}`
				}
			})

			login(userRes.data, accessToken)

			navigate("/lobby")
		} catch (err) {
			setError(err.response?.data || err.message)
		}
	}

	return (
		<div className="container d-flex justify-content-center align-items-center vh-100">
			<div className="card shadow p-4" style={{ width: "400px" }}>

				<h3 className="text-center mb-4">Login to Chess Game</h3>

				{error && (
					<div className="alert alert-danger">
						{error}
					</div>
				)}

				<form onSubmit={handleSubmit}>

					<div className="mb-3">
						<label className="form-label">Username</label>
						<input
							type="username"
							name="username"
							className="form-control"
							value={form.username}
							onChange={handleChange}
						/>
					</div>

					<div className="mb-3">
						<label className="form-label">Password</label>
						<input
							type="password"
							name="password"
							className="form-control"
							value={form.password}
							onChange={handleChange}
						/>
					</div>

					<button
						type="submit"
						className="btn btn-warning w-100"
						disabled={loading}
					>
						{loading ? "Logging in..." : "Login"}
					</button>

				</form>

				<p className="text-center mt-3">
					Don't have an account?{" "}
					<Link to="/register">Register</Link>
				</p>

			</div>
		</div>
	)
}