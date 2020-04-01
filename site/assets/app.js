var appendNext = false;
var apiSocket;
var inputKeys = [];
const inputKeyPattern = /\[(\S)\]\s\S*/g;
const inputLinePattern = /\(default is .*\):$/;


document.addEventListener('DOMContentLoaded', function() {
    // Connect to the api
    apiSocket = connect();
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
        data = "&nbsp;";
    }

    let elem = document.querySelector("#output");

    if(data == "\n") {
        data = "&nbsp;"
    }

    for(line of data.split("\n")) {
        if(line.match(inputLinePattern)) {
            inputKeys = Array.from(line.matchAll(inputKeyPattern)).map(match => match[1])
            enableInput()
        }

        elem.innerHTML += `<div class="line">${line}</div>`;
        elem.scrollTop = elem.scrollHeight;
    }
}

function enableInput() {
    document.addEventListener('keyup', handeKeyEvent);
}

function disableInput() {
    document.removeEventListener('keyup', handeKeyEvent);
}

function handeKeyEvent(event) {
    if (event.defaultPrevented) {
        return;
    }
    
    let allowedKeys = inputKeys;
    disableInput()
    console.log(allowedKeys);


    var key = event.key || event.keyCode;
    if (inputKeys.includes(key.toUpperCase())) {
        appendNext = true;
        apiSocket.send(`${key}\n`);
    }

    inputKeys = [];
    console.log(allowedKeys);
}