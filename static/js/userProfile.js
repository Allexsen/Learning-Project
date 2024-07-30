document.addEventListener('DOMContentLoaded', function() {
    const searchForm = document.getElementById('searchForm');
    const addRecordButton = document.getElementById('addRecordButton');
    const addRecordModal = document.getElementById('addRecordModal');
    const addRecordForm = document.getElementById('addRecordForm');
    const feedbackDiv = document.getElementById('feedback');

    function loadUserProfile() {
        const userData = localStorage.getItem('userData');
        if (userData) {
            const parsedData = JSON.parse(userData);
            const user = parsedData.user;
            document.getElementById('userName').textContent = `Full Name: ${user.firstName} ${user.lastName}`;
            document.getElementById('userEmail').textContent = `Email: ${user.email}`;
            document.getElementById('userUsername').textContent = `Username: ${user.username}`;
            document.getElementById('totalTimeWorked').textContent = `Total Time Worked: ${user.total_hours} hours, ${user.total_minutes} minutes`;
            document.getElementById('logCount').textContent = `Log Count: ${user.log_count}`;
            updateWorkLogTable(user.worklog);
        }
    }

    searchForm.addEventListener('submit', function(event) {
        event.preventDefault();
        const cred = document.getElementById('cred').value;
        retrieveUserProfile(cred);
    });

    addRecordButton.addEventListener('click', function() {
        const userToken = localStorage.getItem('userToken');
        if (!userToken) {
            window.location.href = '/statics/html/login.html';
            return;
        }
        addRecordModal.style.display = 'block';
    });

    addRecordForm.addEventListener('submit', function(event) {
        event.preventDefault();
        const userData = localStorage.getItem('userData');
        if (userData) {
            const parsedData = JSON.parse(userData);
            const user = parsedData.user;

            const firstName = user.firstName;
            const lastName = user.lastName;
            const email = user.email;
            const hours = document.getElementById('hours').value;
            const minutes = document.getElementById('minutes').value;
            addRecord(firstName, lastName, email, hours, minutes);
            addRecordModal.style.display = 'none';
        }
    });

    document.querySelector('.close').addEventListener('click', function() {
        addRecordModal.style.display = 'none';
    });

    window.addEventListener('click', function(event) {
        if (event.target == addRecordModal) {
            addRecordModal.style.display = 'none';
        }
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
                localStorage.setItem('userData', JSON.stringify({ user: data.user }));
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

    function updateWorkLogTable(workLog) {
        const workLogTable = document.getElementById('workLogTable').getElementsByTagName('tbody')[0];
        workLogTable.innerHTML = ''; // Clear existing rows

        if (!workLog) return;

        workLog.forEach(entry => {
            const row = workLogTable.insertRow();
            row.classList.add('table-row');
            row.insertCell(0).textContent = entry.dateTime;
            row.insertCell(1).textContent = entry.hours;
            row.insertCell(2).textContent = entry.minutes;
            row.insertCell(3).textContent = entry.totalTime;
            const deleteCell = row.insertCell(4);
            deleteCell.innerHTML = '<span class="delete-button">Delete</span>';
            deleteCell.querySelector('.delete-button').addEventListener('click', function() {
                deleteRecord(entry.id);
            });
        });
    }

    function addRecord(firstName, lastName, email, hours, minutes) {
        const userToken = localStorage.getItem('userToken');
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
                loadUserProfile();
                showFeedback('Record added successfully!', 'success');
            } else {
                showFeedback('Failed to add record.', 'error');
            }
        })
        .catch(error => {
            showFeedback('An error occurred while adding the record.', 'error');
        });
    }

    function deleteRecord(id) {
        const userToken = localStorage.getItem('userToken');
        fetch('/record/delete', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${userToken}`
            },
            body: JSON.stringify({ id })
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                localStorage.setItem('userData', JSON.stringify({ user: data.user }));
                loadUserProfile();
                showFeedback('Record deleted successfully!', 'success');
            } else {
                showFeedback('Failed to delete record.', 'error');
            }
        })
        .catch(error => {
            showFeedback('An error occurred while deleting the record.', 'error');
        });
    }

    function showFeedback(message, type) {
        feedbackDiv.textContent = message;
        feedbackDiv.className = `feedback ${type}`;
        feedbackDiv.style.display = 'block';
        setTimeout(() => {
            feedbackDiv.style.display = 'none';
        }, 5000);
    }

    loadUserProfile();
});
