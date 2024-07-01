document.getElementById('retrieveForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const email = document.getElementById('email').value;
    fetch('/user/retrieve', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            localStorage.setItem('userData', JSON.stringify({ user: data.user }));
            window.location.href = '/statics/html/userProfile.html';
        } else {
            alert('Failed to retrieve user data.');
        }
    })
    .catch(error => {
        alert('An error occurred while retrieving user data.');
    });
});