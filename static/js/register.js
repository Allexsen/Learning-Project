document.getElementById('registerForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const username = document.getElementById('username').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    const firstName = document.getElementById('firstName').value;
    const lastName = document.getElementById('lastName').value;

    if (password !== confirmPassword) {
        alert('Passwords do not match');
        return;
    }

    fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, email, password, firstName, lastName })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            localStorage.setItem('userToken', data.token);
            alert('Registration successful');
            window.location.href = '/statics/html/userProfile.html';
        } else {
            alert(`Registration failed: ${data.message}`);
        }
    })
    .catch(error => {
        alert('An error occurred while registering.');
    });
});
