const API_BASE_URL_PUBLIC = window.config.API_BASE_URL_PUBLIC;

document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault(); 

    const email = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    let role=document.querySelector('input[name="role"]:checked').value;

    fetch(`${API_BASE_URL_PUBLIC}/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email, password,role })
    })
    .then(response => response.json())
    .then(data => {
        if (data.token) {
            localStorage.setItem('token', data.token); // Save the token
            localStorage.setItem('UserId',data.userId)
            if (role.toLowerCase() === 'admin') {
                window.location.href = '../html/properties.html'; // Redirect to admin page
            } else {
                window.location.href = '../html/booking.html'; // Redirect to user page
            }
        } else {
            alert('Invalid credentials'); // Show error message
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('An error occurred. Please try again.');
    });
});

