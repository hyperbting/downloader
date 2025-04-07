const numDeltaInput = document.getElementById('modifyNumberInput');
const modNumberButton = document.getElementById('modifyNumberSubmit');


const modNumberPlusOneButton = document.getElementById('modifyNumberPlusOneSubmit');
const modNumberMinusOneButton = document.getElementById('modifyNumberMinusOneSubmit');

const numberLengthInput = document.getElementById('numberLength');
const numberInput = document.getElementById('number');
// const numberPrefixInput = document.getElementById('numberPrefix');
// const numberSuffixInput = document.getElementById('numberSuffix');

const optionalNameInput = document.getElementById('name')

const formSubmit = document.getElementById('formSubmit');

let currentNumberLength = 3;

// Event listeners for Keyboards
document.addEventListener('keyup', (e) => {
	// special case: if focus on numberInput
	if (document.activeElement === numberInput) {
		return;
	}

	let skipFocus = false;
	switch (e.code) {
		case "ArrowUp":
			modifyStringWithPadding(10, currentNumberLength);
			break;
		case "ArrowDown":
			modifyStringWithPadding(-10, currentNumberLength);
			break;
		case "ArrowLeft":
			modifyStringWithPadding(-1, currentNumberLength);
			break;
		case "ArrowRight":
			modifyStringWithPadding(1, currentNumberLength);
			break;
		default:
			// Handle other keys or do nothing
			skipFocus = true;
			break;
	}

	if(!skipFocus){
		optionalNameInput.focus();
	}

});

// Event listeners for the buttons
modNumberButton.addEventListener('click', () => {
	let currentModValue = parseInt(numDeltaInput.value, 10) || 0;
	modifyStringWithPadding(currentModValue, currentNumberLength);
});

modNumberPlusOneButton.addEventListener('click', () => {
	modifyStringWithPadding(1, currentNumberLength);
});

modNumberMinusOneButton.addEventListener('click', () => {
	modifyStringWithPadding(-1, currentNumberLength);
});

numberInput.addEventListener('change', () => {
    modifyStringWithPadding(0, currentNumberLength);
});


numberLengthInput.addEventListener('change', () => {
	currentNumberLength = parseInt(numberLengthInput.value, 10) || 3;
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

	//console.log(jsonData);
	jsonData["number"] = String(jsonData["numberPrefix"] || "") +	String(jsonData["number"] || "") + String(jsonData["numberSuffix"] || "");

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