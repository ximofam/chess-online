import { useState } from "react"
import { fenToBoard } from "../utils/FuncUtils"

function ChessBoard({ fen, playerRole, turn, validMoves, isPlaying, handleMove }) {
  const [selected, setSelected] = useState(null)
  const pieceMoves = selected
    ? validMoves.filter(m => m.startsWith(selected))
    : []
  const board = fenToBoard(fen)

  const files = ["a", "b", "c", "d", "e", "f", "g", "h"]
  function toSquare(r, c) {
    return files[c] + (8 - r)
  }

  const dir = '/pieces/images/'
  const pieceImages = {
    p: `${dir}Chess_pdt60.png`,
    r: `${dir}Chess_rdt60.png`,
    n: `${dir}Chess_ndt60.png`,
    b: `${dir}Chess_bdt60.png`,
    q: `${dir}Chess_qdt60.png`,
    k: `${dir}Chess_kdt60.png`,

    P: `${dir}Chess_plt60.png`,
    R: `${dir}Chess_rlt60.png`,
    N: `${dir}Chess_nlt60.png`,
    B: `${dir}Chess_blt60.png`,
    Q: `${dir}Chess_qlt60.png`,
    K: `${dir}Chess_klt60.png`,
  }

  function handleSquareClick(r, c) {
    if (!isPlaying) return

    const square = toSquare(r, c)
    const piece = board[r][c]

    if (selected) {
      const move = selected + square
      const isValid = pieceMoves.some(m => m === move)

      if (isValid) {
        handleMove(move)
        setSelected(null)
        return
      }
    }

    if (piece) {
      if (playerRole === "white" && piece !== piece.toUpperCase()) return
      if (playerRole === "black" && piece !== piece.toLowerCase()) return

      if (turn !== playerRole) return

      setSelected(square)
    }
  }

  return (
    <div
      style={{
        width: "560px",
        border: "3px solid #333"
      }}
    >
      {board.map((row, r) => (
        <div className="d-flex" key={r}>

          {row.map((sq, c) => {

            const isDark = (r + c) % 2 === 1

            const square = toSquare(r, c)
            const isLegal = pieceMoves.some(m => m.slice(2, 4) === square)
            const isSelected = selected === square

            return (
              <div
                key={`${r}-${c}`}
                onClick={() => handleSquareClick(r, c)}
                style={{
                  width: "70px",
                  height: "70px",
                  backgroundColor:
                    isSelected
                      ? "#f6f669"
                      : isLegal
                        ? "#baca44"
                        : isDark
                          ? "#769656"
                          : "#eeeed2",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center"
                }}
              >

                {sq && (
                  <img
                    src={pieceImages[sq]}
                    alt={sq}
                    style={{
                      width: "55px",
                      height: "55px"
                    }}
                  />
                )}

              </div>
            )

          })}

        </div>
      ))}
    </div>
  )
}

export default ChessBoard