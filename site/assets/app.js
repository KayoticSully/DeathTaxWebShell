
// Setup message handling
function displayOutput(msg) {
    let elem = document.querySelector("#output");

    let data = msg.data;
    if(data == "\n") {
        data = "<br />"
    }

    elem.innerHTML += `<div class="line">${data}</div>`;
    elem.scrollTop = elem.scrollHeight;
}

function connect() {
    let ws = new WebSocket("wss://deathtax.kayotic.io/api");
    ws.onmessage = displayOutput;
    return ws;
}

// Connect to the api
let apiSocket = connect();

// Listen for input keypresses
const InputKeys = ['Y', 'N', 'S'];

document.addEventListener('keyup', function (event) {
    if (event.defaultPrevented) {
        return;
    }

    var key = event.key || event.keyCode;
    if (InputKeys.includes(key.toUpperCase())) {
        apiSocket.send(`${key}\n`);
    }
});