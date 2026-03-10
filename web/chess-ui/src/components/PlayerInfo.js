const PlayerInfo = ({ player, color, time }) => {

  function formatTime(ns) {
    const totalSeconds = Math.floor(ns / 1e9)

    const minutes = Math.floor(totalSeconds / 60)
    const seconds = totalSeconds % 60

    return `${minutes}:${seconds.toString().padStart(2, "0")}`
  }

  return (
    <div className="card p-2 mb-2">
      <div className="d-flex justify-content-between align-items-center">

        <div>
          {player ? <strong>{player.username}</strong> : <strong>Waiting...</strong>}
          <span className="ms-2 badge bg-secondary">
            {color}
          </span>
        </div>

        <div className="fs-5">
          {formatTime(time)}
        </div>

      </div>
    </div>
  )
}

export default PlayerInfo