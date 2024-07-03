const API_BASE_URL_PUBLIC = window.config.API_BASE_URL_PUBLIC;

document.addEventListener('DOMContentLoaded', () => {
    // Assuming the user is logged in and the token is stored in localStorage
    const token = localStorage.getItem('token');
    if (token) {
        fetchProperties();
    } else {
        window.location.href = 'login.html';
        alert('Please Login');
        console.log('User not logged in');
    }
});

let properties=null;

function fetchProperties() {
    fetch(`${API_BASE_URL_PUBLIC}/Property/`+localStorage.getItem('UserId'), {
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
    })
    .then(response => response.json())
    .then(data => {
        localStorage.setItem('properties', JSON.stringify(data));
        renderProperties(data);
    })
    .catch(error => {
        console.error('Error fetching properties:', error);
    });
}

const propertiesContainer = document.getElementById('propertiesContainer');
const propertyModal = document.getElementById('propertyModal');
const modalContent = document.getElementById('modalContent');
const bookingModal = document.getElementById('bookingModal');
let currentProperty = null;

function renderProperties(propertiesData) {
    propertiesData.forEach(property => {
        const card = document.createElement('div');
        card.classList.add('card');
        card.dataset.propertyId = property._id;

        const img = document.createElement('img');
        img.src = `data:image/png;base64, ${property.image}`;
        img.alt = property.title;

        const title = document.createElement('div');
        title.classList.add('card-title');
        title.textContent = property.title;

        const price = document.createElement('div');
        price.classList.add('card-price');
        price.textContent = property.price;

        const buttonContainer = document.createElement('div');
        buttonContainer.classList.add('button-container');

        const button = document.createElement('button');
        button.classList.add('button');
        button.textContent = 'Edit Property';
        button.onclick = (event) => {
            event.stopPropagation();
            currentProperty=property
            showBookingForm();
        };

        buttonContainer.appendChild(button);

        card.appendChild(img);
        card.appendChild(title);
        card.appendChild(price);
        card.appendChild(buttonContainer);

        card.onclick = () => showPropertyDetails(property);

        propertiesContainer.appendChild(card);
    });
}

function showBookingForm() {
    propertyModal.classList.remove('active');
    bookingModal.classList.add('active');
}

function showPropertyDetails(property) {
    currentProperty=property
    modalContent.innerHTML = `
        <h2>${property.title}</h2>
        <p>${property.description}</p>
        <p><strong>Accommodates:</strong> ${property.accomodates}</p>
        <p><strong>Price:</strong> ${property.price}</p>
        <p><strong>Address:</strong> ${property.streetAddr}, ${property.city}, ${property.country}</p>
        <p><strong>Amenities:</strong> AC: ${property.amenities.ac}, Heater: ${property.amenities.heater}, TV: ${property.amenities.tv}, WiFi: ${property.amenities.wifi}</p>
        <p><strong>Spaces:</strong> Closets: ${property.spaces.closets}, Gym: ${property.spaces.gym}, Kitchen: ${property.spaces.kitchen}, Parking: ${property.spaces.parking}, Pool: ${property.spaces.pool}</p>
    `;

    const bookNowButton = document.createElement('button');
    bookNowButton.classList.add('button');
    bookNowButton.textContent = 'Edit Property';
    bookNowButton.onclick = (event) => {
        event.stopPropagation();
        showBookingForm();
    };

    modalContent.appendChild(bookNowButton);

    propertyModal.classList.add('active');
}

function closeBookingModal() {
    bookingModal.classList.remove('active');
}

function propertyData() {
    const form = document.getElementById('RegisterProperty');
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
}

function submitPropertyData(propertyData) {
    fetch(`${API_BASE_URL_PUBLIC}/Property/`+currentProperty.propertyId, {
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
        fetchProperties();
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Failed to register property.');
    });
}

function closeModal() {
    propertyModal.classList.remove('active');
}
// function showUserBookings(){
//     window.location.href = '../html/myBookings.html';
// }

// function confirmBooking(bookingDetails) {

//     const confirmation = confirm(`Do you want to book the property "?`);
//     if (confirmation) {
//         fetch('http://localhost:8081/Booking', {
//             method: 'POST',
//             headers: {
//                 'Authorization': `Bearer ${localStorage.getItem('token')}`,
//                 'Content-Type': 'application/json'
//             },
//             body: JSON.stringify({ bookingDetails })
//         })
//         .then(response => response.json())
//         .then(data => {
//             renderProperties(data);
//         })
//         .catch(error => {
//             console.error('Error fetching properties:', error);
//         });
//         alert('Booking confirmed!');
//         closeBookingModal();
//     }
// }

// export function showUserBookings(){

//     window.location.href = '../html/myBookings.html';
//     // fetch('http://localhost:8081/Booking/66751fcdcbb32f0b46e7e214',
//     //     //+localStorage.getItem('userId')+'}', 
//     //     {
//     //     headers: {
//     //         'Authorization': `Bearer ${localStorage.getItem('token')}`
//     //     }
//     // })
//     // .then(response => response.json())
//     // .then(data => {
//     //     renderProperties(data);
//     //     window.location.href = '/myBookings.html';
//     // })
//     // .catch(error => {
//     //     console.error('Error fetching user booking properties:', error);
//     // });
// }

//export const propertiesData=properties

