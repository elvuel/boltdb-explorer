export async function download(data, resp, success, failed) {
    await fetch(`${import.meta.env.VITE_API_URL}/download`, {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    }).then(
        response => resp(response) // if the response is a JSON object
    ).then(
        data => success(data) // Handle the success response object
    ).catch(
        error => failed(error) // Handle the error response object
    );
}