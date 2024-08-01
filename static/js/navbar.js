document.addEventListener('DOMContentLoaded', function() {
    const navbarRight = document.querySelector('.navbar-right');
    const userToken = localStorage.getItem('userToken');

    if (userToken) {
        navbarRight.innerHTML = `
            <a href="/statics/html/userProfile.html">Profile</a>
            <a href="/statics/html/ws.html">WebSocket</a>
            <a href="/statics/html/rooms.html">Rooms</a>
            <a href="#" id="logout">Logout</a>
        `;

        document.getElementById('logout').addEventListener('click', function() {
            localStorage.removeItem('userToken');
            localStorage.removeItem('userData');
            window.location.href = '/statics/html/login.html';
        });
    } else {
        const currentPath = encodeURIComponent(window.location.pathname);
        navbarRight.innerHTML = `
            <a href="/statics/html/login.html?redirect=${currentPath}">Login</a>
            <a href="/statics/html/register.html?redirect=${currentPath}">Register</a>
            <a href="/statics/html/rooms.html">Rooms</a>
        `;
    }
});