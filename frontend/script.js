// Initialize the Vanta Fog effect
var setVanta = () => {
    if (window.VANTA) window.VANTA.FOG({
        el: ".vanta-bg",  // Targeting the 'vanta-bg' div
        mouseControls: true,
        touchControls: true,
        gyroControls: false,
        minHeight: 200.00,
        minWidth: 200.00
    });
};

// Trigger the Vanta effect when the page loads
setVanta();

// Function to handle URL shortening
function shortenUrl() {
    const url = document.getElementById('url-input').value;
    if (url) {
        fetch('http://localhost:8000/api/v1', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ url: url }),
        })
        .then(response => response.json())
        .then(data => {
            alert('Shortened URL: ' + data.url);
            fadeOutBox();
        })
        .catch(error => {
            console.error('Error:', error);
        });
    } else {
        alert('Please enter a URL.');
    }
}

// Function to fade out the input box
function fadeOutBox() {
    const box = document.getElementById('url-box');
    box.classList.add('fade-out');
}
