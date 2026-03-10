import { useState, useRef, useEffect } from "react"
import { formatTime } from "../utils/FuncUtils"

export default function ChatBox({ messages, onSend, height = 300 }) {

  const [text, setText] = useState("")
  const bottomRef = useRef(null)

  const send = () => {

    if (!text.trim()) return

    onSend(text)
    setText("")
  }

  const handleKey = (e) => {
    if (e.key === "Enter") {
      send()
    }
  }

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" })
  }, [messages])

  return (
    <div className="card">

      <div className="card-header">
        Chat
      </div>

      <div
        className="card-body"
        style={{ height: height, overflowY: "auto" }}
      >
        {messages.map((msg, i) => <div key={i} className="mb-1">
          <small className="text-muted me-2">
            {formatTime(msg.send_at)}
          </small>

          <strong>{msg.from}</strong>
          <span className="text-muted">:</span>

          <span className="ms-1">
            {msg.text}
          </span>
        </div>)}

        <div ref={bottomRef}></div>
      </div>

      <div className="card-footer d-flex">

        <input
          className="form-control"
          placeholder="Type message..."
          value={text}
          onChange={e => setText(e.target.value)}
          onKeyDown={handleKey}
        />

        <button
          className="btn btn-primary ms-2"
          onClick={send}
        >
          Send
        </button>

      </div>

    </div>
  )
}