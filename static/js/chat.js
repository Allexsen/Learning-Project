const sendBtn = document.getElementById("sendBtn");
const messageInput = document.getElementById("messageInput");
const messagesDiv = document.getElementById("messages");
const participantsList = document.getElementById("participantsList");
const userData = JSON.parse(localStorage.getItem('userData'));
const urlParams = new URLSearchParams(window.location.search);
const roomId = urlParams.get('roomId');

let socket;

function connectWebSocket() {
    const userToken = localStorage.getItem('userToken');

    if (!userToken) {
        alert('User is not logged in.');
        return;
    }

    const wsUrl = `ws://localhost:8080/ws/rooms/${roomId}?token=${userToken}`
    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
        console.log("Connected to WebSocket");
        sendBtn.disabled = false;
    };

    socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        handleWebSocketMessage(message);
    };

    socket.onclose = () => {
        console.log("Disconnected from WebSocket");
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

function displayMessage(message) {
    const messageElement = document.createElement('div');
    messageElement.textContent = `${message.sender_id}: ${message.content}`;
    messagesDiv.appendChild(messageElement);
}

function updateParticipantsList(participants) {
    participantsList.innerHTML = '';
    participants.forEach(participant => {
        const participantElement = document.createElement('div');
        participantElement.textContent = participant.name;
        participantsList.appendChild(participantElement);
    });
}

function sendMessage() {
    const messageContent = messageInput.value;

    if (messageContent && socket && socket.readyState === WebSocket.OPEN) {
        const message = {
            id: Date.now(), // Placeholder message ID
            sender_id: userData.userId,
            content: messageContent,
            timestamp: Date.now(),
            status: "sent" // TBD
        };
        socket.send(JSON.stringify(message));
        messageInput.value = "";
    }
}

sendBtn.addEventListener("click", sendMessage);

connectWebSocket();