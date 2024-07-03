const API_BASE_URL_PUBLIC = window.config.API_BASE_URL_PUBLIC;

document.addEventListener('DOMContentLoaded', () => {
    // Assuming the user is logged in and the token is stored in localStorage
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = 'login.html';
        alert('Please Login');
        console.log('User not logged in'); 
    }
});

document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('RegisterProperty');
    
    form.addEventListener('submit', (event) => {
        event.preventDefault(); // Prevent the form from submitting the default way

        const formData = new FormData(form);

        const propertyData = {
            userId: localStorage.getItem('UserId'), 
            title: formData.get('title'),
            description: formData.get('description'),
            streetAddr: formData.get('streetAddr'),
            city: formData.get('city'),
            country: formData.get('country'),
            zipCode: formData.get('zipCode'),
            bedRooms: formData.get('bedRooms'),
            bathRooms: formData.get('bathRooms'),
            accomodates: formData.get('accomodates'),
            currency: formData.get('currency'),
            price: formData.get('price'),
            minStay: formData.get('minStay'),
            maxStay: formData.get('maxStay'),
            propertyType: {
                privateBed: formData.get('propertyType') === 'privateBed',
                whole: formData.get('propertyType') === 'whole',
                shared: formData.get('propertyType') === 'shared'
            },
            amenities: {
                ac: formData.get('amenities_ac') === 'on',
                heater: formData.get('amenities_heater') === 'on',
                tv: formData.get('amenities_tv') === 'on',
                wifi: formData.get('amenities_wifi') === 'on'
            },
            spaces: {
                closets: formData.get('spaces_closets') === 'on',
                gym: formData.get('spaces_gym') === 'on',
                kitchen: formData.get('spaces_kitchen') === 'on',
                parking: formData.get('spaces_parking') === 'on',
                pool: formData.get('spaces_pool') === 'on'
            }
        };

        const imageUpload = formData.get('imageUpload');
        if (imageUpload && imageUpload.size > 0) {
            const reader = new FileReader();
            reader.readAsDataURL(imageUpload);
            reader.onload = () => {
                propertyData.image = reader.result.split(',')[1]; // Extract base64 string after comma
                submitPropertyData(propertyData);
            };
        } else {
            submitPropertyData(propertyData);
        }
    });
});

function submitPropertyData(propertyData) {
    fetch('http://localhost:8081/Property', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(propertyData)
    })
    .then(response => response.json())
    .then(data => {
        console.log('Success:', data);
        alert('Property registered successfully!');
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Failed to register property.');
    });
}

function showUserProperties(){
    window.location.href = '../html/myProperties.html';
}
