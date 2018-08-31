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
    console.log('event id is', ev.data);
    if (ev.data == EventUpdate) {
        fetch(window.location.protocol+'//'+window.location.host+'/user/viewAll')
            .then(function (response) {
                if(response.ok) {
                    return response.json();
                }

                throw new Error('Network response was not ok.')
            })
            .then(function (json) {
                console.log(json);
                for(var i = 0; json.length; i++) {
                    var path = window.location.protocol+"//"+window.location.host+"/"+json[i]["IconPath"];
                    $(".users").append('<div class="user-box"><div class="bg-image" style="background-image: url('+path+')"></div><span class="image-title">'+json[i]["UserName"]+'</span></div>');
                }
            })
    } else {
        console.log("event failed");
    }
});

// connection close
ws.addEventListener('close', function (ev) {
    console.log('connection close', ev.code)
});