import liff from "@line/liff";
import { useEffect, useState } from "react";
import "./App.css";

function App() {
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [text, setText] = useState("");

  useEffect(() => {
    const initLiff = async () => {
      try {
        await liff.init({
          liffId: import.meta.env.VITE_LIFF_ID,
        });

        setMessage("LIFF init succeeded.");

        // クエリパラメータを取得
        const queryParams = new URLSearchParams(window.location.search);
        const queryText = queryParams.get("text");
        const messageNo = queryParams.get("messageNo");

        if (queryText) {
          // setText(queryText);
          setText("hogehoge");
          console.log(`Query parameter "text": ${queryText}`);
        }

        if (messageNo) {
          console.log(`Message No.: ${messageNo}`);
        }
      } catch (e) {
        setMessage("LIFF init failed.");
        setError(`${e}`);
      }
    };

    initLiff();
  }, []);

  const handleCopyToClipboard = async () => {
    try {
      if (text) {
        await navigator.clipboard.writeText("hogehoge");
        alert(`Copied to clipboard: ${text}`);
      } else {
        alert("No text available to copy.");
      }
    } catch (error) {
      console.error("Failed to copy text to clipboard:", error);
      alert("Clipboard operation failed.");
    }
  };

  return (
    <div className="App">
      <h1>LIFF App</h1>
      {message && <p>{message}</p>}
      {error && (
        <p>
          <code>{error}</code>
        </p>
      )}
      <div>
        <p>Text to copy: {text || "No text provided"}</p>
        {/* biome-ignore lint/a11y/useButtonType: <explanation> */}
        <button onClick={handleCopyToClipboard}>Copy to Clipboard</button>
      </div>
    </div>
  );
}

export default App;
