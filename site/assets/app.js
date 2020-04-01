var appendNext = false;

// Setup message handling
function writeOutput(msg) {
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
        elem.innerHTML += `<div class="line">${line}</div>`;
        elem.scrollTop = elem.scrollHeight;
    }
}

function connect() {
    let ws = new WebSocket("wss://deathtax.kayotic.io/api");
    ws.onmessage = writeOutput;
    return ws;
}

// Connect to the api
let apiSocket = connect();

// Listen for input keypresses
const InputKeys = ['Y', 'N', 'S', '?'];

document.addEventListener('keyup', function (event) {
    if (event.defaultPrevented) {
        return;
    }

    var key = event.key || event.keyCode;
    if (InputKeys.includes(key.toUpperCase())) {
        appendNext = true;
        apiSocket.send(`${key}\n`);
    }
});