

var conn = null;
function sendMessage() {
    var message = document.getElementById("message").value;
    var username = document.getElementById("username").value;
    var messageObject = {
        "message": message,
        "username": username
    };
    var messageJSON = JSON.stringify(messageObject);
    
    if (conn != null) {
        conn.send(messageJSON);
    }
}

window.onload = function () {
    connect();
};

function connect() {
    if (window["WebSocket"]) {
        // Send message
        conn = new WebSocket("ws://" + document.location.host + "/cr/main");
        conn.onopen = function (evt) {
            var chat = document.getElementById("chat");
            var username = document.getElementById("username").value;
            if (username == "") {
                username = "anon";
            }
            chat.innerHTML += "<br/>Connection established";
            chat.send({"login": username});
        }

        conn.onmessage = function (evt) {
            var chat = document.getElementById("chat");
            var message = evt.data;
            var messageObject = JSON.parse(message);
            chat.innerHTML += "<br/>" + messageObject.username + ": " + messageObject.message;
        }
    }
}
