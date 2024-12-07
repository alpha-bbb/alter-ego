import { useEffect, useState } from "react";
import liff from "@line/liff";
import "./App.css";
import { Button } from "@mui/material";
import SendIcon from "@mui/icons-material/Send";
import picture from "./images/alter-ego-pic.png";

function App() {
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    // liff.ready.then(() => {
    liff
      .init({
        liffId: import.meta.env.VITE_LIFF_ID
      })
      .then(() => {
        setMessage("LIFF init succeeded.");
      })
      .catch((e: Error) => {
        setMessage("LIFF init failed.");
        setError(`${e}`);
      });
    // })
  });

  const handleShareTargetPicker = () => {
    if (liff.isApiAvailable("shareTargetPicker")) {
      liff.shareTargetPicker([
        {
          type: "text",
          text: "Hello, World!",
        },
      ]);
    } else {
      //「シェアターゲットピッカー」が無効になっている場合
      setMessage("invalid")
    }
  };

  return (
    <div className="App">
      <h1>Alter Ego</h1>
      <img src={picture} alt="decoration" style={{ width: '50%', borderRadius: '10px' }} />
      {message && <p>{message}</p>}
      {error && (
        <p>
          <code>{error}</code>
        </p>
      )}
      <Button onClick={handleShareTargetPicker} variant="contained" endIcon={<SendIcon />} size="large" sx={{ borderRadius: '20px', backgroundColor: '#469ac8', color: '#fff', width: '80%'  }}>
        送信先の選択
      </Button>
    </div>
  );
}

export default App;
