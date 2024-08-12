const createRoomBtn = document.getElementById("createRoomBtn");
const roomNameInput = document.getElementById("roomNameInput");
const roomsList = document.getElementById("roomsList");
const messageInput = document.getElementById("messageInput");
const messagesList = document.getElementById('messagesList');
const leaveRoomBtn = document.getElementById("leaveRoomBtn");
const userData = JSON.parse(localStorage.getItem('userData'));
const userToken = localStorage.getItem('userToken');

let socket;
let currentRoomId = null;

function connectWebSocket(roomId) {
    const userToken = localStorage.getItem('userToken');
    socket = new WebSocket(`ws://localhost:8080/rooms/${roomId}/ws?token=${userToken}`);

    socket.onopen = function() {
        console.log('Connected to WebSocket');
    };

    socket.onmessage = function(event) {
        const message = JSON.parse(event.data);
        handleWebSocketMessage(message);
    };

    socket.onclose = function() {
        console.log('WebSocket connection closed');
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

function fetchRooms() {
    fetch('/rooms/', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${userToken}`
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            updateRoomsList(data.rooms);
        } else {
            console.error('Error fetching rooms:', data.message);
        }
    })
    .catch(error => console.error('Error fetching rooms:', error));
}

function updateRoomsList(rooms) {
    roomsList.innerHTML = ''; // Clear the existing list
    rooms.forEach(room => {
        const roomElement = document.createElement('div');
        roomElement.textContent = room.name;
        roomElement.addEventListener('click', () => joinRoom(room.id));
        roomsList.appendChild(roomElement);
    });
}

function createRoom() {
    const roomName = roomNameInput.value;
    if (roomName) {
        fetch('/rooms/new', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${userToken}`
            },
            body: JSON.stringify({ name: roomName })
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                roomNameInput.value = '';
                fetchRooms(); // Refresh the rooms list
            } else {
                console.error('Error creating room:', data.message);
            }
        })
        .catch(error => console.error('Error creating room:', error));
    }
}

function joinRoom(roomId) {
    fetch(`/rooms/join/${roomId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${userToken}`
        },
        body: JSON.stringify({ UserDTO: userData.user })
    })
    .then(response => {
        if (response.ok) {
            return response.json();
        } else {
            throw new Error('Error joining room: ' + response.statusText);
        }
    })
    .then(data => {
        if (data.success) {
            currentRoomId = roomId;
            window.location.href = `/statics/html/room.html?roomId=${roomId}`;
        } else {
            console.error('Error joining room:', data.message);
        }
    })
    .catch(error => console.error('Error joining room:', error));
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
    messagesList.appendChild(messageElement);
}

function leaveRoom() {
    if (socket) {
        socket.close();
        currentRoomId = null;
        window.location.href = '/statics/html/rooms.html';
    }
}

createRoomBtn.addEventListener("click", createRoom);
messageInput.addEventListener("keypress", (e) => {
    if (e.key === 'Enter') {
        sendMessage();
    }
});
leaveRoomBtn.addEventListener("click", leaveRoom);

fetchRooms();