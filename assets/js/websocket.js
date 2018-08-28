// create web socket
var WS_URL = "ws://" + location.hostname + ":8888/ws";
var ws = new WebSocket(WS_URL);

// connection open
ws.addEventListener('open', function () {
    ws.send('hello server');
});

// listen message
ws.addEventListener('message', function (ev) {
    console.log('message from server', ev.data);
});

// connection close
ws.addEventListener('close', function (ev) {
    console.log('connection close', ev.code)
})