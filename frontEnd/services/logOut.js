// logout.js
document.addEventListener('DOMContentLoaded', (event) => {
    const logoutButton = document.getElementById('logoutButton');

    if (logoutButton) {
        logoutButton.addEventListener('click', () => {
            // Implement your logout logic here
            // For example, removing JWT token from storage
            localStorage.removeItem('token');
            localStorage.removeItem('UserId');

            // Redirect to login page
            window.location.href = 'login.html';
        });
    }
});
