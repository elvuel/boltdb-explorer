export async function upload(file, success, failed) {
    let data = new FormData();
    data.append('file', file)

    await fetch(`${import.meta.env.VITE_API_URL}/upload`, {
        method: 'POST',
        body: data
    }).then(
        response => response.json() // if the response is a JSON object
    ).then(
        data => success(data) // Handle the success response object
    ).catch(
        error => failed(error) // Handle the error response object
    );
}