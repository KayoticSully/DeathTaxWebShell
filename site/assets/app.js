'use strict';

var apiSocket;
var inputKeys = [];
const inputKeyPattern = /\[(\S)\]\s\S*/g;
const inputLinePattern = /\(default is \"\S\"\):\s*$/;
const lastLineSelector = "#output .line:last-child";

document.addEventListener('DOMContentLoaded', function() {
    // Connect to the api
    apiSocket = connect();
    document.addEventListener('keyup', handleKeyEvent);
}, false);

function connect() {
    let ws = new WebSocket("wss://deathtax.kayotic.io/api");
    ws.onmessage = handleMessage;
    ws.onclose = function(msg) {
      console.log(msg);
    };
    return ws;
}

// Setup message handling
function handleMessage(msg) {
    console.log(msg);
    let data = msg.data;

    // A single newline will create two newlines on 
    // a split("\n"). Convert it to a blank string to
    // replicate a single new line output.
    if(data == "\n") {
        data = "";
    }

    for(let line of data.split("\n")) {
        if(line.match(inputLinePattern)) {
            let inputMatches = Array.from(line.matchAll(inputKeyPattern));

            for(const match of inputMatches) {
                inputKeys.push(match[1]);
                let inputButton = `<span class="key-input-button" onclick="clickKeyInput('${match[1]}')">${match[0]}</span>`;
                line = line.replace(match[0], inputButton);
            }

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

function clickKeyInput(key) {
    handleKeyEvent({
        data: key,
    });
}

function handleKeyEvent(event) {
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
    } else {
        // re-enable input if key was not valid input
        inputKeys = allowedKeys;
    }
}