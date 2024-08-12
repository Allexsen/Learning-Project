const messageInput = document.getElementById("messageInput");
const messagesList = document.getElementById("messagesList");
const sendBtn = document.getElementById("sendBtn");
const leaveRoomBtn = document.getElementById("leaveRoomBtn");
const userToken = localStorage.getItem('userToken');
const roomId = new URLSearchParams(window.location.search).get('roomId');

let socket;

function connectWebSocket(roomId) {
    socket = new WebSocket(`ws://localhost:8080/rooms/${roomId}/ws?token=${userToken}`);

    socket.onopen = function() {
        console.log('Connected to WebSocket');
        sendBtn.disabled = false;
    };

    socket.onmessage = function(event) {
        const message = JSON.parse(event.data);
        handleWebSocketMessage(message);
    };

    socket.onclose = function() {
        console.log('WebSocket connection closed');
        sendBtn.disabled = true;
    };
}

function handleWebSocketMessage(message) {
    switch (message.type) {
        case 'chatMessage':
            displayMessage(message);
            break;
        case 'participantsList':
            updateParticipantsList(message.participants);
            break;
        case 'roomDeleted':
            alert('Room has been deleted');
            window.location.href = '/statics/html/rooms.html';
            break;
        // Handle other message types as needed
    }
}

function sendMessage() {
    const message = messageInput.value;
    if (message && socket) {
        socket.send(JSON.stringify({ type: 'chatMessage', content: message }));
        messageInput.value = '';
    }
}

function displayMessage(message) {
    const messageElement = document.createElement('div');
    messageElement.textContent = message.content;
    messageElement.title = `User: ${message.user.name}, ID: ${message.user.id}`;

    if (message.id === 0) {
        messageElement.classList.add('system-message');
    } else if (message.user.id === userToken) {
        messageElement.classList.add('my-message');
    } else {
        messageElement.classList.add('other-message');
    }

    messagesList.appendChild(messageElement);
}

function leaveRoom() {
    if (socket) {
        socket.close();
    }
    window.location.href = '/statics/html/rooms.html';
}

sendBtn.addEventListener("click", sendMessage);
messageInput.addEventListener("keypress", (e) => {
    if (e.key === 'Enter') {
        sendMessage();
    }
});
leaveRoomBtn.addEventListener("click", leaveRoom);

connectWebSocket(roomId);