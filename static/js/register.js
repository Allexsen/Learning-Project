document.getElementById('registerForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const username = document.getElementById('username').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    const firstName = document.getElementById('firstName').value;
    const lastName = document.getElementById('lastName').value;
    const urlParams = new URLSearchParams(window.location.search);
    const redirectUrl = urlParams.get('redirect') || `/`;

    if (password !== confirmPassword) {
        alert('Passwords do not match');
        return;
    }

    const registrationData = {
        username,
        email,
        password,
        firstName,
        lastName
    };

    fetch('/user/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(registrationData)
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
                body: JSON.stringify({ email: email })
            });
        } else {
            alert(`Registration failed: ${data.message}`);
            throw new Error('Registration failed');
        }
    })
    .then(response => response.json())
    .then(profileData => {
        if (profileData.success) {
            localStorage.setItem('userData', JSON.stringify({ user: profileData.user }));
            alert('Registration successful');
            window.location.href = redirectUrl;
        } else {
            alert('Failed to retrieve user profile.');
        }
    })
    .catch(error => {
        alert('An error occurred while registering.');
        console.error(error);
    });
});
