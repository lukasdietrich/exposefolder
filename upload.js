document.addEventListener('dragover', preventDefault);
document.addEventListener('drop', handleDrop);

function preventDefault(event) {
	event.preventDefault();
}

async function handleDrop(event) {
	event.preventDefault();

	const { dataTransfer } = event;
	dataTransfer.dropEffect = 'copy';

	for (let file of dataTransfer.files) {
		await uploadFile(file);
	}

	location.reload();
}

async function uploadFile(file, method) {
	const formData = new FormData();
	formData.append("file", file, file.name);

	const res = await fetch(location.pathname, {
		method: method || 'POST' ,
		body: formData,
	});

	if (res.status === 409) {
		if (confirm(`${file.name} already exists. Overwrite?`)) {
			await uploadFile(file, 'PUT');
		}
	}
}
