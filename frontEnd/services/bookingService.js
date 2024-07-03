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
    fetch(`${API_BASE_URL_PUBLIC}/Property`, {
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
        button.textContent = 'Book Now';
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
    bookNowButton.textContent = 'Confirm Booking';
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

function submitBooking() {
    const form = document.getElementById('bookingForm');
    const formData = new FormData(form);
    const noOfGuests= formData.get('noOfGuests');
    const checkInDate= formData.get('checkInDate');
    const checkOutDate= formData.get('checkOutDate');
    const email= formData.get('email');
    const message= formData.get('message');
    const amount= formData.get('amount');
    const propertyId= currentProperty._id;
    const userId=localStorage.getItem('UserId');
    // const bookingDetails = {
    //     noOfGuests: formData.get('noOfGuests'),
    //     checkInDate: formData.get('checkInDate'),
    //     checkOutDate: formData.get('checkOutDate'),
    //     email: formData.get('email'),
    //     message: formData.get('message'),
    //     amount: formData.get('amount'),
    //     propertyId: currentProperty._id,
    //     userId:localStorage.getItem('userId'),
    // };

    //console.log('Booking Details:', bookingDetails);

   //confirmBooking(bookingDetails)

   const confirmation = confirm(`Do you want to book the property "?`);
    if (confirmation) {
        fetch(`${API_BASE_URL_PUBLIC}/Booking`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ noOfGuests,checkOutDate,checkInDate,email,message,amount,propertyId,userId})
        })
        .then(response => response.json())
        .then(data => {
            renderProperties(data);
        })
        .catch(error => {
            console.error('Error fetching properties:', error);
        });
        alert('Booking confirmed!');
        closeBookingModal();
    }
}

function closeModal() {
    propertyModal.classList.remove('active');
}
function showUserBookings(){
    window.location.href = '../html/myBookings.html';
}

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

