document.addEventListener('DOMContentLoaded', function() {
    const searchForm = document.getElementById('searchForm');
    const addRecordButton = document.getElementById('addRecordButton');
    const addRecordModal = document.getElementById('addRecordModal');
    const addRecordForm = document.getElementById('addRecordForm');
    const feedbackDiv = document.getElementById('feedback');
    
    function loadUserProfile() {
        const userData = localStorage.getItem('userData');
        if (userData) {
            const data = JSON.parse(userData);
            user = data.user
            document.getElementById('userName').textContent = `Full Name: ${user.name}`;
            document.getElementById('userEmail').textContent = `Email: ${user.email}`;
            document.getElementById('userUsername').textContent = `Username: ${user.username}`;
            document.getElementById('totalTimeWorked').textContent = `Total Time Worked: ${user.total_hours}`;
            document.getElementById('logCount').textContent = `Log Count: ${user.log_count}`;
            updateWorkLogTable(user.worklog);
        }
    }

    searchForm.addEventListener('submit', function(event) {
        event.preventDefault();
        const email = document.getElementById('email').value;
        retrieveUserProfile(email);
    });

    addRecordButton.addEventListener('click', function() {
        addRecordModal.style.display = 'block';
    });

    addRecordForm.addEventListener('submit', function(event) {
        event.preventDefault();
        const hours = document.getElementById('hours').value;
        const minutes = document.getElementById('minutes').value;
        addRecord(hours, minutes);
        addRecordModal.style.display = 'none';
    });

    document.querySelector('.close').addEventListener('click', function() {
        addRecordModal.style.display = 'none';
    });

    window.addEventListener('click', function(event) {
        if (event.target == addRecordModal) {
            addRecordModal.style.display = 'none';
        }
    });

    function retrieveUserProfile(email) {
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
                const user = data.user;
                document.getElementById('userName').textContent = `Full Name: ${user.name}`;
                document.getElementById('userEmail').textContent = `Email: ${user.email}`;
                document.getElementById('userUsername').textContent = `Username: ${user.username}`;
                document.getElementById('totalTimeWorked').textContent = `Total Time Worked: ${user.total_hours}`;
                document.getElementById('logCount').textContent = `Log Count: ${user.log_count}`;
                updateWorkLogTable(user.worklog);
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

    function addRecord(hours, minutes) {
        fetch('/record/add', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ hours, minutes })
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                updateWorkLogTable(data.workLog);
                showFeedback('Record added successfully!', 'success');
            } else {
                showFeedback('Failed to add record.', 'error');
            }
        })
        .catch(error => {
            showFeedback('An error occurred while adding the record.', 'error');
        });
    }

    function deleteRecord(recordId) {
        fetch('/record/delete', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ recordId })
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                updateWorkLogTable(data.workLog);
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

    loadUserProfile(); // Load user profile if data is available
});
