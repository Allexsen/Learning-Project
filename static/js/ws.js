document.addEventListener("DOMContentLoaded", () => {
    let socket;
    const connectBtn = document.getElementById("connectBtn");
    const disconnectBtn = document.getElementById("disconnectBtn");
    const sendBtn = document.getElementById("sendBtn");
    const connectionStatus = document.getElementById("connectionStatus");
    const messages = document.getElementById("messages");
    const messageInput = document.getElementById("messageInput");

    connectBtn.addEventListener("click", () => {
        const userData = localStorage.getItem('userData');
        let username = 'Guest';
        if (userData) {
            const parsedData = JSON.parse(userData);
            username = parsedData.user.username;
        }

        // Include the username as a query parameter in the WebSocket URL
        socket = new WebSocket(`ws://localhost:8080/ws?username=${encodeURIComponent(username)}`);

        socket.onopen = () => {
            connectionStatus.textContent = "Connected as " + username;
            connectBtn.disabled = true;
            disconnectBtn.disabled = false;
            sendBtn.disabled = false;
        };

        socket.onmessage = (event) => {
            const message = document.createElement("div");
            message.textContent = event.data;
            messages.appendChild(message);
        };

        socket.onclose = () => {
            connectionStatus.textContent = "Disconnected";
            connectBtn.disabled = false;
            disconnectBtn.disabled = true;
            sendBtn.disabled = true;
        };

        socket.onerror = (error) => {
            console.error("WebSocket Error: ", error);
        };
    });

    disconnectBtn.addEventListener("click", () => {
        if (socket) {
            socket.close();
        }
    });

    sendBtn.addEventListener("click", () => {
        if (socket && socket.readyState === WebSocket.OPEN) {
            const message = messageInput.value;
            socket.send(message);
            messageInput.value = "";
        }
    });
});