var apiSocket;
var inputKeys = [];
const inputKeyPattern = /\[(\S)\]\s\S*/g;
const inputLinePattern = /\(default is \"\S\"\):\s*$/;
const lastLineSelector = "#output .line:last-child";

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
    for(line of msg.data.split("\n")) {
        if(line.match(inputLinePattern)) {
            inputKeys = Array.from(line.matchAll(inputKeyPattern)).map(match => match[1]);
        } else if(line.length == 1) {
            let lastLineText = document.querySelector(lastLineSelector).innerHTML;

            if(line == lastLineText[lastLineText.length-1]) {
                line = "&nbsp;";
            }
        } else if(line == "") {
            line = "&nbsp;";
        }

        addLine(line);
    }
}

function addLine(line) {
    let elem = document.querySelector("#output");
    elem.innerHTML += `<div class="line">${line}</div>`;
    elem.scrollTop = elem.scrollHeight;
}

function appendToLastLine(str) {
    let elem = document.querySelector(lastLineSelector);
    elem.innerHTML += str;
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
        appendToLastLine(key);
        apiSocket.send(`${key}\n`);
    }
}