const sendBtn = document.getElementById("sendBtn");
const messageInput = document.getElementById("messageInput");
const messagesDiv = document.getElementById("messages");
const connectionStatus = document.getElementById("connectionStatus");
const userData = JSON.parse(localStorage.getItem('userData'));

let socket;

connectBtn.addEventListener("click", () => {
    const userToken = localStorage.getItem('userToken');

    if (!userToken) {
        alert('User is not logged in.');
        return;
    }

    const wsUrl = `ws://localhost:8080/ws?token=${userToken}`;
    socket = new WebSocket(wsUrl);
    
    socket.onopen = () => {
        connectionStatus.textContent = "Connected";
        connectBtn.disabled = true;
        disconnectBtn.disabled = false;
        sendBtn.disabled = false;
    };
    
    socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        const messageElement = document.createElement('div');
        messageElement.textContent = `${message.sender_id}: ${message.content}`;
        messagesDiv.appendChild(messageElement);
    };

    socket.onclose = () => {
        connectionStatus.textContent = "Disconnected";
        connectBtn.disabled = false;
        disconnectBtn.disabled = true;
        sendBtn.disabled = true;
    };
});

disconnectBtn.addEventListener("click", () => {
    if (socket) {
        socket.close();
    }
});

sendBtn.addEventListener("click", () => {
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
});