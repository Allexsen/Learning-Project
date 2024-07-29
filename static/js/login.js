document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const cred = document.getElementById('cred').value;
    const password = document.getElementById('password').value;
    const urlParams = new URLSearchParams(window.location.search);
    const redirectUrl = urlParams.get('redirect') || `/`;

    fetch('/user/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ cred, password })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            localStorage.setItem('userToken', data.token);
            return fetch('/user/retrieve', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${data.token}`
                },
                body: JSON.stringify({ cred: cred })
            });
        } else {
            alert('Invalid credentials');
            throw new Error('Invalid credentials');
        }
    })
    .then(response => response.json())
    .then(profileData => {
        if (profileData.success) {
            localStorage.setItem('userData', JSON.stringify({ user: profileData.user }));
            window.location.href = redirectUrl;
        } else {
            alert('Failed to retrieve user profile.');
        }
    })
    .catch(error => {
        alert('An error occurred while logging in.');
        console.error(error);
    });
});
