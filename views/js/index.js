const numDeltaInput = document.getElementById('modifyNumberInput');
const modNumberButton = document.getElementById('modifyNumberSubmit');

const numberInput = document.getElementById('number');

// Event listeners for the buttons
modNumberButton.addEventListener('click', () => {
	let currentModValue = parseInt(numDeltaInput.value, 10) || 0;
	modifyStringWithPadding(currentModValue, 3);
});

numberInput.addEventListener('change', () => {
    modifyStringWithPadding(0, 3);
});

function modifyStringWithPadding(delta, minLength) {
    let currentStr = numberInput.value;
    let currentLen = currentStr.length;

    // Update minLength if the current string is longer than the provided minLength
    if (currentLen > minLength) {
        minLength = currentLen;
    }

    // Convert the current string to an integer or default to 0
    let currentValue = parseInt(currentStr, 10) || 0;
    currentValue += delta;

    // Ensure the value is non-negative
    if (currentValue < 0) {
        currentValue = 0;
    }

    // Convert back to string and pad to the required length, then set back to input
    numberInput.value = currentValue.toString().padStart(minLength, '0');
}

document.getElementById('jsonForm').addEventListener('submit', async function (e) {
	e.preventDefault(); // Prevent the default form submission behavior

	// Get form data as JSON
	const formData = new FormData(this);

	// Clear values if corresponding checkboxes are checked
	if (document.getElementById('clearGroup').checked) {
		document.getElementById('group').value = '';
	}
	if (document.getElementById('clearNumber').checked) {
		numberInput.value = '';
	}
	if (document.getElementById('clearName').checked) {
		document.getElementById('name').value = '';
	}

	document.getElementById('result').value = '';

	const jsonData = {};
	formData.forEach((value, key) => {
		jsonData[key] = value;
	});

	// Send JSON via fetch
	try {
		//const response = await fetch('./download', {
		const response = await fetch(submitPath, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(jsonData),
		});

		if (response.ok) {
			const result = await response.json();
			//alert('Form submitted successfully: ' + JSON.stringify(result));
			document.getElementById('result').value = 'Form submitted successfully: ' + JSON.stringify(result);
		} else {
			//alert('Error submitting form: ' + response.statusText);
			document.getElementById('result').value = 'Error submitting form: ' + response.statusText;
		}
	} catch (error) {
		//alert('Network error: ' + error.message);
		document.getElementById('result').value = 'Network error: ' + error.message;
	}
});