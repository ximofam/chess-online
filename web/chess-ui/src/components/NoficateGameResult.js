import { useEffect, useState } from "react"
import { formatResult } from "../utils/FuncUtils"

const NoficateGameResult = ({ gameResult, handleLeaveRoom }) => {
  const [countdown, setCountdown] = useState(10)

  useEffect(() => {
    if (!gameResult) return

    setCountdown(10)

    const timer = setInterval(() => {
      setCountdown(prev => {
        if (prev <= 1) {
          clearInterval(timer)
          handleLeaveRoom()
          return 0
        }
        return prev - 1
      })
    }, 1000)

    return () => clearInterval(timer)
  }, [gameResult])

  return <>
    <div className="modal fade show d-block">
      <div className="modal-dialog modal-dialog-centered">
        <div className="modal-content">

          <div className="modal-header">
            <h5 className="modal-title">Game Over</h5>
          </div>

          <div className="modal-body text-center">

            <h3>
              {formatResult(gameResult.result, gameResult.method)}
            </h3>

            <p className="text-muted mt-3">
              Leaving room automatically in <b>{countdown}</b> seconds
            </p>

          </div>

          <div className="modal-footer">

            <button
              className="btn btn-danger"
              onClick={handleLeaveRoom}
            >
              Leave Now
            </button>

          </div>

        </div>
      </div>
    </div>

    <div className="modal-backdrop fade show"></div>
  </>
}

export default NoficateGameResult