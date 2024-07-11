document.addEventListener('DOMContentLoaded', function () {
    const emailField = document.getElementById('email');
    const usernameField = document.getElementById('username');
    const emailFeedback = document.getElementById('emailFeedback');
    const usernameFeedback = document.getElementById('usernameFeedback');
    
    let isEmailAvailable = false;
    let isUsernameAvailable = false;

    emailField.addEventListener('blur', function () {
        const email = emailField.value.trim();
        if (email === '') {
            emailFeedback.textContent = 'Email cannot be empty';
            emailFeedback.style.color = 'red';
            isEmailAvailable = false;
            return;
        }

        if (!isValidEmail(email)) {
            emailFeedback.textContent = 'Invalid email format';
            emailFeedback.style.color = 'red';
            isEmailAvailable = false;
            return;
        }

        fetch('/user/check-email', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email })
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                isEmailAvailable = true;
                emailFeedback.textContent = 'Email is available';
                emailFeedback.style.color = 'green';
            } else {
                isEmailAvailable = false;
                emailFeedback.textContent = 'Email is already taken';
                emailFeedback.style.color = 'red';
            }
        })
        .catch(error => {
            isEmailAvailable = false;
            console.error('Error checking email:', error);
            emailFeedback.textContent = 'Error checking email';
            emailFeedback.style.color = 'red';
        });
    });

    usernameField.addEventListener('blur', function () {
        const username = usernameField.value.trim();
        if (username === '') {
            usernameFeedback.textContent = 'Username cannot be empty';
            usernameFeedback.style.color = 'red';
            isUsernameAvailable = false;
            return;
        }

        if (!isValidUsername(username)) {
            usernameFeedback.textContent = 'Username can only contain letters, numbers, dashes, and underscores';
            usernameFeedback.style.color = 'red';
            isUsernameAvailable = false;
            return;
        }
        
        if (username.length <= 3) {
            usernameFeedback.textContent = 'Username must be longer than 3 characters';
            usernameFeedback.style.color = 'red';
            isUsernameAvailable = false;
            return;
        }

        fetch('/user/check-username', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username })
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                isUsernameAvailable = true;
                usernameFeedback.textContent = 'Username is available';
                usernameFeedback.style.color = 'green';
            } else {
                isUsernameAvailable = false;
                usernameFeedback.textContent = 'Username is already taken';
                usernameFeedback.style.color = 'red';
            }
        })
        .catch(error => {
            isUsernameAvailable = false;
            console.error('Error checking username:', error);
            usernameFeedback.textContent = 'Error checking username';
            usernameFeedback.style.color = 'red';
        });
    });

    document.getElementById('registerForm').addEventListener('submit', function (event) {
        event.preventDefault();
        
        const username = usernameField.value.trim();
        const email = emailField.value.trim();
        
        if (username === '' || email === '') {
            alert('Please fill in both email and username fields.');
            return;
        }

        if (!isEmailAvailable || !isUsernameAvailable) {
            alert('Please make sure that both email and username are available.');
            return;
        }

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

    function isValidEmail(email) {
        const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailPattern.test(email);
    }

    function isValidUsername(username) {
        const usernamePattern = /^[a-zA-Z0-9-_]+$/;
        return usernamePattern.test(username);
    }
});
