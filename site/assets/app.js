
function displayOutput(msg) {
    let elem = document.querySelector("#output");
    elem.innerHTML += `<div class="line">${msg.data}</div>`;
}

function connect() {
    let ws = new WebSocket("wss://deathtax.kayotic.io/api");
    ws.onmessage = displayOutput;
    return ws;
}

let apiSocket = connect();
