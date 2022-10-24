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
function getConfig() {
  const configEl = document.getElementById("config");

  const config: any = {}; //TODO fix type
  for (const el of configEl.querySelectorAll("input")) {
    const key = el.getAttribute("data-key");
    const value = el.value;
    config[key] = value;
  }

  return config;
}

document.getElementById("run").addEventListener("click", run);

function run() {
  const config = getConfig();
  const code = editorInstance.getValue();

  fetch("http://localhost:3333/run", {
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
