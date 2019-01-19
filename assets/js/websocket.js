// const eventId
const EventUpdate = 0;

// create web socket
var WS_URL = "ws://" + location.hostname + ":8888/ws";
var ws = new WebSocket(WS_URL);

// connection open
ws.addEventListener('open', function () {
    ws.send('hello server');
});

// listen message
ws.addEventListener('message', function (ev) {
    var s = JSON.parse(ev.data);
    console.log(s);
    $(".users").html("");
    for(var i = 0; s.length; i++) {
        var path = window.location.protocol+"//"+window.location.host+"/"+s[i]["IconPath"];
        $(".users").append('<div class="user-box"><div class="bg-image" style="background-image: url('+path+')"></div><span class="image-title">'+s[i]["UserName"]+'</span></div>');
    }
});

// connection close
ws.addEventListener('close', function (ev) {
    console.log('connection close', ev.code);
});