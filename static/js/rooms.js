const createRoomBtn = document.getElementById("createRoomBtn");
const roomNameInput = document.getElementById("roomNameInput");
const roomsList = document.getElementById("roomsList");
const userData = JSON.parse(localStorage.getItem('userData'));
const userToken = localStorage.getItem('userToken');

let socket;
function connectWebSocket() {
    // WebSocket connection logic
}

function handleWebSocketMessage(message) {
    // Handle incoming WebSocket messages
}

function fetchRooms() {
    fetch('/rooms/', {
        headers: {
            'Authorization': `Bearer ${userToken}`
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            console.log('Fetched rooms:', data.rooms); // Log the fetched data
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

function removeRoomFromList(roomId) {
    const roomElement = document.querySelector(`.room-card[data-room-id="${roomId}"]`);
    if (roomElement) {
        roomElement.remove();
    }
}

function createRoom() {
    const roomName = roomNameInput.value;
    if (roomName) {
        fetch('/rooms/new', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
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
            'Authorization': `Bearer ${userToken}`
        },
        body: JSON.stringify({ UserDTO: userData.user})
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
            window.location.href = `/statics/html/chat.html?roomId=${roomId}`;
        } else {
            console.error('Error joining room:', data.message);
        }
    })
    .catch(error => console.error('Error joining room:', error));
}

createRoomBtn.addEventListener("click", createRoom);

connectWebSocket();
fetchRooms();