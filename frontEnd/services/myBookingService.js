const API_BASE_URL_PUBLIC = window.config.API_BASE_URL_PUBLIC;

const propertiesData=JSON.parse(localStorage.getItem('properties'))

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


function fetchProperties() {
    fetch(`${API_BASE_URL_PUBLIC}/Booking/`+localStorage.getItem('UserId'), {
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
    })
    .then(response => response.json())
    .then(data => {
        renderProperties(data);
    })
    .catch(error => {
        console.error('Error fetching properties:', error);
    });
}

function renderProperties(mypropertiesData) {
    mypropertiesData.forEach(myProperty=>{
        propertiesData.forEach(property => {
            if(myProperty.propertyid==property._id){
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
                button.textContent = 'Cancel Booking';
                button.onclick = (event) => {
                    event.stopPropagation();
                    cancelBooking(myProperty._id);
                };
    
                buttonContainer.appendChild(button);
        
                card.appendChild(img);
                card.appendChild(title);
                card.appendChild(price);
                card.appendChild(buttonContainer);
        
                card.onclick = () => showPropertyDetails(property,myProperty._id);
        
                propertiesContainer.appendChild(card);
            }
        });
    })
}

function closeModal() {
    propertyModal.classList.remove('active');
}

// function closeBookingModal() {
//     bookingModal.classList.remove('active');
// }


function showBookingForm() {
    propertyModal.classList.remove('active');
    bookingModal.classList.add('active');
}

function showPropertyDetails(property,bookingId) {
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
    bookNowButton.textContent = 'Cancel Booking';
    bookNowButton.onclick = (event) => {
        event.stopPropagation();
        cancelBooking(bookingId);
    };

    modalContent.appendChild(bookNowButton);

    propertyModal.classList.add('active');
}

function cancelBooking(bookingId){
    const confirmation = confirm(`Do you want to cancel the booking "?`);
    if (confirmation) {
        fetch(`${API_BASE_URL_PUBLIC}\\`+bookingId, {
            method:'DELETE',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
                //'Content-Type': 'application/json'
            },
            //body: JSON.stringify({ noOfGuests,checkOutDate,checkInDate,email,message,amount,propertyId,userId})
        })
        .catch(error => {
            console.error('Error cancelling booking:', error);
        });
        alert('Booking cancelled!');
        closeModal();
        fetchProperties();
    }
}