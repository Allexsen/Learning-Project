document.getElementById('addRecordForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const name = document.getElementById('name').value;
    const email = document.getElementById('email').value;
    const hours = document.getElementById('hours').value;
    const minutes = document.getElementById('minutes').value;
    fetch('/record/add', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name, email, hours, minutes })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            localStorage.setItem('userData', JSON.stringify(data));
            window.location.href = '/statics/html/userProfile.html';
        } else {
            alert('Failed to add record.');
        }
    })
    .catch(error => {
        alert('An error occurred while adding the record.');
    });
});