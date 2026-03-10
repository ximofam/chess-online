export function fenToBoard(fen) {
  const position = fen.split(" ")[0]
  const rows = position.split("/")

  return rows.map((row) => {
    const squares = []

    for (const char of row) {
      if (!isNaN(char)) {
        const empty = Number(char)
        for (let i = 0; i < empty; i++) {
          squares.push("")
        }
      } else {
        squares.push(char)
      }
    }

    return squares
  })
}

export function formatTime(time) {
  const date = new Date(time)
  return date.toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit"
  })
}

export function formatResult(result, method) {
  let winner = ""

  if (result === "1-0") winner = "White Wins"
  else if (result === "0-1") winner = "Black Wins"
  else winner = "Draw"

  const methods = {
    checkmate: "by Checkmate",
    resign: "by Resignation",
    timeout: "on Time",
    stalemate: "by Stalemate",
    draw_agreement: "by Agreement",
    threefold_repetition: "by Threefold Repetition",
    insufficient_material: "by Insufficient Material"
  }

  return `${winner} ${methods[method] || ""}`
}