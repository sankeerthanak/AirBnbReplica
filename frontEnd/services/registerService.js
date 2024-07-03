const API_BASE_URL_PUBLIC = window.config.API_BASE_URL_PUBLIC;

document.getElementById('SignUpForm').addEventListener('submit', function(event) {
    event.preventDefault(); 

    const email = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const cnfrmPassword = document.getElementById('cnfrmpassword').value;
    const firstName = document.getElementById('firstName').value;
    const lastName = document.getElementById('lastName').value;
    const roles = document.querySelectorAll('input[name="role"]:checked');
    const roleName = [];

    roles.forEach((role) => {
        roleName.push(role.value);
    });

    //console.log(selectedRoles); // This will log an array of selected roles


    if(password===cnfrmPassword){
        fetch(`${API_BASE_URL_PUBLIC}/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password, firstName, lastName, roleName })
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                alert('Trouble creating the user'); 
            } else {
                alert('User successfully created'); // Show error message
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('An error occurred. Please try again.');
        });

    }else {
            alert('Passwords does not match'); // Show error message
    }
});