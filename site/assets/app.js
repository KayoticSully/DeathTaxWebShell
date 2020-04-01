var appendNext = false;
var apiSocket;
var inputKeys = [];
const inputKeyPattern = /\[(\S)\]\s\S*/g;
const inputLinePattern = /\(default is \"\S\"\):\s*$/;

document.addEventListener('DOMContentLoaded', function() {
    // Connect to the api
    apiSocket = connect();
    document.addEventListener('keyup', handeKeyEvent);
}, false);

function connect() {
    let ws = new WebSocket("wss://deathtax.kayotic.io/api");
    ws.onmessage = handleMessage;
    return ws;
}

// Setup message handling
function handleMessage(msg) {
    let data = msg.data;

    if(appendNext) {
        let elem = document.querySelector("#output .line:last-child");
        elem.innerHTML += data;
        appendNext = false;

        // force a newline after an append
        data = "\n";
    }

    let elem = document.querySelector("#output");

    for(line of data.split("\n")) {
        if(line.match(inputLinePattern)) {
            inputKeys = Array.from(line.matchAll(inputKeyPattern)).map(match => match[1]);
        }

        if(line == "") {
            line = "&nbsp;";
        }

        elem.innerHTML += `<div class="line">${line}</div>`;
        elem.scrollTop = elem.scrollHeight;
    }
}

function handeKeyEvent(event) {
    if (event.defaultPrevented) {
        return;
    }
    
    // Prevent any further input than the one keystroke
    let allowedKeys = Array.from(inputKeys);
    inputKeys = [];
    
    var key = event.key || event.keyCode;
    if (allowedKeys.includes(key.toUpperCase())) {
        appendNext = true;
        apiSocket.send(`${key}\n`);
    }
}