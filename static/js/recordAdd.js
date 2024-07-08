document.getElementById('addRecordForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const userToken = localStorage.getItem('userToken');
    if (!userToken) {
        window.location.href = '/statics/html/login.html';
        return;
    }

    const firstName = document.getElementById('firstname').value;
    const lastName = document.getElementById('lastname').value;
    const email = document.getElementById('email').value;
    const hours = document.getElementById('hours').value;
    const minutes = document.getElementById('minutes').value;

    fetch('/record/add', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${userToken}`
        },
        body: JSON.stringify({ firstName, lastName, email, hours, minutes })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            localStorage.setItem('userData', JSON.stringify({ user: data.user }));
            window.location.href = '/statics/html/userProfile.html';
        } else {
            alert('Failed to add record.');
        }
    })
    .catch(error => {
        alert('An error occurred while adding the record.');
    });
});
