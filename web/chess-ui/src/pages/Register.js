import { useState } from "react";
import Apis from "../configs/Apis";

const Register = () => {
	const [form, setForm] = useState({
		username: "",
		password: "",
		confirmPassword: ""
	});
	const [success, setSuccess] = useState("");
	const [error, setError] = useState("");

	const handleChange = (e) => {
		setForm({
			...form,
			[e.target.name]: e.target.value
		});
	};

	const handleSubmit = async (e) => {
		e.preventDefault();

		if (form.password !== form.confirmPassword) {
			setError("The password doesn't match the confirmation password.");
			return;
		}

		setError("");
		setSuccess("");

		try {
			const res = await Apis.post('/auth/register', {
				username: form.username,
				password: form.password
			});

			setSuccess("Register successful!");
			setForm({
				username: "",
				password: "",
				confirmPassword: ""
			});

		} catch (err) {
			setError(err.response?.data || err.message);
		}
	};
	return (
		<div className="container mt-5">

			<div className="row justify-content-center">
				<div className="col-md-4">

					<div className="card shadow">
						<div className="card-body">

							<h3 className="text-center mb-4">Register</h3>

							<form onSubmit={handleSubmit}>

								<div className="mb-3">
									<label className="form-label">Username</label>
									<input
										type="text"
										name="username"
										className="form-control"
										value={form.username}
										onChange={handleChange}
										required
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
										required
									/>
								</div>

								<div className="mb-3">
									<label className="form-label">Confirm Password</label>
									<input
										type="password"
										name="confirmPassword"
										className="form-control"
										value={form.confirmPassword}
										onChange={handleChange}
										required
									/>
								</div>

								{success && (
									<div className="alert alert-success">
										{success}
									</div>
								)}

								{error && (
									<div className="alert alert-danger">
										{error}
									</div>
								)}

								<button className="btn btn-primary w-100">
									Register
								</button>

							</form>

						</div>
					</div>

				</div>
			</div>

		</div>
	);
}

export default Register;