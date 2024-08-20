document.addEventListener('DOMContentLoaded', function() {
    const searchForm = document.getElementById('searchForm');

    function loadUserProfile() {
        const userData = localStorage.getItem('userData');
        if (userData) {
            const parsedData = JSON.parse(userData);
            const user = parsedData.user;
            document.getElementById('userName').textContent = `Full Name: ${user.firstName} ${user.lastName}`;
            document.getElementById('userEmail').textContent = `Email: ${user.email}`;
            document.getElementById('userUsername').textContent = `Username: ${user.username}`;
        }
    }

    searchForm.addEventListener('submit', function(event) {
        event.preventDefault();
        const cred = document.getElementById('cred').value;
        retrieveUserProfile(cred);
    });

    function retrieveUserProfile(cred) {
        fetch('/user/retrieve', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ cred })
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                const user = data.user;
                localStorage.setItem('userData', JSON.stringify({ user: user }));
                loadUserProfile();
                showFeedback('User profile retrieved successfully!', 'success');
            } else {
                showFeedback('Failed to retrieve user profile.', 'error');
            }
        })
        .catch(error => {
            showFeedback('An error occurred while retrieving user profile.', 'error');
        });
    }

    loadUserProfile();
});
