
let socket, connectButton, screenImg

function connect() {

    protocol = window.location.protocol == "http:" ? "ws:" : "wss:"
    ws_uri = `${protocol}//${window.location.host}/ws/`

    socket = new WebSocket(ws_uri);

    socket.addEventListener('open', function (event) {
        socket.send(window.location.pathname)
    })

    socket.addEventListener('message', function (event) {
        screenImg.src = event.data
    })

}

function disconnect() {
    if(socket) {
        socket.close();
    }
}

(function(){

    connectButton = document.getElementById("connect_button")
    screenImg = document.getElementById("screen_img")

    connectButton.onclick = function (event) {

        if (connectButton.innerHTML == "Connect") {
            connectButton.innerHTML = "Disconect"
            connect()
        } else {
            connectButton.innerHTML = "Connect"
            disconnect()
        }
    }

})();