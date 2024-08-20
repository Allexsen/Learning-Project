document.getElementById('retrieveForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const cred = document.getElementById('cred').value;
    retrieveUserProfile(cred);
});

function retrieveUserProfile(cred) {
    const userToken = localStorage.getItem('userToken');
    if (!userToken) {
        window.location.href = '/statics/html/login.html';
        return;
    }

    fetch(`/user/profile?cred=${encodeURIComponent(cred)}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${userToken}`
        },
        body: JSON.stringify({ cred: cred })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            localStorage.setItem('userData', JSON.stringify({ user: data.user }));
            window.location.href = '/statics/html/userProfile.html';
        } else {
            alert('Failed to retrieve user profile.');
        }
    })
    .catch(error => {
        alert('An error occurred while retrieving the user profile.');
    });
}