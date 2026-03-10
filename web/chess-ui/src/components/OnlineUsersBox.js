import { useAuth } from "../context/AuthContext"

const OnlineUsersBox = ({ title, onlineUsers }) => {
  const { user } = useAuth()

  return (
    <div className="col-md-4">
      <div className="card shadow-sm">
        <div className="card-body">
          <h5>{title}</h5>

          <div
            className="mt-3"
            style={{
              maxHeight: "400px",
              overflowY: "auto"
            }}
          >
            <ul className="list-group">
              {onlineUsers.map((u) => u.id !== user.id && (
                <li
                  key={u.id}
                  className="list-group-item d-flex justify-content-between"
                >
                  <span>{u.username}</span>
                  <span className="badge bg-success">
                    online
                  </span>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
    </div>
  )
}

export default OnlineUsersBox