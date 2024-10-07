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
            if (data && data.url) {
                displayShortenedUrl(data.url);
            } else {
                alert('Error: Invalid response from server.');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('An error occurred while shortening the URL.');
        });
    } else {
        alert('Please enter a URL.');
    }
}

// Function to display the shortened URL
function displayShortenedUrl(shortUrl) {
    // Hide the URL input box
    const urlBox = document.getElementById('url-box');
    urlBox.style.display = 'none';

    // Show the shortened URL box
    const shortenedUrlBox = document.getElementById('shortened-url-box');
    const shortenedUrlLink = document.getElementById('shortened-url');

    shortenedUrlLink.href = shortUrl;
    shortenedUrlLink.textContent = shortUrl;

    shortenedUrlBox.style.display = 'block';
}

// Function to reset the form for another URL
function resetForm() {
    // Clear the input field
    document.getElementById('url-input').value = '';

    // Hide the shortened URL box
    document.getElementById('shortened-url-box').style.display = 'none';

    // Show the URL input box
    document.getElementById('url-box').style.display = 'block';
}
