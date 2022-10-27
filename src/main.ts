import { editor } from "monaco-editor";
import { Store, defaults } from "./default";

const store: Store = getStorageFromBrowser()
  ? getStorageFromBrowser()
  : defaults;

const editorInstance = editor.create(document.getElementById("editor"), {
  value: store.javascript,
  language: "javascript",
  theme: "vs-dark",
});
window.onresize = function () {
  editorInstance.layout();
};

console.log("Hello World");

let loading = false;
editorInstance.getModel().onDidChangeContent((e) => {
  if (loading) return;
  loading = true;
  setTimeout(() => {
    saveStorageInBrowser();
    loading = false;
  }, 1000);
});

document.getElementById("save").addEventListener("click", () => {
  const blob = new Blob([JSON.stringify(store)], { type: "text/plain" });
  const a = document.createElement("a");
  a.href = URL.createObjectURL(blob);
  a.download = "store.json";
  a.click();
});

document.getElementById("import").addEventListener("click", () => {
  const input = document.createElement("input");
  input.type = "file";
  input.accept = ".json";
  input.onchange = (e) => {
    const file = (e.target as HTMLInputElement).files[0];
    const reader = new FileReader();
    reader.onload = (e) => {
      const text = e.target.result as string;
      const store = JSON.parse(text);
      editorInstance.setValue(store.javascript);
    };
    reader.readAsText(file);
  };
  input.click();
});

function saveStorageInBrowser() {
  localStorage.setItem("store", JSON.stringify(store));
}

function getStorageFromBrowser() {
  const store = localStorage.getItem("store");
  if (store) return JSON.parse(store);
  return null;
}


document.getElementById("run").addEventListener("click", run);

function run() {
  const config = getConfig();
  const code = editorInstance.getValue();

  fetch("/run", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      code,
      config,
    }),
  })
    .then((res) => res.json())
    .then((res) => {
      document.getElementById("result").innerText = JSON.stringify(
        res,
        null,
        2
      );
      console.log(res);
    });
}

var url = new URL('/ws', window.location.href);
url.protocol = url.protocol.replace('http', 'ws');
const socket = new WebSocket(url.href);

socket.addEventListener("open", (event) => {
  console.log("Sending message to server");

  socket.send(
    JSON.stringify({
      data: "Hello Server! I am a client. I am sending you a message. I hope you get it.",
    })
  );
});
socket.addEventListener("close", (event) => {
  console.log("Connection closed", event);
});
socket.addEventListener("message", (event) => {
  const data = JSON.parse(event.data);
  console.log("Message from server ", data);
});

interface Config {
  nodes: number;
  msgs_s_node: number;
  datasets: number;
  datasets_s: number;
  seeds: number;
  iterations: number;
  timeout: number;
}

function getConfig() : Config | Error {
  const configEl = document.getElementById("config");

  const config: Config = {
    nodes: 0,
    msgs_s_node: 0,
    datasets: 0,
    datasets_s: 0,
    seeds: 0,
    iterations: 0,
    timeout: 0,
  };
  for (const el of configEl.querySelectorAll("input")) {
    const key = el.getAttribute("data-key");
    const value = parseInt(el.value);
    if (typeof value !== "number" || isNaN(value)) throw new Error("Invalid value");
    if ((config as any)[key] === undefined) throw new Error("Invalid key");
    
    (config as any)[key] = value;
  }

  return config;
}

console.log(getConfig());