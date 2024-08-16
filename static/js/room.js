const participantsList = document.getElementById('participantsList');
const messageInput = document.getElementById("messageInput");
const messagesList = document.getElementById("messagesList");
const sendBtn = document.getElementById("sendBtn");
const leaveRoomBtn = document.getElementById("leaveRoomBtn");
const userToken = localStorage.getItem('userToken');
const roomId = new URLSearchParams(window.location.search).get('roomId');
const userData = JSON.parse(localStorage.getItem('userData'));

let socket;

document.addEventListener("DOMContentLoaded", () => {
    fetch(`/rooms/${roomId}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${userToken}`
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            document.title = `${data.room.name} - Rooms`;
            document.getElementById('roomTitle').innerText = data.room.name;
            fetchParticipants(roomId, userToken);
        } else {
            console.error('Error fetching room details:', data.message);
        }
    })
    .catch(error => console.error('Error fetching room details:', error));
});

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
        case 'participantJoined':
            addParticipant(message);
            break;
        case 'participantLeft':
            removeParticipant(message);
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
    const timestamp = new Date(message.timestamp).toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' });
    messageElement.title = `${timestamp}`;

    if (message.sender_id === 0) {
        messageElement.classList.add('system-message');
    } else if (message.sender_id === userData.user.id) {
        messageElement.classList.add('my-message');
    } else {
        messageElement.classList.add('other-message');
    }

    messagesList.appendChild(messageElement);
}

function addParticipant(message) {
    const participantElement = document.createElement('div');
    participantElement.textContent = message.content;
    participantElement.id = `participant-${message.sender_id}`;
    participantsList.appendChild(participantElement);

    // Transform into a system message
    message.type = 'chatMessage';
    message.content = `${message.content} has joined the room`,
    message.sender_id = 0
    displayMessage(message);
}

function removeParticipant(message) {
    const participantElement = document.getElementById(`participant-${message.sender_id}`);
    if (participantElement) {
        participantsList.removeChild(participantElement);
    }

    // Transform into a system message
    message.type = 'chatMessage';
    message.content = `${message.content} has left the room`,
    message.sender_id = 0
    displayMessage(message);
}

function fetchParticipants(roomId, userToken) {
    fetch(`/rooms/${roomId}/participants`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${userToken}`
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            participantsList.innerHTML = '';
            data.participants.forEach(participant => {
                const participantElement = document.createElement('div');
                participantElement.innerText = participant.username;
                participantsList.appendChild(participantElement);
            });
        } else {
            console.error('Error fetching participants:', data.message);
        }
    })
    .catch(error => console.error('Error fetching participants:', error));
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